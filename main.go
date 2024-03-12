package main

import (
	"fmt"
	"log"
	"os"
	"sec-infra/firewall"
	"sec-infra/proxy"
	"sec-infra/vpn"

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

	containerName := os.Getenv("CONTAINER_NAME")

	fmt.Println(welcomeMessage)

	proxy.Run()

	firewall.SetupFirewall(containerName)

	vpn.Run()

}
