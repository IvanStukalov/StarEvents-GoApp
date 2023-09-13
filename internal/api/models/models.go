package models

type Item struct {
	Name        string
	Description string
	Distance    float32
	Magnitude   float32
	Image       string
}

type List struct {
	Items []Item
}
