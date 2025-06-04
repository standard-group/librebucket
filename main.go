package main

import (
	"librebucket/internal/db"
	"librebucket/web"
	"log"
)

func main() {
	if err := db.InitDB("users.db"); err != nil {
		log.Fatalf("Failed to initialize user DB: %v", err)
	}
	web.StartServer()
}
