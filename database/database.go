package database

import (
	"log"

	"github.com/parthavpovil/sora/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb()  {
	dsn := "host=localhost user=postgres password=yourpassword dbname=soradb port=5434 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&models.User{}, &models.Message{})
	DB=db
	
}
