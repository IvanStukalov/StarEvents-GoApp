package models

type Item struct {
	ID          int
	Name        string
	Description string
	Distance    float32
	Magnitude   float32
	Image       string
	Age         float32
}

type List struct {
	Items []Item
}
