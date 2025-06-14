package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"librebucket/internal/db"

	"github.com/go-chi/chi/v5"
)

// UserRegisterHandler handles POST /api/v1/users/register
func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
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

// UserLogInHandler handles POST /api/v1/users/login
func UserLogInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON or missing fields")
		return
	}
	user, err := db.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"token":  user.Token,
	})
}

// UserAPIKeyHandler handles POST /api/v1/users/{username}/apikeys
func UserAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	targetUsername := chi.URLParam(r, "username")
	if targetUsername == "" {
		writeJSONError(w, http.StatusNotFound, "Not found: username missing")
		return
	}

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

// --- Helpers (copied from web/server.go, but only if not used elsewhere) ---

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

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
