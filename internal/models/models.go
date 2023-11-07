package models

import "time"

const (
	SUN_MAGNITUDE           = -26.74
	MIN_MAGNITUDE           = SUN_MAGNITUDE
	UNIVERSAL_AGE           = 13.8
	VISIBLE_UNIVERSE_RADIUS = 4.65e10
)

type Star struct {
	ID          int     `json:"star_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Distance    float32 `json:"distance"`
	Age         float32 `json:"age"`
	Magnitude   float32 `json:"magnitude"`
	Image       string  `json:"image"`
	IsActive    bool    `json:"is_active"`
}

type Event struct {
	ID             int       `json:"event_id" gorm:"primaryKey;column:event_id;not null"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	ModeratorID    int       `json:"moderator_id"`
	CreatorID      int       `json:"creator_id"`
}

type EventDetails struct {
	Event     Event  `json:"event"`
	StarsList []Star `json:"stars_list"`
}

type StarEvents struct {
	StarID  int `json:"star_id"`
	EventID int `json:"event_id"`
}

type User struct {
	UserID      int    `json:"user_id" gorm:"primaryKey;column:user_id;not null"`
	Name        string `json:"name"`
	IsModerator bool   `json:"is_moderator"`
}
