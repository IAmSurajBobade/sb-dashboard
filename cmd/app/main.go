package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/IAmSurajBobade/sb-dashboard/internal/handlers"
	"github.com/IAmSurajBobade/sb-dashboard/internal/storage"
	"github.com/gorilla/mux"
)

//go:embed templates/*
var content embed.FS

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8100"
	}

	muxRouter := mux.NewRouter()
	location := time.UTC
	if l, err := time.LoadLocation("Asia/Kolkata"); err == nil {
		location = l
	}

	ctrl := handlers.NewController(location)

	// Custom 404 handler
	muxRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// redirect to home page
		fmt.Printf("404: %s %s\n", r.Method, r.URL.Path)
		http.Redirect(w, r, "/events/", http.StatusSeeOther)
	})
	muxRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
	muxRouter.Use(recoveryMiddleware)
	// muxRouter.Use(loggingMiddleware)

	eventsRouter := muxRouter.PathPrefix("/events").Subrouter()

	eventsRouter.HandleFunc("/", ctrl.HomeHandler(content)).Methods(http.MethodGet)
	// eventsRouter.HandleFunc("/users/", ctrl.UsersHandler(content)).Methods(http.MethodGet)
	// eventsRouter.HandleFunc("/users/{user_id}/lists", ctrl.GetListHandler(content)).Methods(http.MethodGet)
	eventsRouter.HandleFunc("/users/{user_id}/lists/{list_name}", ctrl.ListHandler(content)).Methods(http.MethodGet)
	eventsRouter.HandleFunc("/users/{user_id}/lists/{list_name}/edit/{id}", ctrl.EditHandler(content)).Methods(http.MethodGet)
	eventsRouter.HandleFunc("/users/{user_id}/lists/{list_name}/delete/{id}", ctrl.DeleteHandler).Methods(http.MethodGet)
	eventsRouter.HandleFunc("/save", ctrl.SaveHandler).Methods(http.MethodPost)
	// Serve static files
	//muxRouter.PathPrefix("/css/").Handler(http.FileServer(http.FS(content)))

	// Create a channel to listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		fmt.Println("Server started on port:", port)
		if err := http.ListenAndServe(":"+port, muxRouter); err != nil {
			fmt.Println("Server stopped with error:", err)
		}
	}()

	// Start a goroutine to perform periodic backups
	go periodicBackup(location)

	// Wait for a termination signal
	sig := <-sigChan
	fmt.Println("Received signal:", sig)

	if err := backup("crash", location); err != nil {
		fmt.Println("Error saving backup:", err)
	}

	fmt.Println("Server stopped gracefully")
}

// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		next.ServeHTTP(w, r)
// 		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
// 	})
// }

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Recovered from panic: %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Start a goroutine to perform periodic backups
func periodicBackup(l *time.Location) {
	for {
		time.Sleep(1 * time.Hour * 24) // Adjust the interval as needed
		if storage.DataModified() {
			if err := backup("periodic", l); err == nil { // no error
				storage.ResetDataModified()
			}
		}
	}
}

func backup(backupType string, l *time.Location) error {
	backupDir := "backup"
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		fmt.Println("Error creating backup directory:", err)
		return err
	}

	backupFileName := filepath.Join(backupDir, fmt.Sprintf("%s_%s.json", backupType, time.Now().In(l).Format("20060102150405")))
	if err := storage.SaveToFile(backupFileName); err != nil {
		return fmt.Errorf("error saving backup: %v", err.Error())
	} else {
		fmt.Println("Backup saved to:", backupFileName)
		return nil
	}
}
