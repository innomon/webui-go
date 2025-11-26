package database

import (
	"backend/config"
	"backend/models"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

// ConnectDB connect to db with retry logic
func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		log.Println("Error parsing port")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))

	// Retry logic
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn))
		if err == nil {
			fmt.Println("Connection Opened to Database")
			// Migrate the schema
			DB.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{}, &models.File{}, &models.Folder{}, &models.Knowledge{}, &models.Model{}, &models.Prompt{}, &models.Tool{})
			fmt.Println("Database Migrated")
			return
		}
		log.Printf("Attempt %d to connect to database failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	panic("failed to connect database after multiple retries")
}
