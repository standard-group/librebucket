package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"librebucket/internal/db"
	"librebucket/internal/git"
)

// APICreateRepoHandler handles POST /api/v1/git/create
func APICreateRepoHandler(w http.ResponseWriter, r *http.Request) {
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

	// Authenticate the user trying to create the repo (get user by the token)
	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.Header.Get("X-Auth-Token")
	}
	if token == "" {
		token = r.URL.Query().Get("token")
	}

	var user db.User
	var err error
	if strings.HasPrefix(token, "Bearer ") {
		user, err = db.GetUserByBearerToken(token)
	} else {
		user, err = db.GetUserByToken(token)
	}
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Authentication required to create repository")
		return
	}
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

// --- Helpers (copied from user.go if not used elsewhere) ---
// Only include if not already present in user.go or other shared location
// func writeJSONError ...
// func getBasicAuth ...
