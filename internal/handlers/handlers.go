package handlers

import (
	"embed"
	"html/template"
	"net/http"
	"time"

	"github.com/IAmSurajBobade/sb-dashboard/internal/storage/inmemory"
)

// func HomeHandler(content embed.FS) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tmpl, _ := template.ParseFS(content, "templates/home.html")
// 		tmpl.Execute(w, nil)
// 	}
// }

// func CreatePageHandler(content embed.FS) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tmpl, _ := template.ParseFS(content, "templates/create.html")
// 		tmpl.Execute(w, nil)
// 	}
// }

// func CreateHandler(w http.ResponseWriter, r *http.Request) {
// 	var event models.Event
// 	err := json.NewDecoder(r.Body).Decode(&event)
// 	if err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	if len(event.Password) != 4 {
// 		http.Error(w, "Password must be 4 digits", http.StatusBadRequest)
// 		return
// 	}

// 	event.ID = storage.SaveEvent(event)
// 	response := map[string]string{"id": event.ID}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

func ListHandler(content embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identifier := r.FormValue("identifier")
		password := r.FormValue("password")

		events, err := inmemory.GetEvents(identifier, password)
		if err != nil {
			http.Error(w, "Invalid identifier or password", http.StatusUnauthorized)
			return
		}

		for i := range events {
			events[i].DaysSince = int(time.Since(events[i].Date).Hours() / 24)
		}

		tmpl, _ := template.ParseFS(content, "templates/list.html")
		tmpl.Execute(w, map[string]interface{}{"Events": events})
	}
}
