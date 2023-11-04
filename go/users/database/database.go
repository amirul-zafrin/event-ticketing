package database

import (
	"fmt"
	"log"

	"github.com/amirul-zafrin/event-ticketing/users.git/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	Db *gorm.DB
}

var Database DBInstance

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("apiv2.db"), &gorm.Config{})
	if err != nil {
		// log.Fatal("Failed to connect to database! \n", err.Error())
		// os.Exit(1)
		errorString := fmt.Sprintf("Failed to connect to database!\n%v", err.Error())
		panic(errorString)
	}
	log.Println("Connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	db.AutoMigrate(&models.User{})

	Database = DBInstance{Db: db}
}
