package server

import (
	"embed"
	"net/http"

	"github.com/IAmSurajBobade/sb-dashboard/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func Start(content embed.FS) {
	// r := chi.()

	// r.HandleFunc("/events/get", handlers.ListHandler(content)).Methods(http.MethodGet)

	// // http.ListenAndServe(":8080", mux)

	// // // Serve static files
	// // r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(content))))

	// // // Define routes
	// // r.HandleFunc("/", handlers.HomeHandler(content)).Methods(http.MethodGet)
	// // r.HandleFunc("/create", handlers.CreatePageHandler(content)).Methods(http.MethodGet)
	// // r.HandleFunc("/create", handlers.CreateHandler).Methods(http.MethodPost)
	// // r.HandleFunc("/list", handlers.ListHandler(content)).Methods(http.MethodGet)

	// // http.Handle("/", r)

	// http.ListenAndServe(":8080", r)

	// chi router
	r := chi.NewRouter()

	r.HandleFunc("/", handlers.ListHandler(content))

	http.ListenAndServe(":8080", r)

	// mux router
	// r := mux.NewRouter()
}
