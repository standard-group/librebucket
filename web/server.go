package web

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"librebucket/internal/db"
	"librebucket/internal/git"
)

// StartServer starts the web server for LibreBucket
func StartServer() {
	mux := http.NewServeMux()

	// Middleware: Request logging
	withLogging := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
			h(w, r)
			log.Printf("Completed in %v", time.Since(start))
		}
	}

	// API endpoints
	mux.Handle("/api/v1/users/register", withLogging(userRegisterHandler))
	mux.Handle("/api/v1/users/", withLogging(userAPIKeyHandler)) // Catches /api/v1/users/{username}/apikeys
	mux.Handle("/api/v1/git/create", withLogging(apiCreateRepoHandler))

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))

	// Generic handler for root, repo pages, and Git HTTP services
	mux.Handle("/", withLogging(gitAndWebHandler))
	mux.Handle("/login", withLogging(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "static/login.html")
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}))

	log.Println("Starting server on :3000...")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// gitAndWebHandler routes requests for Git HTTP services and general web UI
func gitAndWebHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")

	if path == "" {
		http.ServeFile(w, r, "static/index.html")
		return
	}

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		// Not enough parts for a /username/repo structure, serve index or 404
		http.ServeFile(w, r, "static/index.html") // Fallback
		return
	}

	username := parts[0]
	repoSegment := parts[1] // Can be "reponame" or "reponame.git"

	// Remove .git suffix for the actual bare repository name
	bareRepoName := repoSegment
	if strings.HasSuffix(repoSegment, ".git") {
		bareRepoName = strings.TrimSuffix(repoSegment, ".git")
	}
	repoPath := filepath.Join("repos", username, bareRepoName+".git")

	// Check if this is a Git HTTP transport request for info/refs or git-service
	if len(parts) >= 3 {
		service := parts[2] // e.g., "info", "git-upload-pack", "git-receive-pack"

		// Handle /username/repo[.git]/info/refs (GET or HEAD)
		if service == "info" && len(parts) >= 4 && parts[3] == "refs" && (r.Method == "GET" || r.Method == "HEAD") {
			handleGitInfoRefs(w, r, repoPath)
			return
		}

		// Handle /username/repo[.git]/git-upload-pack (POST)
		if service == "git-upload-pack" && r.Method == "POST" {
			handleGitService(w, r, repoPath, "git-upload-pack")
			return
		}

		// Handle /username/repo[.git]/git-receive-pack (POST)
		if service == "git-receive-pack" && r.Method == "POST" {
			handleGitService(w, r, repoPath, "git-receive-pack")
			return
		}
	}

	// If it's not a Git service request, check if the repository exists for web UI
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		http.NotFound(w, r) // Repository not found, serve 404
		return
	}

	// Serve dynamic web UI for a repository
	// This generates a simple HTML page for the repo
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>%s/%s - LibreBucket</title>
	<link rel="stylesheet" href="/css/style.css">
	<script src="/js/repo.js"></script>
</head>
<body>
	<div id="app"></div>
	<script>
		window.repoInfo = {
			username: '%s',
			repoName: '%s',
			cloneUrl: 'http://%s/%s/%s.git'
		};
	</script>
