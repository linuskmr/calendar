package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

// CalendarView fetches a calendar and its events within a date range
func (s *Server) CalendarView(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	beforeString := r.URL.Query().Get("before")
	afterString := r.URL.Query().Get("after")
	timezoneString := r.URL.Query().Get("timezone")

	timezone := time.UTC
	if timezoneString != "" {
		var err error
		timezone, err = time.LoadLocation(timezoneString)
		if err != nil {
			http.Error(w, "Invalid timezone", http.StatusBadRequest)
			return
		}
	}

	before, err := time.Parse(time.RFC3339, beforeString)
	if err != nil {
		http.Error(w, "Invalid before date", http.StatusBadRequest)
		return
	}
	after, err := time.Parse(time.RFC3339, afterString)
	if err != nil {
		http.Error(w, "Invalid after date", http.StatusBadRequest)
		return
	}

	var calendar Calendar
	result := s.db.Where("name = ?", name).First(&calendar)
	if result.Error != nil {
		log.Println("Error fetching calendar", result.Error)
		http.Error(w, "Error fetching calendar", http.StatusInternalServerError)
		return
	}

	var events []Event
	result = s.db.Where("calendar_id = ? AND start_date > ? AND end_date < ?", calendar.ID, after, before).Find(&events)
	if result.Error != nil {
		log.Println("Error fetching events", result.Error)
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		return
	}

	calendarView := CalendarView {
		After: after.In(timezone),
		Before: before.In(timezone),
		Calendar: Calendar {
			Name: name,
		},
		Events: events,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(calendarView)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}


// Event fetches an event by ID
func (s *Server) Event(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var event Event
	result := s.db.First(&event, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		} else {
			log.Println("Error fetching event", result.Error)
			http.Error(w, "Error fetching event", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}