package models

import "time"

type (
	// User Events
	User struct {
		ID        string    `json:"id"`
		Password  string    `json:"password"`
		Events    []Event   `json:"events"` // events identifier
		CreatedAt time.Time `json:"created_at"`
	}
	Event struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		Date      time.Time `json:"date"`
		DaysSince int       `json:"days_since"`
	}

	EventGetReq struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}
	EventGetResp struct {
		Events []Event `json:"events"`
	}
)
