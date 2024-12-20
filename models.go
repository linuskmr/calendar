package main

import (
	"time"

	"gorm.io/gorm"
)


// EventList in the range between start and end
type Events struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Events   []Event   `json:"events"`
}

// Event in database
type Event struct {
	gorm.Model
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}