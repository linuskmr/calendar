package main

import (
	"log"
	"net/http"
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

	http.HandleFunc("/calendar/{name}", server.CalendarView)
	http.HandleFunc("/event/{id}", server.Event)
	port := ":8080"
	log.Println("Listening on port", port)
	http.ListenAndServe(port, nil)
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