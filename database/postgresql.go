package database

import (
	"helpdesk-ticketing-system/internal/helper"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres() *gorm.DB {
	dsn := helper.GetConnectionString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}
