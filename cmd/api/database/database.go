package database

import (
	"GeoApi/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase() *gorm.DB {

	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&models.Zone{}, &models.Device{}, &models.Position{})
	if err != nil {
		return nil
	}

	return database
}
