package db

import (
	"log"
	"os"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	var err error

	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Connecting to database failed")
	}

	err = db.AutoMigrate(
		&models.Role{},
	)
	if err != nil {
		log.Fatal("Database auto migration failed")
	}
}

func Get() *gorm.DB {
	return db
}