</body>
</html>`, username, bareRepoName, username, bareRepoName, r.Host, username, bareRepoName)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// apiCreateRepoHandler handles POST /api/v1/git/create
func apiCreateRepoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var req struct {
		Username string `json:"username"`
		RepoName string `json:"reponame"`
		Public   bool   `json:"public"`
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil || req.Username == "" || req.RepoName == "" {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON or missing fields")
		return
	}

	// Authenticate the user trying to create the repo
	authUsername, authPassword, authOK := getBasicAuth(r) // Correctly get password from Basic Auth
	if !authOK {
		writeJSONError(w, http.StatusUnauthorized, "Authentication required to create repository")
		return
	}

	// Use the password extracted from the Basic Auth header
	user, err := db.AuthenticateUser(authUsername, authPassword)
	if err != nil || user.Username != authUsername { // Double-check username if basic auth didn't match
		writeJSONError(w, http.StatusUnauthorized, "Invalid authentication for repo creation")
		return
	}
	// User must be the owner of the repo being created
	if user.Username != req.Username {
		writeJSONError(w, http.StatusForbidden, "Cannot create repository for another user")
		return
	}

	repoPath := filepath.Join("repos", req.Username, req.RepoName+".git")
	if err := git.CreateRepo(repoPath, req.Username, req.Public); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Use 201 Created for successful creation
	json.NewEncoder(w).Encode(map[string]any{
		"status":    "success",
		"clone_url": fmt.Sprintf("http://%s/%s/%s.git", r.Host, req.Username, req.RepoName),
	})
}

// userRegisterHandler handles POST /api/v1/users/register
func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"is_admin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON or missing fields")
		return
	}
	token, _ := GenerateToken()
	user, err := db.CreateUser(req.Username, req.Password, req.IsAdmin, token)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"token":  user.Token,
		"user":   user,
	})
}

// userAPIKeyHandler handles POST /api/v1/users/{username}/apikeys
func userAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 5 || parts[2] != "users" || parts[4] != "apikeys" {
		writeJSONError(w, http.StatusNotFound, "Not found")
		return
	}
	targetUsername := parts[3]

	// Authenticate the caller (must be the user themselves or an admin)
	authUsername, authPassword, authOK := getBasicAuth(r)
	if !authOK {
		writeJSONError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	authedUser, err := db.AuthenticateUser(authUsername, authPassword)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if authedUser.Username != targetUsername && !authedUser.IsAdmin {
		writeJSONError(w, http.StatusForbidden, "Not authorized to generate API key for this user")
		return
	}

	// Generate API key
	key, _ := GenerateToken()
	err = db.GenerateAPIKey(authedUser.ID, key)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"api_key": key,
	})
}

// handleGitInfoRefs handles GET/HEAD /username/repo[.git]/info/refs
func handleGitInfoRefs(w http.ResponseWriter, r *http.Request, repoPath string) {
	// Check if repository exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	service := r.URL.Query().Get("service")
	if service == "" {
		// Dumb HTTP protocol is typically not supported by modern Git clients
		// For smart protocol, service must be provided in query.
		http.Error(w, "Service parameter is required for smart HTTP protocol.", http.StatusBadRequest)
		return
	}

	var gitService string
	var action string // "pull" or "push"
	switch service {
	case "git-upload-pack":
		gitService = "upload-pack"
		action = "pull"
	case "git-receive-pack":
		gitService = "receive-pack"
		action = "push"
	default:
		http.Error(w, "Invalid Git service requested.", http.StatusBadRequest)
		return
	}

	// Check authorization
	username := filepath.Base(filepath.Dir(repoPath)) // Extract username from repoPath
	if !checkRepoAuth(r, repoPath, action, username) {
		w.Header().Set("WWW-Authenticate", `Basic realm="LibreBucket"`)
		w.WriteHeader(http.StatusUnauthorized)
		// No body for 401 on Git info/refs
		return
	}

	// Execute git command
	// --stateless-rpc --advertise-refs is for smart HTTP protocol for info/refs
	// Pass the repository path relative to the working directory
	cmd := exec.Command("git", gitService, "--stateless-rpc", "--advertise-refs", filepath.Base(repoPath))
	// Change working directory to the parent of the repository path
	cmd.Dir = filepath.Dir(repoPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Internal server error: failed to create stdout pipe.", http.StatusInternalServerError)
		return
	}
	stderr, err := cmd.StderrPipe() // Capture stderr for logging potential git errors
	if err != nil {
		http.Error(w, "Internal server error: failed to create stderr pipe.", http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		http.Error(w, "Internal server error: failed to start git process.", http.StatusInternalServerError)
		return
	}

	// Read stderr in background to avoid blocking
	go func() {
		if errText, err := io.ReadAll(stderr); err == nil && len(errText) > 0 {
			log.Printf("Git info/refs stderr: %s", string(errText))
		}
	}()

	// Set correct headers for Git info/refs advertisement
	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	// Write Git protocol service header and flush packet
	serviceHeader := fmt.Sprintf("# service=%s\n", service)
	w.Write(packetWrite(serviceHeader))
	w.Write(packetWrite("")) // Flush packet

	// Copy git command output to response
	io.Copy(w, stdout)

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		log.Printf("Git info/refs command (%s) exited with error: %v", service, err)
	}
}

// handleGitService handles POST /username/repo[.git]/git-upload-pack and git-receive-pack
func handleGitService(w http.ResponseWriter, r *http.Request, repoPath, service string) {
	// Check if repository exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	var gitServiceCmd string // Actual git command, e.g., "upload-pack" or "receive-pack"
	var action string        // "pull" or "push" for auth purposes
	switch service {
	case "git-upload-pack":
		gitServiceCmd = "upload-pack"
		action = "pull"
	case "git-receive-pack":
		gitServiceCmd = "receive-pack"
		action = "push"
	default:
		http.Error(w, "Invalid Git service requested.", http.StatusBadRequest)
		return
	}

	// Check authorization
	username := filepath.Base(filepath.Dir(repoPath)) // Extract username from repoPath
	if !checkRepoAuth(r, repoPath, action, username) {
		w.Header().Set("WWW-Authenticate", `Basic realm="LibreBucket"`)
		// Send 401 Unauthorized without a body for git service requests
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Execute git command
	// Pass the repository path relative to the working directory
	cmd := exec.Command("git", gitServiceCmd, "--stateless-rpc", filepath.Base(repoPath))
	// Change working directory to the parent of the repository path
	cmd.Dir = filepath.Dir(repoPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		http.Error(w, "Internal server error: failed to create stdin pipe.", http.StatusInternalServerError)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Internal server error: failed to create stdout pipe.", http.StatusInternalServerError)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		http.Error(w, "Internal server error: failed to create stderr pipe.", http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		http.Error(w, "Internal server error: failed to start git process.", http.StatusInternalServerError)
		return
	}

	// Set correct headers for Git service response
	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", service))
	w.WriteHeader(http.StatusOK)

	var wg sync.WaitGroup
	wg.Add(2) // Wait for both stdin and stdout/stderr pipes to finish

	// Goroutine to copy request body (client's data) to git's stdin
	go func() {
		defer wg.Done()
		defer stdin.Close() // Important: close stdin pipe when done copying
		if _, err := io.Copy(stdin, r.Body); err != nil {
			log.Printf("Error copying request body to git stdin: %v", err)
		}
	}()

	// Goroutine to copy git's stdout and stderr to the HTTP response
	go func() {
		defer wg.Done()
		// Git protocol requires streaming stdout and stderr to the client.
		// Prioritize stdout, then ensure stderr is also sent.
		if _, err := io.Copy(w, stdout); err != nil {
			log.Printf("Error copying git stdout to response: %v", err)
		}
		if _, err := io.Copy(w, stderr); err != nil {
			log.Printf("Error copying git stderr to response: %v", err)
		}
	}()

	wg.Wait() // Wait for both goroutines to complete I/O

	// Wait for the git command to finish.
	// If it exits with an error, the client will see it through the streamed stderr.
	if err := cmd.Wait(); err != nil {
		log.Printf("Git command (%s) exited with error: %v", gitServiceCmd, err)
	}
}

// packetWrite formats a Git protocol packet line
func packetWrite(s string) []byte {
	if s == "" {
		return []byte("0000")
	}
	length := len(s) + 4
	return []byte(fmt.Sprintf("%04x%s", length, s))
}

// writeJSONError writes a JSON error response
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// getBasicAuth extracts username and password from HTTP Basic Auth
func getBasicAuth(r *http.Request) (username, password string, ok bool) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		return
	}
	payload, err := base64.StdEncoding.DecodeString(auth[len("Basic "):])
	if err != nil {
		return
	}
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return
	}
	return pair[0], pair[1], true
}

// isOwnerAuthenticated checks if the request is authenticated as the repo owner using db
func isOwnerAuthenticated(r *http.Request, meta git.RepoMeta) bool {
	// 1. Try Basic Auth
	if username, password, ok := getBasicAuth(r); ok {
		user, err := db.AuthenticateUser(username, password)
		if err == nil && user.Username == meta.Owner {
			return true
		}
	}

	// 2. Try API Token (from query or header)
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("X-Auth-Token")
	}
	if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		token = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	}

	if token != "" {
		user, err := db.GetUserByToken(token)
		if err == nil && user.Username == meta.Owner {
			return true
		}
	}

	return false
}

// checkRepoAuth enforces public/private and owner rules for pull/push
func checkRepoAuth(r *http.Request, repoPath, action, expectedOwner string) bool {
	meta, err := git.LoadRepoMeta(repoPath)
	if err != nil {
		// If repo meta cannot be loaded, treat as unauthorized or non-existent
		return false
	}

	if action == "pull" { // Clone/Fetch (git-upload-pack)
		if meta.Public {
			return true // Public repos can be pulled by anyone
		}
		// Private repo: only owner can pull
		return isOwnerAuthenticated(r, meta)
	}

	// For push (git-receive-pack) and other write actions: only owner can push
	return isOwnerAuthenticated(r, meta)
}

// GenerateToken creates a 32-character alphanumeric token
func GenerateToken() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b), nil
}
