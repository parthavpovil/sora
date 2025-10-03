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

	err = db.AutoMigrate(&models.User{})
    if err != nil {
        log.Printf("Failed to migrate User: %v", err)
    }

    err = db.AutoMigrate(&models.Room{})
    if err != nil {
        log.Printf("Failed to migrate Room: %v", err)
    }

    err = db.AutoMigrate(&models.Message{})
    if err != nil {
        log.Printf("Failed to migrate Message: %v", err)
    }

    err = db.AutoMigrate(&models.RoomMembers{})
    if err != nil {
        log.Printf("Failed to migrate RoomMembers: %v", err)
    }
	
	DB=db
	
}
