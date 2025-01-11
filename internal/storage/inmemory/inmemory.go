package inmemory

import (
	"errors"
	"sync"
	"time"

	"github.com/IAmSurajBobade/sb-dashboard/internal/models"
)

var (
	dataStore = map[string][]models.Event{
		"suraj": {
			{
				ID:   1,
				Name: "Suraj's Birthday",
				Date: time.Date(1995, time.May, 28, 0, 0, 0, 0, time.Local),
			},
			{
				ID:   2,
				Name: "Suraj's Wedding Anniversary",
				Date: time.Date(2024, time.July, 14, 0, 0, 0, 0, time.Local),
			},
			{
				ID:   3,
				Name: "Suraj's Engagement Anniversary",
				Date: time.Date(2024, time.June, 28, 0, 0, 0, 0, time.Local),
			},
			{
				ID:   4,
				Name: "Dad's Birthday",
				Date: time.Date(1963, time.January, 2, 0, 0, 0, 0, time.Local),
			},
			{
				ID:   5,
				Name: "Mom's Birthday",
				Date: time.Date(1974, time.January, 1, 0, 0, 0, 0, time.Local),
			},
			{
				ID:   6,
				Name: "Shubhangi's Birthday",
				Date: time.Date(1999, time.December, 25, 0, 0, 0, 0, time.Local),
			},
		},
	}
	userStore = map[string]models.User{
		"suraj": {
			ID:        "suraj",
			Password:  "1234",
			CreatedAt: time.Date(2024, time.December, 29, 0, 0, 0, 0, time.Local),
		},
	}
	mutex = &sync.Mutex{}
)

const (
	dateFormat = "2006-01-02 15:04:05 MST"
)

func init() {
	// dataStore
}

func GetEvents(userID, password string) ([]models.Event, error) {

	mutex.Lock()
	user, ok := userStore[userID]
	if !ok {
		mutex.Unlock()
		return []models.Event{}, errors.New("user not found")
	}

	if user.Password != password {
		mutex.Unlock()
		return []models.Event{}, errors.New("invalid password")
	}

	events, ok := dataStore[userID]
	mutex.Unlock()

	if !ok {
		return []models.Event{}, errors.New("no events found")
	}
	return events, nil
}
