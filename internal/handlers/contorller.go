package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/IAmSurajBobade/sb-dashboard/internal/models"
	"github.com/IAmSurajBobade/sb-dashboard/internal/storage"
	"github.com/gorilla/mux"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) ListHandler(content embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]
		listName := vars["list_name"]

		events, err := storage.GetEvents(userID, listName)
		if err != nil {
			http.Error(w, "Error retrieving events", http.StatusInternalServerError)
			return
		}

		for i := range events {
			events[i].DaysSince = int(time.Since(events[i].EventDate).Hours() / 24)
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

	w.Header().Set("HX-Redirect", fmt.Sprintf("/user/%s/list/%s", userID, listName))
}
