package db

import (
	"fmt"

	"gorm.io/driver/postgres"

	"github.com/EugeneKrivoshein/qa_api_service/models"
	"gorm.io/gorm"
)

func Connect(host, user, password, dbname, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автоматическая миграция моделей
	if err := db.AutoMigrate(&models.Question{}, &models.Answer{}); err != nil {
		return nil, err
	}

	return db, nil
}
