package main

import (
	"log"
	"os"
	"path/filepath"

	"librebucket/cmd/db"
	"librebucket/cmd/web"
)

// main initializes the application's data directory and user database, then starts the web server.
// It terminates execution with a fatal log if any critical setup step fails.
func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// ./config/data
	dataDir := filepath.Join(wd, "config", "data")

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	dbPath := filepath.Join(dataDir, "users.db")

	if err := db.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize user DB: %v", err)
		return
	}

	log.Println("Working dir:", wd)
	log.Println("DB initialized at:", dbPath)

	web.StartServer()
}
