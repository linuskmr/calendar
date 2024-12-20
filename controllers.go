package main

import (
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func (s *Server) CalendarViewController(w http.ResponseWriter, name, startString, endString, timezoneString string) *CalendarView {
	timezone := time.UTC
	if timezoneString != "" {
		var err error
		timezone, err = time.LoadLocation(timezoneString)
		if err != nil {
			http.Error(w, "Invalid timezone", http.StatusBadRequest)
			return nil
		}
	}

	start, err := time.Parse(time.RFC3339, startString)
	if err != nil {
		http.Error(w, "Invalid start timestamp", http.StatusBadRequest)
		return nil
	}
	end, err := time.Parse(time.RFC3339, endString)
	if err != nil {
		http.Error(w, "Invalid end timestamp", http.StatusBadRequest)
		return nil
	}

	var calendar Calendar
	result := s.db.Where("name = ?", name).First(&calendar)
	if result.Error != nil {
		log.Println("Error fetching calendar", result.Error)
		http.Error(w, "Error fetching calendar", http.StatusInternalServerError)
		return nil
	}

	var events []Event
	// Fetch events during the specified timestamp range (start - end).
	// Therefore, we explicitly also want to include events that are only partially within the range.
	// Consequently, the condition is either that they have started before the end or have not ended before the start.
	result = s.db.Where("calendar_id = ? AND (start <= ? OR end >= ?)", calendar.ID, end, start).Find(&events)
	if result.Error != nil {
		log.Println("Error fetching events", result.Error)
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		return nil
	}

	return &CalendarView{
		Start: start.In(timezone),
		End:   end.In(timezone),
		Calendar: Calendar{
			Name: name,
		},
		Events: events,
	}
}

func (s *Server) AddEventController(w http.ResponseWriter, addEvent AddEvent) bool {
	// Find the calendar id for the specified calendar name
	var calendar Calendar
	result := s.db.Where("name = ?", addEvent.CalendarName).Take(&calendar)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Calendar for specified calendar name not found", http.StatusNotFound)
			return false
		} else {
			log.Println("Error fetching corresponding calendar", result.Error)
			http.Error(w, "Error fetching corresponding calendar", http.StatusInternalServerError)
			return false
		}
	}
	log.Println("calendar name", addEvent.CalendarName, "calendar id", calendar.ID)

	// Transform the AddEvent to an Event
	event := Event{
		CalendarID:  calendar.ID,
		Title:       addEvent.Title,
		Start:       addEvent.Start,
		End:         addEvent.End,
		Location:    addEvent.Location,
		Description: addEvent.Description,
	}

	// Insert the event into the database
	result = s.db.Create(&event)
	if result.Error != nil {
		log.Println("Error creating event", result.Error)
		http.Error(w, "Error creating event", http.StatusInternalServerError)
		return false
	}

	return true
}
