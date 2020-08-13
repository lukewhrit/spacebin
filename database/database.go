package database

import (
	"log"

	"github.com/gobuffalo/pop"
)

// DBConn holds the current connection to the database
var DBConn *pop.Connection

// Init opens a connection to the database
func Init() {
	var err error

	DBConn, err = pop.Connect("main")

	if err != nil {
		log.Fatalf("Failed to connect to database: %e", err)
	}
}
