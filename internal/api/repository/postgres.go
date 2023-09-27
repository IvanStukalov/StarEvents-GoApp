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

	err := r.db.First(&star, "star_id = ?", strconv.Itoa(starId)).Error
	if err != nil {
		return star, err
	}

	return star, nil
}

// ////////////////////////////////////////////////////////////////////////////////
func (r *Repository) GetStarByName(name string) (models.Star, error) {
	star := models.Star{}

	err := r.db.First(&star, "name = ?", name).Error
	if err != nil {
		return star, err
	}
	return star, nil
}

func (r *Repository) GetStars() ([]models.Star, error) {
	stars := make([]models.Star, 0, 4)

	r.db.Where("is_active = ?", true).Find(&stars)

	return stars, nil
}

func (r *Repository) CreateProduct(star models.Star) error {
	return r.db.Create(star).Error
}
