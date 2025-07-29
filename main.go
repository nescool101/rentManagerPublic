package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nescool101/rentManager/controller"
	"github.com/nescool101/rentManager/service"
	"github.com/nescool101/rentManager/storage"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using default configurations.")
	}

	// Ensure payers file exists but DO NOT load them here
	storage.InitializePayersFile()

	// Load payers into memory
	service.LoadPayers()

	// Start scheduler
	go service.StartScheduler()

	// Start HTTP server
	if err := controller.StartHTTPServer(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
