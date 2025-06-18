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
	return err
}

// User represents a user account
type User struct {
	ID           int
	Username     string
	PasswordHash string
	Token        string
	IsAdmin      bool
}

// CreateUser creates a new user with a hashed password and random token
func CreateUser(username, password string, _ bool, token string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	res, err := db.Exec(`INSERT INTO users (username, password_hash, token, is_admin) VALUES (?, ?, ?, 0)`, username, string(hash), token)
	if err != nil {
		return User{}, err
	}
	id, _ := res.LastInsertId()
	return User{ID: int(id), Username: username, PasswordHash: string(hash), Token: token, IsAdmin: false}, nil
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
	return User{}, errors.New("invalid API key")
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
