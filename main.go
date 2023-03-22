package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//Start logging
	go StartLogging()

	//Get the router and start server on port 8080
	router := SetupRouter()
	router.Run(":8080")
}
