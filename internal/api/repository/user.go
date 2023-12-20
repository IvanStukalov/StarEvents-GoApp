package repository

func (r *Repository) GetCreatorId() int {
	return 2
}

func (r *Repository) GetModeratorId() int {
	return 1
}

func (r *Repository) GetCreator() string {
	return "Алексей"
}

func (r *Repository) GetModerator() string {
	return "Владимир"
}
