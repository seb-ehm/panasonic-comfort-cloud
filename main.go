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

	deviceID := os.Getenv("PANASONIC_DEVICE_ID")
	if deviceID == "" {
		log.Fatal("PANASONIC_DEVICE_ID is not set in the .env file")
	}
	//auth := comfortcloud.NewAuthentication(username, password, nil)
	//err = auth.GetNewToken()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("Refreshing token")
	//err = auth.RefreshToken()

	c := comfortcloud.NewClient(username, password, ".panasonic-oauth-token")
	c.Login()
	err = c.FetchGroupsAndDevices()
	if err != nil {
		fmt.Println(err)
	}
	device, err := c.GetDevice("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(device)
	fmt.Println("#######")

	err = c.SetDevice(deviceID,
		comfortcloud.WithPower(comfortcloud.PowerOn),
		comfortcloud.WithTemperature(22.0))
	if err != nil {
		fmt.Println(err)
	}
}
