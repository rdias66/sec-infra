package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Access variables from the environment
	welcomeMessage := os.Getenv("WELCOME_MESSAGE")

	fmt.Println(welcomeMessage)

}
