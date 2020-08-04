package main

import (
	"log"

	"github.com/spacebin-org/curiosity/config"
	"github.com/spacebin-org/curiosity/database"
	"github.com/spacebin-org/curiosity/server"
)

func main() {
	// Load config
	if err := config.Load(); err != nil {
		log.Fatalf("Couldn't load configuration file: %v", err)
	}

	// Start server and initialize database
	server.Start()
	database.Init()
}
