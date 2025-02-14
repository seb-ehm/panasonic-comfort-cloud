package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/seb-ehm/panasonic-comfort-cloud/comfortcloud"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	username := os.Getenv("PANASONIC_USER")
	if username == "" {
		log.Fatal("PANASONIC_USER is not set in the .env file")
	}

	password := os.Getenv("PANASONIC_PASSWORD")
	if password == "" {
		log.Fatal("PANASONIC_PASSWORD is not set in the .env file")
	}
	auth := comfortcloud.NewAuthentication(username, password, nil, false)
	err = auth.GetNewToken()
	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println("Hello World")
}
