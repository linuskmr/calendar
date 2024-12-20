package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func (s *Server) ApiEventList(w http.ResponseWriter, r *http.Request) {
	startString := r.URL.Query().Get("start")
	endString := r.URL.Query().Get("end")
	timezoneString := r.URL.Query().Get("timezone")

	calendarView := s.EventListController(w, startString, endString, timezoneString)
	if calendarView == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(calendarView)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// ApiGetEvent fetches an event by ID
func (s *Server) ApiGetEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var event Event
	result := s.Db.Take(&event, id)
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

func (s *Server) ApiAddOrUpdateEvent(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Insert or update the event in the database
	result := s.Db.Save(&event)
	if result.Error != nil {
		log.Println("Error creating event", result.Error)
		http.Error(w, "Error creating event", http.StatusInternalServerError)
		return
	}

	// Return created event, including its id
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
