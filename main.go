package main

import (
	"librebucket/internal/db"
	"librebucket/web"
	"log"
	"os"
)

func main() {
	if err := db.InitDB("users.db"); err != nil {
		log.Fatalf("Failed to initialize user DB: %v", err)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println(dir)

	web.StartServer()
}
