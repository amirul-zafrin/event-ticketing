package database

import (
	"log"

	"github.com/amirul-zafrin/event-ticketing/orders.git/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstances struct {
	Db *gorm.DB
}

var Database DBInstances

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("events.db"), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database")

	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running migrations")
	db.AutoMigrate(&models.Orders{})

	Database = DBInstances{Db: db}
}
