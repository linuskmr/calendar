package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/calendar/{name}", calendarView)
	http.HandleFunc("/event/{id}", event)
	port := ":8080"
	log.Println("Listening on port", port)
	http.ListenAndServe(port, nil)
}

type Repository struct {
	calendars []Calendar
}

type Calendar struct {
	Name string `json:"name"`
}

// CalendarView is a list of events for a calendar within a date range
type CalendarView struct {
	After time.Time `json:"after"`
	Before time.Time `json:"before"`
	Calendar Calendar `json:"calendar"`
	Events []Event `json:"events"`
}

type Event struct {
	Id int `json:"id"`
	Title string `json:"title"`
	StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date"`
	Location string `json:"location"`
	Description string `json:"description"`
}

func calendarView(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	beforeString := r.URL.Query().Get("before")
	afterString := r.URL.Query().Get("after")

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

	demoEvent := Event{
		Id: 123,
		Title: "My Event",
		StartDate: time.Now(),
		EndDate: time.Now().Add(24 * time.Hour),
		Location: "My House",
		Description: "This is my event",
	}

	events := []Event{}
	for _ = range 100 {
		events = append(events, demoEvent)
	}

	calendarView := CalendarView {
		After: after,
		Before: before,
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

func event(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event := Event {
		Id: id,
		Title: "My Event",
		StartDate: time.Now(),
		EndDate: time.Now().Add(24 * time.Hour),
		Location: "My House",
		Description: "This is my event",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}