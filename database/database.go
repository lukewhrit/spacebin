package database

import (
	"github.com/jinzhu/gorm"
)

var (
	// DBConn holds the current connection to the database
	DBConn *gorm.DB
)
