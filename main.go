package main

import (
	"log"

	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	err = db.Connect()
	if err != nil {
		log.Fatal("failed to laod database")
	}
}
