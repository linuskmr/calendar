package main

import (
	"time"

	"gorm.io/gorm"
)

type Calendar struct {
	gorm.Model
	Name string `json:"name"`
	Events []Event `gorm:"foreignKey:CalendarID" json:"-"`
}

// CalendarView is a list of events within a date range of a calendar
type CalendarView struct {
	After time.Time `json:"after"`
	Before time.Time `json:"before"`
	Calendar Calendar `json:"calendar"`
	Events []Event `json:"events"`
}

type Event struct {
	gorm.Model
	CalendarID uint `json:"-"`
	Title string `json:"title"`
	StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date"`
	Location string `json:"location"`
	Description string `json:"description"`
}