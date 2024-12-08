package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Calendar{}, &Event{})

	seeding(db)

	server := &Server{db: db}

	http.HandleFunc("/calendar/{name}", server.calendarView)
	http.HandleFunc("/event/{id}", server.event)
	port := ":8080"
	log.Println("Listening on port", port)
	http.ListenAndServe(port, nil)
}

type Server struct {
	db *gorm.DB
}

type Calendar struct {
	gorm.Model
	Name string `json:"name"`
	Events []Event `gorm:"foreignKey:CalendarID" json:"-"`
}

// CalendarView is a list of events for a calendar within a date range
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

func (s *Server) calendarView(w http.ResponseWriter, r *http.Request) {
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

	for _, event := range events {
		event.StartDate = event.StartDate.In(timezone)
		event.EndDate = event.EndDate.In(timezone)
		fmt.Println(event.StartDate, event.EndDate)
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

func (s *Server) event(w http.ResponseWriter, r *http.Request) {
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


func seeding(db *gorm.DB) {
	calendar := Calendar{
		Name: "test-calendar",
	}
	result := db.FirstOrCreate(&calendar, calendar)
	if result.Error != nil {
		panic(result.Error)
	}

	start, err := time.Parse(time.RFC3339, "2024-12-08T11:00:00+01:00")
	if err != nil {
		panic(err)
	}
	end, err := time.Parse(time.RFC3339, "2024-12-08T13:00:00+01:00")
	if err != nil {
		panic(err)
	}

	event := Event{
		Title: "Test Event",
		CalendarID: calendar.ID,
		StartDate: start,
		EndDate: end,
		Location: "Test Location",
		Description: "Test Description",
	}
	result = db.FirstOrCreate(&event, event)
	if result.Error != nil {
		panic(result.Error)
	}
}