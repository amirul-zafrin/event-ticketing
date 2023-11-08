package database

import (
	"fmt"
	"log"

	"github.com/amirul-zafrin/event-ticketing/events.git/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	Db *gorm.DB
}

var Database DBInstance

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("events.db"), &gorm.Config{})
	if err != nil {
		errorString := fmt.Sprintf("Failed to connect to database!\n%v", err.Error())
		panic(errorString)
	}
	log.Println("Connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	db.AutoMigrate(&models.Events{}, &models.Prices{})

	Database = DBInstance{Db: db}
}
