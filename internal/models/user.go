package models

type User struct {
	UserID      int    `json:"user_id" gorm:"primaryKey;column:user_id;not null"`
	Name        string `json:"name"`
	IsModerator bool   `json:"is_moderator"`
}
