package api

import (
	"time"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
)

type Repo interface {
	GetFilteredStars(substring string) ([]models.Star, error)
	GetStarByID(starId int) (models.Star, error)
	DeleteStarByID(starId int) error
	UpdateStar(star models.Star) error
	GetStarImageById(starId int) (string, error)
	CreateStar(star models.Star) error
	PutIntoEvent(starEvent models.StarEvents) error
	RemoveFromEvent(starEvent models.StarEvents) error

	GetEventList(status string, startFormation time.Time, endFormation time.Time) ([]models.Event, error)
	GetEventByID(eventId int) (models.EventDetails, error)
	UpdateEvent(event models.Event) error
	CreateEvent(event models.Event) (models.Event, error)
	FormEvent(eventId int) (models.Event, error)
	CompleteEvent(eventId int) (models.Event, error)
	RejectEvent(eventId int) (models.Event, error)
	DeleteEvent(eventId int) error

	GetCreator() (models.User, error)
	GetModerator() (models.User, error)
	GetCreatorId() (int, error)
	GetModeratorId() (int, error)
}
