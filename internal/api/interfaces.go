package api

import "github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"

type Repo interface {
	GetStars() ([]models.Star, error)
	GetStarByID(threatId int) (models.Star, error)
	DeleteStarById(starId int) error
}
