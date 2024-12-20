package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

// ApiCalendarView fetches a calendar and its events within a date range
func (s *Server) ApiCalendarView(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	startString := r.URL.Query().Get("start")
	endString := r.URL.Query().Get("end")
	timezoneString := r.URL.Query().Get("timezone")

	calendarView := s.CalendarViewController(w, name, startString, endString, timezoneString)
	if calendarView == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(calendarView)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func (s *Server) CalendarView(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	startString := r.URL.Query().Get("start")
	endString := r.URL.Query().Get("end")
	timezoneString := r.URL.Query().Get("timezone")

	if startString == "" {
		// Use today's date at midnight
		now := time.Now()
		startString = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
	}
	if endString == "" {
		// Use today's date at midnight + 7 days
		endString = time.Now().Add(time.Hour * 24 * 7).Format(time.RFC3339)
	}

	calendarView := s.CalendarViewController(w, name, startString, endString, timezoneString)
	if calendarView == nil {
		return
	}

	calendarViewTemplate, err := template.ParseFiles("templates/base.html", "templates/calendar.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = calendarViewTemplate.Execute(w, calendarView)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
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
	result := s.db.Take(&event, id)
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

func (s *Server) ApiAddEvent(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request
	var addEvent AddEvent
	err := json.NewDecoder(r.Body).Decode(&addEvent)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	if !s.AddEventController(w, addEvent) {
		return // Error is already written
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(addEvent)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func (s *Server) ShowAddEventForm(w http.ResponseWriter, r *http.Request) {
	addEventFormTemplate, err := template.ParseFiles("templates/base.html", "templates/addEvent.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = addEventFormTemplate.Execute(w, struct {
		CalendarName string
	}{
		CalendarName: r.PathValue("name"),
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func (s *Server) AddFormEvent(w http.ResponseWriter, r *http.Request) {
	// Decode the form-encoded body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	start, err := time.Parse(time.RFC3339, r.FormValue("start"))
	if err != nil {
		http.Error(w, "Invalid start timestamp", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, r.FormValue("end"))
	if err != nil {
		http.Error(w, "Invalid end timestamp", http.StatusBadRequest)
		return
	}

	addEvent := AddEvent{
		CalendarName: r.FormValue("calendar_name"),
		Title:        r.FormValue("title"),
		Start:        start,
		End:          end,
		Location:     r.FormValue("location"),
		Description:  r.FormValue("description"),
	}

	if !s.AddEventController(w, addEvent) {
		return // Error was already written
	}

	http.Redirect(w, r, "/calendar/"+addEvent.CalendarName+"/events", http.StatusSeeOther)
}
