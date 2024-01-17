package repository

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"

	"StarEvent-GoApp/internal/models"
)

// get list of events
func (r *Repository) GetEventList(status string, startFormation time.Time, endFormation time.Time, creatorId int, isAdmin bool) ([]models.Event, error) {
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
		wasPrevCond = true
	}

	if creatorId != 0 && !isAdmin {
		if wasPrevCond {
			queryCondition += " AND "
		}
		queryCondition += fmt.Sprintf("creator_id = %v", creatorId)
		wasPrevCond = true
	}

	r.db.Where("NOT status = ?", models.StatusDeleted).Order("event_id").Find(&events, queryCondition)

	for i := range events {
		var creator models.User
		var moderator models.User
		r.db.Find(&creator, "user_id = ?", events[i].CreatorID)
		events[i].Creator = creator.Login

		if events[i].ModeratorID != 0 {
			r.db.Find(&moderator, "user_id = ?", events[i].ModeratorID)
			events[i].Moderator = moderator.Login
		}
	}

	return events, nil
}

// get event by ID
func (r *Repository) GetEventByID(eventId int, creatorId int, isAdmin bool) (models.Event, []models.Star, error) {
	event := models.Event{}
	var err error
	if isAdmin {
		r.db.Find(&event, "event_id = ?", strconv.Itoa(eventId))
	} else {
		err = r.db.Where("creator_id = ?", strconv.Itoa(creatorId)).Find(&event, "event_id = ?", strconv.Itoa(eventId)).Error
	}

	if event.ID == 0 {
		return models.Event{}, []models.Star{}, err
	}

	var starEvents []models.StarEvents
	r.db.Find(&starEvents, "event_id = ?", strconv.Itoa(event.ID))

	var stars []models.Star
	for i := 0; i < len(starEvents); i++ {
		star := models.Star{}
		row := r.db.Find(&star, "star_id = ? AND is_active = ?", starEvents[i].StarID, true)
		if row.RowsAffected > 0 {
			stars = append(stars, star)
		}
	}

	var creator models.User
	var moderator models.User
	r.db.Find(&creator, "user_id = ?", event.CreatorID)
	event.Creator = creator.Login

	if event.ModeratorID != 0 {
		r.db.Find(&moderator, "user_id = ?", event.ModeratorID)
		event.Moderator = moderator.Login
	}

	return event, stars, nil
}

func (r *Repository) UpdateEvent(eventId int, name string) error {
	var event models.Event
	r.db.Where("event_id = ?", eventId).First(&event)

	event.Name = name

	res := r.db.Save(&event)

	return res.Error
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

func (r *Repository) FormEvent(creatorId int) (error, int) {
	var event models.Event
	err := r.db.Where("status = ?", models.StatusCreated).First(&event, "creator_id = ?", creatorId)
	if err.Error != nil {
		return err.Error, 0
	}

	var starEvent = []models.StarEvents{}
	err = r.db.Find(&starEvent, "event_id = ?", event.ID)
	if err.Error != nil {
		return err.Error, 0
	}

	if len(starEvent) == 0 {
		return fmt.Errorf("заявка пуста"), 0
	}

	event.Status = models.StatusFormed
	event.FormationDate = time.Now()

	res := r.db.Save(&event)
	if res.Error != nil {
		return res.Error, 0
	}

	return nil, event.ID
}

func (r *Repository) ChangeEventStatus(eventId int, status string, moderatorId int) error {
	var event models.Event

	err := r.db.Where("status = ?", models.StatusFormed).First(&event, "event_id = ?", eventId)
	if err.Error != nil {
		return err.Error
	}

	event.Status = status
	event.CompletionDate = time.Now()
	event.ModeratorID = moderatorId

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

func (r *Repository) SaveScannedPercent(eventAsync models.EventAsync) error {
	var event models.Event
	err := r.db.First(&event, "event_id = ?", eventAsync.ID)
	if err.Error != nil {
		return err.Error
	}
	event.ScannedPercent = eventAsync.ScannedPercent
	res := r.db.Save(&event)
	return res.Error
}
