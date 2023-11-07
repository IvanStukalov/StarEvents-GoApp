package pkg

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
)

func CastEvent(event models.Event, stars []models.Star) models.EventDetails {
	eventDetails := models.EventDetails{}

	eventDetails.Event.ID = event.ID
	eventDetails.Event.Name = event.Name
	eventDetails.Event.Status = event.Status
	eventDetails.Event.CreationDate = event.CreationDate
	eventDetails.Event.FormationDate = event.FormationDate
	eventDetails.Event.CompletionDate = event.CompletionDate
	eventDetails.Event.ModeratorID = event.ModeratorID
	eventDetails.Event.CreatorID = event.CreatorID
	eventDetails.StarsList = stars

	return eventDetails
}
