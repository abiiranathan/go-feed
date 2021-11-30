package db

import (
	"log"
	"os"

	"github.com/abiiranathan/gin_api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

// init postgres connection
func Connect() {
	DSN := os.Getenv("DSN")
	var err error
	DBConn, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	log.Printf("Connected to database...")

	DBConn.AutoMigrate(&models.Feed{})

}
