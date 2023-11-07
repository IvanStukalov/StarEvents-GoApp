package repository

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
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
	r.db.Where("NOT status = ?", models.StatusDeleted).Order("event_id").Find(&events, queryCondition)

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

	eventDetails := models.EventDetails{
		Event:     event,
		StarsList: stars,
	}

	return eventDetails, nil
}

func (r *Repository) ChangeEvent(eventId int, name string) error {
	var event models.Event
	r.db.Where("event_id = ?", eventId).First(&event)

	event.Name = name

	res := r.db.Save(&event)

	return res.Error
}

// put star into event
func (r *Repository) PutIntoEvent(eventMsg models.EventMsg) error {
	var draft models.Event
	r.db.Where("creator_id = ?", eventMsg.CreatorID).Where("status = ?", models.StatusCreated).First(&draft)

	if draft.ID == 0 {
		newEvent := models.Event{
			CreatorID:    eventMsg.CreatorID,
			ModeratorID:  r.GetModeratorId(),
			Status:       models.StatusCreated,
			CreationDate: time.Now(),
		}
		res := r.db.Create(&newEvent)
		if res.Error != nil {
			return res.Error
		}
		draft = newEvent
	}

	starEvent := models.StarEvents{
		EventID: draft.ID,
		StarID:  eventMsg.StarID,
	}

	res := r.db.Create(&starEvent)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *Repository) DeleteEvent(creatorId int) error {
	var event models.Event
	res := r.db.Where("status = ?", models.StatusCreated).First(&event, "creator_id = ?", creatorId)
	if res.Error != nil {
		return res.Error
	}

	event.Status = models.StatusDeleted
	res = r.db.Save(event)
	return res.Error
}

func (r *Repository) FormEvent(creatorId int) error {
	var event models.Event
	err := r.db.Where("status = ?", models.StatusCreated).First(&event, "creator_id = ?", creatorId)
	if err.Error != nil {
		return err.Error
	}

	event.Status = models.StatusFormed
	event.FormationDate = time.Now()

	res := r.db.Save(&event)

	return res.Error
}

func (r *Repository) ChangeEventStatus(eventId int, status string) error {
	var event models.Event

	err := r.db.Where("status = ?", models.StatusFormed).First(&event, "event_id = ?", eventId)
	if err.Error != nil {
		return err.Error
	}

	event.Status = status
	event.CompletionDate = time.Now()
	res := r.db.Save(&event)

	return res.Error
}

func (r *Repository) GetDraft(creatorId int) (int, error) {
	var event models.Event

	err := r.db.Where("status = ?", models.StatusCreated).First(&event, "creator_id = ?", creatorId)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return 0, err.Error
	}

	return event.ID, nil
}
