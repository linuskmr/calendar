package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	port := flag.String("addr", ":8080", "Address to listen on")
	dbFile := flag.String("db", "test.db", "Database file")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open(*dbFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Calendar{}, &Event{})

	seeding(db)

	server := &Server{db: db}

	http.HandleFunc("GET /api/calendar/{name}/events", server.ApiCalendarView)
	http.HandleFunc("GET /api/event/{id}", server.ApiGetEvent)
	http.HandleFunc("POST /api/event", server.ApiAddEvent)

	http.HandleFunc("GET /calendar/{name}/events/add", server.ShowAddEventForm)
	http.HandleFunc("POST /calendar/{name}/events", server.AddFormEvent)
	http.HandleFunc("GET /calendar/{name}/events", server.CalendarView)

	log.Println("Listening on port", *port)
	http.ListenAndServe(*port, nil)
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
		Title:       "Test Event",
		CalendarID:  calendar.ID,
		Start:       start,
		End:         end,
		Location:    "Test Location",
		Description: "Test Description",
	}
	result = db.FirstOrCreate(&event, event)
	if result.Error != nil {
		panic(result.Error)
	}
}
