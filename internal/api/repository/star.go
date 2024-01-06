package repository

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"StarEvent-GoApp/internal/models"
)

// get stars with filter
func (r *Repository) GetFilteredStars(substring string, distTop string, distBot string, ageTop string, ageBot string, magTop string, magBot string) ([]models.Star, error) {
	var star []models.Star
	var queryCondition = "is_active = true"

	if distTop != "" {
		queryCondition = fmt.Sprintf(queryCondition+" AND distance <= '%s'", distTop)
	}

	if distBot != "" {
		queryCondition = fmt.Sprintf(queryCondition+" AND distance >= '%s'", distBot)
	}

	if ageTop != "" {
		queryCondition = fmt.Sprintf(queryCondition+" AND age <= '%s'", ageTop)
	}

	if ageBot != "" {
		queryCondition = fmt.Sprintf(queryCondition+" AND age >= '%s'", ageBot)
	}

	if magTop != "" {
		queryCondition = fmt.Sprintf(queryCondition+" AND magnitude <= '%s'", magTop)
	}

	if magBot != "" {
		queryCondition = fmt.Sprintf(queryCondition+" AND magnitude >= '%s'", magBot)
	}

	if len(substring) != 0 {
		// if query substring exists
		err := r.db.Order("star_id").Where("name ILIKE ?", "%"+substring+"%").Find(&star, queryCondition).Error
		if err != nil {
			log.Println(err)
			return []models.Star{}, err
		}

	} else {
		// if query substring is empty
		err := r.db.Order("star_id").Find(&star, queryCondition).Error
		if err != nil {
			log.Println(err)
			return []models.Star{}, err
		}
	}

	return star, nil
}

// get star by id
func (r *Repository) GetStarByID(starId int) (models.Star, error) {
	var star models.Star

	err := r.db.Find(&star, "star_id = ?", strconv.Itoa(starId)).Error
	if err != nil {
		log.Println(err)
		return models.Star{}, err
	}

	return star, nil
}

// get star image by id
func (r *Repository) GetStarImageById(starId int) (string, error) {
	var star models.Star

	err := r.db.First(&star, "star_id = ?", starId).Error
	if err != nil {
		return "", err
	}

	return star.Image, nil
}

// create star
func (r *Repository) CreateStar(star models.Star) error {
	var newStar models.Star

	newStar.Name = star.Name
	newStar.Description = star.Description
	newStar.Image = star.Image
	newStar.IsActive = true

	if star.Age >= 0 && star.Age <= models.UniversalAge {
		newStar.Age = star.Age
	} else {
		return errors.New("star age must be greater than 0 and less than Universal age (13.8 billion years)")
	}

	if star.Distance >= 0 && star.Distance <= models.VisibleUniverseRadius {
		newStar.Distance = star.Distance
	} else {
		return errors.New("star distance must be greater than 0 and less than visible universe radius (4.65e+10 l.y.)")
	}

	if star.Magnitude >= models.MinMagnitude {
		newStar.Magnitude = star.Magnitude
	} else {
		return errors.New("star magnitude must be greater than minimum possible magnitude (-26.74 - Sun magnitude)")
	}

	err := r.db.Create(&newStar).Error
	if err != nil {
		return err
	}

	return nil
}

// update star
func (r *Repository) UpdateStar(star models.Star, isUpdA bool, isUpdD bool, isUpdM bool) error {
	fmt.Println(star)

	var lastStar models.Star

	err := r.db.First(&lastStar, star.ID).Error
	if err != nil {
		return err
	}

	if star.Name != "" {
		lastStar.Name = star.Name
	}

	if star.Description != "" {
		lastStar.Description = star.Description
	}

	if star.Image != "" {
		lastStar.Image = star.Image
	}

	if star.Age >= 0 && star.Age <= models.UniversalAge && isUpdA {
		lastStar.Age = star.Age
	}

	if star.Distance >= 0 && star.Distance <= models.VisibleUniverseRadius && isUpdD {
		lastStar.Distance = star.Distance
	}

	if star.Magnitude >= models.MinMagnitude && isUpdM {
		lastStar.Magnitude = star.Magnitude
	}

	err = r.db.Save(&lastStar).Error
	if err != nil {
		return err
	}

	return nil
}

// delete star by id
func (r *Repository) DeleteStarByID(starId int) error {
	err := r.db.Exec("UPDATE stars SET is_active=false WHERE star_id = ?", starId).Error
	if err != nil {
		return err
	}
	return nil
}

// put star into event
func (r *Repository) PutIntoEvent(eventMsg models.EventMsg) (int, error) {
	var draft models.Event
	r.db.Where("creator_id = ?", eventMsg.CreatorID).Where("status = ?", models.StatusCreated).First(&draft)


	if draft.ID == 0 {
		newEvent := models.Event{
			CreatorID:    eventMsg.CreatorID,
			Status:       models.StatusCreated,
			CreationDate: time.Now(),
		}
		res := r.db.Create(&newEvent)
		if res.Error != nil {
			return 0, res.Error
		}
		draft = newEvent
	}

	starEvent := models.StarEvents{
		EventID: draft.ID,
		StarID:  eventMsg.StarID,
	}

	res := r.db.Create(&starEvent)
	if res.Error != nil {
		return 0, res.Error
	}

	return draft.ID, nil
}
