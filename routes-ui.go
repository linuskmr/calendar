package main

import (
	"fmt"
	"log"
	"math/bits"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"gorm.io/gorm"
)

// Events returns a list of events in the range between start and end, with timestamps formatted in the specified timezone.
func (s *Server) EventList(w http.ResponseWriter, r *http.Request) {
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

	eventList := s.EventListController(w, startString, endString, timezoneString)
	if eventList == nil {
		return // Error was already responded
	}

	eventListTemplate, err := template.ParseFiles("templates/base.html", "templates/event/list.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = eventListTemplate.Execute(w, eventList)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func (s *Server) ShowAddEventForm(w http.ResponseWriter, r *http.Request) {
	addEventFormTemplate, err := template.ParseFiles("templates/base.html", "templates/event/form.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = addEventFormTemplate.Execute(w, struct {
		FormActionUrl string
		Event Event
	}{
		FormActionUrl: "/events",
		Event: Event{},
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func (s *Server) ShowEditEventForm(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid event id", http.StatusBadRequest)
		return
	}

	// Fetch event from database
	var event Event
	result := s.Db.Where("id = ?", eventId).Take(&event)
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

	editEventFormTemplate, err := template.ParseFiles("templates/base.html", "templates/event/form.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = editEventFormTemplate.Execute(w, struct {
		FormActionUrl string
		Event Event
	}{
		FormActionUrl: "/events/" + fmt.Sprint(eventId),
		Event: Event{},
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func (s *Server) AddOrUpdateFormEvent(w http.ResponseWriter, r *http.Request) {
	// Decode the form-encoded body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	datetimeLocalFormat := "2006-01-02T15:04"

	start, err := time.Parse(datetimeLocalFormat, r.FormValue("start"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid start timestamp", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(datetimeLocalFormat, r.FormValue("end"))
	if err != nil {
		http.Error(w, "Invalid end timestamp", http.StatusBadRequest)
		return
	}

	eventIdString := r.FormValue("event_id")
	var eventId uint64 = 0
	if eventIdString != "" {
		var err error
		eventId, err = strconv.ParseUint(eventIdString, 10, bits.UintSize)
		if err != nil {
			http.Error(w, "Invalid event id", http.StatusBadRequest)
			return
		}
	}

	event := Event{
		ID:          uint(eventId),
		Title:       r.FormValue("title"),
		Start:       start,
		End:         end,
		Location:    r.FormValue("location"),
		Description: r.FormValue("description"),
	}

	// Insert or update the event in the database
	result := s.Db.Save(&event)
	if result.Error != nil {
		log.Println("Error creating event", result.Error)
		http.Error(w, "Error creating event", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/events", http.StatusSeeOther)
}