package repository

import (
	// "strconv"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
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

// func (r *Repository) GetStarByID(starId int) (models.Star, error) {
// 	star := models.Star{}

// 	r.db.Find(&star, "id = ?", strconv.Itoa(starId))

// 	return star, nil
// }

func (r *Repository) GetStarsByNameFilter(substring string) ([]models.Star, error) {
	var star []models.Star

	r.db.Where("name ILIKE ?", "%" + substring + "%").Find(&star, "is_active = ?", true)
	return star, nil
}

// func (r *Repository) DeleteStarById(starId int) error {
// 	err := r.db.Exec("UPDATE stars SET is_active=false WHERE id = ?", starId).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *Repository) UpdateStar(starID int, star models.Star) error {
// 	var lastStar models.Star
// 	r.db.First(&lastStar, starID)
	
// 	if (len(star.Name) != 0) {
// 		lastStar.Name = star.Name
// 	}

// 	if (len(star.Description) != 0) {
// 		lastStar.Description = star.Description
// 	}

// 	if (len(star.Image) != 0) {
// 		lastStar.Image = star.Image
// 	}

// 	if (star.Age != -1) {
// 		lastStar.Age = star.Age
// 	}

// 	if (star.Distance != -1) {
// 		lastStar.Distance = star.Distance
// 	}

// 	if (star.Magnitude != 100) {
// 		lastStar.Magnitude = star.Magnitude
// 	}

// 	r.db.Save(&lastStar)

// 	return nil
// }

func (r *Repository) CreateStar(star models.Star) error {
	var count int64
	r.db.Model(&models.Star{}).Count(&count)

	err := r.db.Create(&star).Error
	if err != nil {
		return err
	}
	
	return nil
}
