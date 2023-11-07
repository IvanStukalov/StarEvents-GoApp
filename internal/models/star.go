package models

type Star struct {
	ID          int     `json:"star_id" gorm:"primaryKey;column:star_id;not null"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Distance    float32 `json:"distance"`
	Age         float32 `json:"age"`
	Magnitude   float32 `json:"magnitude"`
	Image       string  `json:"image"`
	IsActive    bool    `json:"is_active"`
}
