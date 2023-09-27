package repository

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
	log "github.com/sirupsen/logrus"
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
	var star []models.Star

	r.db.Find(&star)
	log.Println("len", len(star))
	return star, nil
}

func (r *Repository) CreateProduct(star models.Star) error {
	return r.db.Create(star).Error
}
