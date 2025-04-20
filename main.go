package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed static/*
var staticFs embed.FS

//go:embed templates/*
var templateFs embed.FS

var templates = template.Must(template.ParseFS(templateFs, "templates/*.html"))

func main() {
	port := flag.String("addr", ":8080", "Address to listen on")
	dbFile := flag.String("db", "sqlite.db", "Database file")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open(*dbFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Event{})

	seeding(db)

	server := &Server{Db: db}

	http.HandleFunc("GET /{$}", http.RedirectHandler("/events", http.StatusFound).ServeHTTP)

	http.HandleFunc("GET /api/events", server.ApiEventList)
	http.HandleFunc("GET /api/events/{event_id}", server.ApiGetEvent)
	http.HandleFunc("POST /api/events", server.ApiAddOrUpdateEvent)
	http.HandleFunc("PATCH /api/events/{id}", server.ApiAddOrUpdateEvent)

	http.HandleFunc("GET /events/add", server.ShowAddEventForm)
	http.HandleFunc("GET /events/{event_id}", server.ShowEvent)
	http.HandleFunc("GET /events/{event_id}/edit", server.ShowEditEventForm)
	http.HandleFunc("POST /events", server.AddOrUpdateFormEvent)
	http.HandleFunc("POST /events/{id}", server.AddOrUpdateFormEvent)
	http.HandleFunc("GET /events", server.EventList)

	staticHttpPath := "/static/"
	http.HandleFunc("GET " + staticHttpPath, http.FileServerFS(staticFs).ServeHTTP)

	log.Println("Listening on port", *port)
	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func seeding(db *gorm.DB) {
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
		Start:       start,
		End:         end,
		Location:    "Test Location",
		Description: "Test Description",
	}
	result := db.FirstOrCreate(&event, event)
	if result.Error != nil {
		panic(result.Error)
	}
}
