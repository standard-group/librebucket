package db

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// InitDB initializes the SQLite database and creates tables if needed
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	// Create users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		token TEXT UNIQUE NOT NULL,
		is_admin INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		return err
	}
	// Create api_keys table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS api_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		key TEXT UNIQUE NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)`)
	return err
}

// User represents a user account
type User struct {
	ID           int
	Username     string
	PasswordHash string
	Token        string
	IsAdmin      bool
	APIKeys      []string
}

// CreateUser creates a new user with a hashed password and random token
func CreateUser(username, password string, isAdmin bool, token string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	res, err := db.Exec(`INSERT INTO users (username, password_hash, token, is_admin) VALUES (?, ?, ?, ?)`, username, string(hash), token, boolToInt(isAdmin))
	if err != nil {
		return User{}, err
	}
	id, _ := res.LastInsertId()
	return User{ID: int(id), Username: username, PasswordHash: string(hash), Token: token, IsAdmin: isAdmin}, nil
}

// AuthenticateUser checks username and password
func AuthenticateUser(username, password string) (User, error) {
	var u User
	row := db.QueryRow(`SELECT id, username, password_hash, token, is_admin FROM users WHERE username = ?`, username)
	var isAdminInt int
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Token, &isAdminInt); err != nil {
		return User{}, errors.New("user not found")
	}
	u.IsAdmin = isAdminInt != 0
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return User{}, errors.New("invalid password")
	}
	u.APIKeys, _ = GetAPIKeys(u.ID)
	return u, nil
}

// GetUserByToken returns a user by login token
func GetUserByToken(token string) (User, error) {
	var u User
	row := db.QueryRow(`SELECT id, username, password_hash, token, is_admin FROM users WHERE token = ?`, token)
	var isAdminInt int
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Token, &isAdminInt); err != nil {
		return User{}, errors.New("user not found for token")
	}
	u.IsAdmin = isAdminInt != 0
	u.APIKeys, _ = GetAPIKeys(u.ID)
	return u, nil
}

// GenerateAPIKey creates a new API key for a user and saves it
func GenerateAPIKey(userID int, key string) error {
	_, err := db.Exec(`INSERT INTO api_keys (user_id, key) VALUES (?, ?)`, userID, key)
	return err
}

// GetAPIKeys returns all API keys for a user
func GetAPIKeys(userID int) ([]string, error) {
	rows, err := db.Query(`SELECT key FROM api_keys WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var keys []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err == nil {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

// ValidateAPIKey checks if the API key is valid and returns the user
func ValidateAPIKey(key string) (User, error) {
	row := db.QueryRow(`SELECT u.id, u.username, u.password_hash, u.token, u.is_admin FROM users u JOIN api_keys k ON u.id = k.user_id WHERE k.key = ?`, key)
	var u User
	var isAdminInt int
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Token, &isAdminInt); err != nil {
		return User{}, errors.New("invalid API key")
	}
	u.IsAdmin = isAdminInt != 0
	u.APIKeys, _ = GetAPIKeys(u.ID)
	return u, nil
}

// GetUserByBearerToken checks for Bearer token in Authorization header
func GetUserByBearerToken(authHeader string) (User, error) {
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return User{}, errors.New("not a bearer token")
	}
	token := authHeader[7:]
	// Try login token first
	user, err := GetUserByToken(token)
	if err == nil {
		return user, nil
	}
	// Try API key
	return ValidateAPIKey(token)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
