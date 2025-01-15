package models

import "time"

type (
	Event struct {
		ID           int       `json:"id"`
		Title        string    `json:"title"`
		EventDate    time.Time `json:"event_date"`
		EventDateStr string    `json:"event_date_str"`
		AgeInDays    int       `json:"-"`
		AgeInYears   string    `json:"-"`
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
