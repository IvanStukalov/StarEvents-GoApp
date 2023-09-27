package models

type Star struct {
	ID          int
	Name        string
	Description string
	Distance    float32
	Age         float32
	Magnitude   float32
	Image       string
	IsActive    bool
}
