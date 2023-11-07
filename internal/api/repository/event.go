package repository

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg"
)

// get list of events
func (r *Repository) GetEventList(status string, startFormation time.Time, endFormation time.Time) ([]models.Event, error) {
	var events []models.Event
	var queryCondition string
	var wasPrevCond bool

	if status != "" {
		queryCondition = fmt.Sprintf("status = '%s'", status)
		wasPrevCond = true
	}

	if !startFormation.IsZero() {
		if wasPrevCond {
			queryCondition += " AND "
		}
		queryCondition += fmt.Sprintf("formation_date > '%v'", startFormation.Format(time.DateTime))
		wasPrevCond = true
	}

	if !endFormation.IsZero() {
		if wasPrevCond {
			queryCondition += " AND "
		}
		queryCondition += fmt.Sprintf("formation_date < '%v'", endFormation.Format(time.DateTime))
	}

	log.Println(queryCondition)
	r.db.Order("event_id").Find(&events, queryCondition)

	return events, nil
}

// get event by ID
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

func (r *Repository) UpdateEvent(event models.Event) error {
	var lastEvent models.Event
	r.db.Select("status").Where("event_id = ?", event.ID).First(&lastEvent)

	event.Status = lastEvent.Status

	if err := r.db.Save(&event).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateEvent(event models.Event) (models.Event, error) {
	if err := r.db.Create(&event).Error; err != nil {
		return models.Event{}, err
	}

	return event, nil
}

func (r *Repository) FormEvent(eventId int) (models.Event, error) {
	event := models.Event{}

	err := r.db.Find(&event, "event_id = ?", eventId).Error
	if err != nil {
		return models.Event{}, err
	}

	event.Status = "formed"
	event.FormationDate = time.Now()

	err = r.db.Save(&event).Error
	if err != nil {
		return models.Event{}, err
	}

	return event, nil
}

func (r *Repository) CompleteEvent(eventId int) (models.Event, error) {
	event := models.Event{}

	err := r.db.Find(&event, "event_id = ?", eventId).Error
	if err != nil {
		return models.Event{}, err
	}

	event.Status = "fulfilled"
	event.CompletionDate = time.Now()

	err = r.db.Save(&event).Error
	if err != nil {
		return models.Event{}, err
	}

	return event, nil
}

func (r *Repository) RejectEvent(eventId int) (models.Event, error) {
	event := models.Event{}

	err := r.db.Find(&event, "event_id = ?", eventId).Error
	if err != nil {
		return models.Event{}, err
	}

	event.Status = "rejected"
	event.CompletionDate = time.Now()

	err = r.db.Save(&event).Error
	if err != nil {
		return models.Event{}, err
	}

	return event, nil
}

func (r *Repository) DeleteEvent(eventId int) error {
	event := models.Event{}

	err := r.db.Where("event_id = ?", eventId).Delete(&event).Error
	if err != nil {
		return err
	}

	return nil
}
