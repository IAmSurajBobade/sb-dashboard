package main

import (
	"embed"
	"fmt"
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
	muxRouter := mux.NewRouter()

	ctrl := handlers.NewController()

	muxRouter.HandleFunc("/user/{user_id}/list/{list_name}", ctrl.ListHandler(content)).Methods("GET")
	muxRouter.HandleFunc("/user/{user_id}/list/{list_name}/edit/{id}", ctrl.EditHandler(content)).Methods(http.MethodGet)
	muxRouter.HandleFunc("/save", ctrl.SaveHandler).Methods("POST")

	// Create a channel to listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := http.ListenAndServe(":8080", muxRouter); err != nil {
			fmt.Println("Server stopped with error:", err)
		}
	}()

	// Start a goroutine to perform periodic backups
	go periodicBackup()

	// Wait for a termination signal
	sig := <-sigChan
	fmt.Println("Received signal:", sig)

	if err := backup("crash"); err != nil {
		fmt.Println("Error saving backup:", err)
	}

	fmt.Println("Server stopped gracefully")
}

// Start a goroutine to perform periodic backups
func periodicBackup() {
	for {
		time.Sleep(1 * time.Hour * 24) // Adjust the interval as needed
		if storage.DataModified() {
			if err := backup("periodic"); err == nil { // no error
				storage.ResetDataModified()
			}
		}
	}
}

func backup(backupType string) error {
	backupDir := "backup"
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		fmt.Println("Error creating backup directory:", err)
		return err
	}

	backupFileName := filepath.Join(backupDir, fmt.Sprintf("%s_%s.json", backupType, time.Now().Format("20060102150405")))
	if err := storage.SaveToFile(backupFileName); err != nil {
		return fmt.Errorf("error saving backup: %v", err.Error())
	} else {
		fmt.Println("Backup saved to:", backupFileName)
		return nil
	}
}
