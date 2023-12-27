package api

import (
	"time"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
)

type Repo interface {
	GetFilteredStars(substring string, distTop string, distBot string, ageTop string, ageBot string, magTop string, magBot string) ([]models.Star, error)
	GetStarByID(starId int) (models.Star, error)
	DeleteStarByID(starId int) error
	UpdateStar(star models.Star) error
	GetStarImageById(starId int) (string, error)
	CreateStar(star models.Star) error
	RemoveFromEvent(starEvent models.StarEvents) error

	GetEventList(status string, startFormation time.Time, endFormation time.Time) ([]models.Event, error)
	GetEventByID(eventId int) (models.Event, []models.Star, error)
	UpdateEvent(eventId int, name string) error
	PutIntoEvent(eventMsg models.EventMsg) error
	DeleteEvent(creatorId int) error
	FormEvent(creatorId int) error
	ChangeEventStatus(eventId int, status string) error
	GetDraft(creatorId int) (int, error)
	SaveScannedPercent(eventAsync models.EventAsync) error

	RemoveStarFromEvent(creatorId int, starId int) (models.Event, []models.Star, error)

	GetCreatorId() int
	GetModeratorId() int
}
