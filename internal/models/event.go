package models

import "time"

type (
	Event struct {
		ID        int       `json:"id"`
		Title     string    `json:"title"`
		EventDate time.Time `json:"event_date"`
		DaysSince int       `json:"days_since"`
	}

	EventList struct {
		ID           string             `json:"id"`
		PasswordHash string             `json:"password_hash"`
		Lists        map[string][]Event `json:"lists"`
	}

	User struct {
		ID           string             `json:"id"`
		PasswordHash string             `json:"password_hash"`
		Lists        map[string][]Event `json:"lists"`
	}

	Backup struct {
		Service string `json:"service"`
		Version string `json:"version"`
		Data    []User `json:"data"`
	}
)
