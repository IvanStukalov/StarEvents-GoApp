package pkg

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
)

func CastEvent(event models.Event, stars []models.Star) models.EventDetails {
	eventDetails := models.EventDetails{}

	eventDetails.ID = event.ID
	eventDetails.Name = event.Name
	eventDetails.Status = event.Status
	eventDetails.CreationDate = event.CreationDate
	eventDetails.FormationDate = event.FormationDate
	eventDetails.CompletionDate = event.CompletionDate
	eventDetails.ModeratorID = event.ModeratorID
	eventDetails.CreatorID = event.CreatorID
	eventDetails.StarsList = stars

	return eventDetails
}
