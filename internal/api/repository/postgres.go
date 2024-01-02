package repository

import (
	"StarEvent-GoApp/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(connectionString string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Star{})
	if err != nil {
		panic("cant migrate db")
	}

	return &Repository{
		db: db,
	}, nil
}
