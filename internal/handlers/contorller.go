package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/IAmSurajBobade/sb-dashboard/internal/models"
	"github.com/IAmSurajBobade/sb-dashboard/internal/storage"
	"github.com/gorilla/mux"
)

type Controller struct {
	mutex     *sync.Mutex
	userCount int
	location  *time.Location
}

func NewController(location *time.Location) *Controller {

	return &Controller{
		mutex:     &sync.Mutex{},
		userCount: 0,
		location:  location,
	}
}

func (c *Controller) ListHandler(content embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]
		listName := vars["list_name"]

		events, logEvent, err := storage.GetEvents(userID, listName)
		if err != nil {
			http.Error(w, "Error retrieving events", http.StatusInternalServerError)
			return
		}

		if logEvent {
			log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, r.UserAgent())
		}

		currTime := c.now()

		for i := range events {
			events[i].EventDateStr = events[i].EventDate.Format("02/Jan/2006")
			// age in years
			events[i].AgeInDays = int(currTime.Sub(events[i].EventDate).Hours() / 24)
			y, m, d := diff(events[i].EventDate, currTime)

			events[i].AgeInYears = setAgeInYears(y, m, d)
		}

		data := map[string]interface{}{
			"UserID":   userID,
			"ListName": listName,
			"Events":   events,
		}

		tmpl, _ := template.ParseFS(content, "templates/user/list.html")
		tmpl.Execute(w, data)
	}
}

func setAgeInYears(y, m, d int) string {
	age := ""
	if y == 1 {
		age = "1 year"
	} else if y > 1 {
		age = fmt.Sprintf("%d years", y)
	}
	if m == 1 {
		age = fmt.Sprintf("%s 1 month", age)
	} else if m > 1 {
		age = fmt.Sprintf("%s %d months", age, m)
	}
	if d == 1 {
		age = fmt.Sprintf("%s 1 day", age)
	} else if d > 1 {
		age = fmt.Sprintf("%s %d days", age, d)
	}
	return age
}

func (c *Controller) EditHandler(content embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		eventID, _ := strconv.Atoi(vars["id"])
		userID := vars["user_id"]
		listName := vars["list_name"]

		event, err := storage.GetEventByID(eventID)
		if err != nil {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		data := map[string]interface{}{
			"ID":           event.ID,
			"Title":        event.Title,
			"EventDateStr": event.EventDate.Format("2006-01-02"),
			"UserID":       userID,
			"ListName":     listName,
		}

		tmpl, _ := template.ParseFS(content, "templates/user/edit_event.html")
		tmpl.Execute(w, data)
	}
}

func (c *Controller) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eventID, _ := strconv.Atoi(vars["id"])
	userID := vars["user_id"]
	listName := vars["list_name"]

	event, err := storage.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	storage.DeleteEvent(userID, listName, event)

	w.Header().Set("HX-Redirect", fmt.Sprintf("/events/users/%s/lists/%s", userID, listName))
}

func (c *Controller) SaveHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	title := r.FormValue("title")
	dateStr := r.FormValue("event_date")
	eventDate, err := time.Parse(time.RFC3339, dateStr+"T00:00:00+05:30")
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	event := models.Event{
		ID:        id,
		Title:     title,
		EventDate: eventDate,
	}

	userID := r.FormValue("user_id")
	listName := r.FormValue("list_name")

	if event.ID == 0 {
		event.ID = storage.SaveEvent(userID, listName, event)
	} else {
		storage.UpdateEvent(userID, listName, event)
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/events/users/%s/lists/%s", userID, listName))
}

func diff(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour := int(h2 - h1)
	min := int(m2 - m1)
	sec := int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, a.Location())
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
