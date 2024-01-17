package api

import (
	"context"
	"time"

	"StarEvent-GoApp/internal/models"
)

type Repo interface {
	GetFilteredStars(substring string, distTop string, distBot string, ageTop string, ageBot string, magTop string, magBot string) ([]models.Star, error)
	GetStarByID(starId int) (models.Star, error)
	DeleteStarByID(starId int) error
	UpdateStar(star models.Star, isUpdA bool, isUpdD bool, isUpdM bool) error
	GetStarImageById(starId int) (string, error)
	CreateStar(star models.Star) error
	PutIntoEvent(eventMsg models.EventMsg) (int, error)

	GetEventList(status string, startFormation time.Time, endFormation time.Time, creatorId int, isAdmin bool) ([]models.Event, error)
	GetEventByID(eventId int, creatorId int, isAdmin bool) (models.Event, []models.Star, error)
	UpdateEvent(eventId int, name string) error
	DeleteEvent(creatorId int) error
	FormEvent(creatorId int) (error, int)
	ChangeEventStatus(eventId int, status string, moderatorId int) error
	GetDraft(creatorId int) (int, error)
	SaveScannedPercent(eventAsync models.EventAsync) error

	RemoveStarFromEvent(creatorId int, starId int) (models.Event, []models.Star, error)

	SignUp(ctx context.Context, newUser models.User) error
	GetByCredentials(ctx context.Context, user models.User) (models.User, error)
	GetUserInfo(ctx context.Context, user models.User) (models.User, error)
}
