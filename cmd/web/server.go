package web

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	api "librebucket/cmd/api/v1"
	"librebucket/cmd/db"
	"librebucket/cmd/git"

	"gopkg.in/yaml.v3"
)

// StartServer initializes and runs the LibreBucket web server, setting up API endpoints, static file serving, Git HTTP protocol handlers, and web UI routes. The server listens on the specified port and terminates with a fatal log message if it fails to start.
func StartServer() {
	port := flag.Int("port", 3000, "Port to listen on")
	flag.Parse()

	r := chi.NewRouter()

	// Middleware: Request logging
	r.Use(middleware.Logger)

	// API endpoints
	r.Post("/api/v1/users/register", api.UserRegisterHandler)
	r.Post("/api/v1/users/login", api.UserLogInHandler)
	// r.Post("/api/v1/users/{username}/apikeys", api.UserAPIKeyHandler)
	r.Post("/api/v1/git/create", api.APICreateRepoHandler)

	// Commits API endpoints (mount ServeMux from api.CommitHandler)
	commitMux := http.NewServeMux()
	api.CommitHandler(commitMux)
	r.Mount("/api/v1/repos", commitMux)

	// Serve static files from the cmd/web/static directory
	fs := http.FileServer(http.Dir("cmd/web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Page handlers
	r.Get("/", homeHandler)
	r.Post("/set-lang", setLangHandler)
	r.Get("/login", loginHandler)

	// Git HTTP services
	r.Get("/{username}/{repoName}.git/info/refs", handleGitInfoRefs)
	r.Head("/{username}/{repoName}.git/info/refs", handleGitInfoRefs)
	r.Post("/{username}/{repoName}.git/git-upload-pack", handleGitService)
	r.Post("/{username}/{repoName}.git/git-receive-pack", handleGitService)

	r.Get("/{username}/{repoName}/info/refs", handleGitInfoRefs)
	r.Head("/{username}/{repoName}/info/refs", handleGitInfoRefs)
	r.Post("/{username}/{repoName}/git-upload-pack", handleGitService)
	r.Post("/{username}/{repoName}/git-receive-pack", handleGitService)

	// Repository web UI pages
	r.Get("/{username}/{repoName}", gitAndWebHandler)
	r.Get("/{username}/{repoName}.git", gitAndWebHandler) // Handles paths with .git suffix

	log.Printf("Starting server on :%d...", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func isSafeComponent(s string) bool {
	return !strings.Contains(s, "/") &&
		!strings.Contains(s, "\\") &&
		!strings.Contains(s, "..")
}

// gitAndWebHandler serves the web UI page for a Git repository if it exists and the path is valid.
//
// Validates the username and repository name from the URL, ensures the repository exists, and renders the repository's web interface with relevant information. Responds with HTTP 400 for invalid paths or 404 if the repository is not found.
func gitAndWebHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	repoName := chi.URLParam(r, "repoName")

	// Ensure repoName does not have .git suffix for consistency
	repoName = strings.TrimSuffix(repoName, ".git")

	if !isSafeComponent(username) || !isSafeComponent(repoName) {
		http.Error(w, "Invalid repo path", http.StatusBadRequest)
		return
	}

	repoPath := filepath.Join("repos", username, repoName+".git")

	// If it's not a Git service request, check if the repository exists for web UI
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		http.NotFound(w, r) // Repository not found, serve 404
		return
	}

	// Data to pass to the template
	data := map[string]any{
		"username": username,
		"repoName": repoName,
		"cloneUrl": fmt.Sprintf("http://%s/%s/%s.git", r.Host, username, repoName),
	}

	// Render the Go HTML template
	RenderTemplate("repo.tmpl", data, w)
}

// handleGitInfoRefs serves the Git smart HTTP protocol's info/refs endpoint for a repository.
//
// Validates the repository path and requested Git service, checks authorization for pull or push actions,
// and executes the corresponding Git command to advertise refs. Streams the Git protocol response to the client
// with appropriate headers. Responds with 401 if authentication is required or 404 if the repository does not exist.
func handleGitInfoRefs(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	repoName := chi.URLParam(r, "repoName")

	// Ensure repoName does not have .git suffix for consistency
	repoName = strings.TrimSuffix(repoName, ".git")

	if !isSafeComponent(username) || !isSafeComponent(repoName) {
		http.Error(w, "Invalid repo path", http.StatusBadRequest)
		return
	}
	repoPath := filepath.Join("repos", username, repoName+".git")

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
	if !checkRepoAuth(r, repoPath, action, username) {
		w.Header().Set("WWW-Authenticate", `Basic realm="LibreBucket"`)
		w.WriteHeader(http.StatusUnauthorized)
		// No body for 401 on Git info/refs
		return
	}

	// Execute git command
	// --stateless-rpc --advertise-refs is for smart HTTP protocol for info/refs
	// Pass the repository path relative to the working directory
	cmd := exec.Command("git", gitService, "--stateless-rpc", "--advertise-refs", "--", filepath.Base(repoPath))
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

// handleGitService processes Git smart HTTP POST requests for upload-pack (pull) and receive-pack (push) services on a repository.
//
// It validates the repository path and service type, checks user authorization for the requested action, and executes the corresponding Git command in stateless RPC mode. The function streams the request body to the Git process (handling gzip compression if present) and relays the Git command's output to the HTTP response with appropriate headers. Responds with relevant HTTP errors for invalid paths, unauthorized access, or internal failures.
func handleGitService(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	repoName := chi.URLParam(r, "repoName")

	// Ensure repoName does not have .git suffix for consistency
	repoName = strings.TrimSuffix(repoName, ".git")

	if !isSafeComponent(username) || !isSafeComponent(repoName) {
		http.Error(w, "Invalid repo path", http.StatusBadRequest)
		return
	}
	repoPath := filepath.Join("repos", username, repoName+".git")

	// Check if repository exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	var gitServiceCmd string
	var action string
	var serviceName string // This will be "git-upload-pack" or "git-receive-pack"

	// correctly determine the service from the request URL path
	if strings.HasSuffix(r.URL.Path, "/git-upload-pack") {
		serviceName = "git-upload-pack"
		gitServiceCmd = "upload-pack"
		action = "pull"
	} else if strings.HasSuffix(r.URL.Path, "/git-receive-pack") {
		serviceName = "git-receive-pack"
		gitServiceCmd = "receive-pack"
		action = "push"
	} else {
		http.Error(w, "Invalid Git service requested.", http.StatusBadRequest)
		return
	}

	// check authorizationo
	if !checkRepoAuth(r, repoPath, action, username) {
		w.Header().Set("WWW-Authenticate", `Basic realm="LibreBucket"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cmd := exec.Command("git", gitServiceCmd, "--stateless-rpc", "--", filepath.Base(repoPath))
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
	// capture stderr for logging
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		http.Error(w, "Internal server error: failed to start git process.", http.StatusInternalServerError)
		return
	}

	// handle gzip compressed request body from the git client*
	go func() {
		defer stdin.Close()
		var reader io.Reader = r.Body
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("Error creating gzip reader for git request: %v", err)
				return
			}
			reader = gz
			defer gz.Close()
		}

		if _, err := io.Copy(stdin, reader); err != nil {
			log.Printf("Error copying request body to git stdin: %v", err)
		}
	}()

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", serviceName))
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	// copy git's stdout to the HTTP response
	if _, err := io.Copy(w, stdout); err != nil {
		log.Printf("Error copying git stdout to response: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Git command (%s) exited with error: %v", gitServiceCmd, err)
		// log any stderr
		if stderr.Len() > 0 {
			log.Printf("Git stderr: %s", stderr.String())
		}
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

// getLang gets language from cookie, defaults to "en"
func getLang(r *http.Request) string {
	cookie, err := r.Cookie("lang")
	if err == nil && len(cookie.Value) == 2 {
		return cookie.Value
	}
	return "en"
}

// loadTranslations loads translations from a YAML file for a given page and language
func loadTranslations(lang, page string) (map[string]any, error) {
	path := fmt.Sprintf("cmd/web/i18n/langs/%s/%s.yaml", lang, page)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var m map[string]any
	if err := yaml.NewDecoder(f).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

// homeHandler serves the home page with translations
func homeHandler(w http.ResponseWriter, r *http.Request) {
	lang := getLang(r)
	trans, err := loadTranslations(lang, "home")
	if err != nil {
		// Fallback to English if the language file is missing or fails to parse
		log.Printf("Could not load translations for lang '%s': %v. Falling back to 'en'.", lang, err)
		trans, _ = loadTranslations("en", "home")
	}

	data := map[string]any{
		"Trans": trans,
		"Lang":  lang,
	}
	RenderTemplate("home.tmpl", data, w)
}

// loginHandler serves the login page with translations
func loginHandler(w http.ResponseWriter, r *http.Request) {
	lang := getLang(r)
	trans, err := loadTranslations(lang, "login")
	if err != nil {
		log.Printf("Could not load translations for lang '%s': %v. Falling back to 'en'.", lang, err)
		trans, _ = loadTranslations("en", "login") // fallback
	}
	RenderTemplate("login.tmpl", map[string]any{
		"Trans": trans,
		"Lang":  lang,
	}, w)
}

// setLangHandler sets the language cookie
func setLangHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err == nil {
		lang := r.FormValue("lang")
		http.SetCookie(w, &http.Cookie{
			Name:   "lang",
			Value:  lang,
			Path:   "/",
			MaxAge: 86400 * 365, // 1 year
		})
	}
	// Redirect back to the previous page
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
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
