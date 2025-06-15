// internal/api/commit.go
package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"librebucket/internal/git"
)

// CommitHandler handles commit-related API endpoints
func CommitHandler(mux *http.ServeMux) {
	// Get commit history for a repository
	mux.HandleFunc("GET /api/v1/repos/{username}/{reponame}/commits", getCommitHistory)

	// Get a specific commit by hash
	mux.HandleFunc("GET /api/v1/repos/{username}/{reponame}/commits/{hash}", getCommit)

	// Get file changes in a commit
	mux.HandleFunc("GET /api/v1/repos/{username}/{reponame}/commits/{hash}/changes", getCommitChanges)

	// Get file content at a specific commit
	mux.HandleFunc("GET /api/v1/repos/{username}/{reponame}/blob/{hash}/{filepath...}", getFileAtCommit) // fix for the getFileAtCommit
}

func getCommitHistory(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	reponame := r.PathValue("reponame")

	repoPath := getRepoPath(username, reponame)
	commits, err := git.GetCommitHistory(repoPath)
	if err != nil {
		http.Error(w, "Failed to get commit history: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(commits); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func getCommit(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	reponame := r.PathValue("reponame")
	hash := r.PathValue("hash")

	repoPath := getRepoPath(username, reponame)
	commit, err := git.GetCommitByHash(repoPath, hash)
	if err != nil {
		http.Error(w, "Failed to get commit: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(commit); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func getCommitChanges(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	reponame := r.PathValue("reponame")
	hash := r.PathValue("hash")

	repoPath := getRepoPath(username, reponame)
	changes, err := git.GetCommitChanges(repoPath, hash)
	if err != nil {
		http.Error(w, "Failed to get commit changes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(changes); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// detectContentType determines the content type of a file based on its extension and/or content
func detectContentType(filePath string) string {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(filePath))

	// Check common extensions first
	switch ext {
	case ".go":
		return "text/x-go"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".md":
		return "text/markdown"
	case ".txt":
		return "text/plain"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".xml":
		return "application/xml"
	case ".zip":
		return "application/zip"
	case ".sh":
		return "application/x-sh"
	case ".py":
		return "text/x-python"
	}

	// If we can't determine from extension, try to read the file
	// This is optional and can be removed if you want to avoid file I/O
	if ext == "" || ext == ".bin" {
		// Open the file
		file, err := os.Open(filePath)
		if err == nil {
			defer file.Close()

			// Read the first 512 bytes to determine content type
			buffer := make([]byte, 512)
			_, err = file.Read(buffer)
			if err == nil {
				// Use http.DetectContentType to get MIME type
				return http.DetectContentType(buffer)
			}
		}
	}

	// Default to binary if we couldn't determine the type
	return "application/octet-stream"
}

func getFileAtCommit(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	reponame := r.PathValue("reponame")
	hash := r.PathValue("hash")
	filepathA := r.PathValue("filepath")

	repoPath := getRepoPath(username, reponame)
	content, err := git.GetFileAtCommit(repoPath, filepathA, hash)
	if err != nil {
		http.Error(w, "Failed to get file at commit: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set appropriate content type based on file extension
	contentType := detectContentType(filepathA)
	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}

// Helper function to construct repository path
func getRepoPath(username, reponame string) string {
	// This should come from your configuration
	repoRoot := "./repos" // Default, should be configurable
	return repoRoot + "/" + username + "/" + reponame
}
