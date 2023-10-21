package repository

import (
	"strconv"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg"
)

func (r *Repository) GetEventList() ([]models.Event, error) {
	var events []models.Event

	r.db.Find(&events)

	return events, nil
}

func (r *Repository) GetEventByID(eventId int) (models.EventDetails, error) {
	event := models.Event{}
	r.db.Find(&event, "event_id = ?", strconv.Itoa(eventId))

	var starEvents []models.StarEvents
	r.db.Find(&starEvents, "event_id = ?", strconv.Itoa(eventId))

	var stars []models.Star
	for i := 0; i < len(starEvents); i++ {
		star := models.Star{}
		row := r.db.Find(&star, "star_id = ? AND is_active = ?", starEvents[i].StarID, true)
		if row.RowsAffected > 0 {
			stars = append(stars, star)
		}
	}

	eventDetails := pkg.CastEvent(event, stars)

	return eventDetails, nil
}
