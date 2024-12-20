package main

import (
	"time"

	"gorm.io/gorm"
)

type Calendar struct {
	gorm.Model
	Name   string  `json:"name"`
	Events []Event `gorm:"foreignKey:CalendarID" json:"-"`
}

// CalendarView is a list of events within a date range of a calendar
type CalendarView struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Calendar Calendar  `json:"calendar"`
	Events   []Event   `json:"events"`
}

// Event in database
type Event struct {
	gorm.Model
	CalendarID  uint      `json:"-"`
	Title       string    `json:"title"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}

// AddEvent is the request to add an event
type AddEvent struct {
	CalendarName string    `json:"calendar_name"`
	Title        string    `json:"title"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	Location     string    `json:"location"`
	Description  string    `json:"description"`
}
