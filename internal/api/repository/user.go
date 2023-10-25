package repository

import "github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"

func (r *Repository) GetCreator() (models.User, error) {
	creator := models.User{}
	
	err := r.db.Find(&creator, "user_id = ?", 2).Error
	if err != nil {
		return models.User{}, err
	}

	return creator, nil
}

func (r *Repository) GetModerator() (models.User, error) {
	moderator := models.User{}
	
	err := r.db.Find(&moderator, "user_id = ?", 1).Error
	if err != nil {
		return models.User{}, err
	}

	return moderator, nil
}

func (r *Repository) GetCreatorId() (int, error) {
	return 2, nil
}

func (r *Repository) GetModeratorId() (int, error) {
	return 1, nil
}
