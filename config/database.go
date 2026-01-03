package config

import (
	"fmt"
	"hotel-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=admin dbname=hotel_db port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal konek ke database...")
	}

	database.AutoMigrate(&models.Room{}, &models.Reservation{})
	DB = database
	fmt.Println("✅ Database Connected & Migrated!")
}
