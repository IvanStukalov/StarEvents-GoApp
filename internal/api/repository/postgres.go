package repository

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
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

func (r *Repository) GetStarByID(starId int) (models.Star, error) {
	star := models.Star{}

	r.db.Find(&star, "id = ?", strconv.Itoa(starId))

	return star, nil
}

func (r *Repository) GetStars() ([]models.Star, error) {
	var star []models.Star

	r.db.Find(&star, "is_active = ?", true)
	return star, nil
}

func (r *Repository) DeleteStarById(starId int) error {
	err := r.db.Exec("UPDATE stars SET is_active=false WHERE id = ?", starId).Error
	if err != nil {
		return err
	}
	return nil
}
