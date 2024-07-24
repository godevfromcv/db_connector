package main

import (
	"db_connector/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("failed to run the application: %v", err)
	}
}
