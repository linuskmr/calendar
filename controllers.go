package main

import (
	"log"
	"net/http"
	"time"
)

func (s *Server) EventListController(w http.ResponseWriter, startString, endString, timezoneString string) *Events {
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

	var events []Event
	// Fetch events during the specified timestamp range (start - end).
	// Therefore, we explicitly also want to include events that are only partially within the range.
	// Consequently, the condition is either that they have started before the end or have not ended before the start.
	result := s.Db.Where("(start >= ? AND start <= ?) OR (end >= ? AND end <= ?)", start, end, start, end).Order("start").Find(&events)
	if result.Error != nil {
		log.Println("Error fetching events", result.Error)
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		return nil
	}

	return &Events{
		Start: start.In(timezone),
		End:   end.In(timezone),
		Events: events,
	}
}