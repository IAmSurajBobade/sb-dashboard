package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/IAmSurajBobade/sb-dashboard/internal/models"
)

var (
	dataAll        models.Backup
	dataStore      []models.User
	mutex          = &sync.Mutex{}
	isDataModified bool
)

func init() {
	loadInitialData()
}

func loadInitialData() {
	backupDir := "backup"
	files, err := os.ReadDir(backupDir)
	if err != nil {
		fmt.Println("Error reading backup directory:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No backup files found")
		return
	}

	// Sort files by name in descending order to get the latest file first
	sort.Slice(files, func(i, j int) bool {
		file1, _ := files[i].Info()
		file2, _ := files[j].Info()
		return file1.ModTime().After(file2.ModTime())
	})

	latestFile := filepath.Join(backupDir, files[0].Name())
	fmt.Println("Loading initial data from:", latestFile)
	file, err := os.Open(latestFile)

	if err != nil {
		fmt.Println("Error loading initial data:", err)
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&dataAll); err != nil {
		fmt.Println("Error decoding initial data:", err)
		return
	}
	dataStore = dataAll.Data
}

func SaveEvent(userID, listName string, event models.Event) int {
	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range dataStore {
		if user.ID == userID {
			event.ID = len(user.Lists[listName]) + 1
			dataStore[i].Lists[listName] = append(dataStore[i].Lists[listName], event)
			isDataModified = true
			return event.ID
		}
	}
	// If user does not exist, create a new user
	newUser := models.User{
		ID:           userID,
		PasswordHash: "", // Set an appropriate default or hash
		Lists:        make(map[string][]models.Event),
	}
	event.ID = 1
	newUser.Lists[listName] = []models.Event{event}
	dataStore = append(dataStore, newUser)
	isDataModified = true
	return event.ID
}

func UpdateEvent(userID, listName string, event models.Event) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range dataStore {
		if user.ID == userID {
			for j, e := range user.Lists[listName] {
				if e.ID == event.ID {
					dataStore[i].Lists[listName][j] = event
					isDataModified = true
					return
				}
			}
		}
	}
}

func GetEventByID(eventID int) (models.Event, error) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, user := range dataStore {
		for _, events := range user.Lists {
			for _, event := range events {
				if event.ID == eventID {
					return event, nil
				}
			}
		}
	}
	return models.Event{}, fmt.Errorf("event not found")
}

func GetEvents(userID, listName string) ([]models.Event, error) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range dataStore {
		if user.ID == userID {
			events, exists := user.Lists[listName]
			if !exists {
				dataStore[i].Lists[listName] = []models.Event{}
				return dataStore[i].Lists[listName], nil

			}
			return events, nil
		}
	}
	// If user does not exist, create a new user with an empty list
	newUser := models.User{
		ID:           userID,
		PasswordHash: "", // Set an appropriate default or hash
		Lists:        make(map[string][]models.Event),
	}
	newUser.Lists[listName] = []models.Event{}
	dataStore = append(dataStore, newUser)
	return newUser.Lists[listName], nil
}

func SaveToFile(filePath string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(models.Backup{
		Service: dataAll.Service,
		Version: dataAll.Version,
		Data:    dataStore,
	})
}
