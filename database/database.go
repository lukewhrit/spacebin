package database

import (
	"log"

	"github.com/spacebin-org/curiosity/config"
	"github.com/spacebin-org/curiosity/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBConn holds the current connection to the database
var DBConn *gorm.DB

// Init opens a connection to the database
func Init() {
	var err error
	var dialect gorm.Dialector

	switch config.Config.Database.Dialect {
	case "sqlite":
		dialect = sqlite.Open(config.Config.Database.ConnectionURI)
	case "postgresql":
		dialect = postgres.Open(config.Config.Database.ConnectionURI)
	case "mysql":
		dialect = mysql.Open(config.Config.Database.ConnectionURI)
	}

	DBConn, err = gorm.Open(dialect, &gorm.Config{})

	DBConn.AutoMigrate(&models.Document{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %e", err)
	}
}
