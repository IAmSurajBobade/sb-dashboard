package storage

import (
	"encoding/json"
	"os"

	"github.com/IAmSurajBobade/sb-dashboard/internal/models"
)

func DataModified() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return isDataModified
}

func ResetDataModified() {
	mutex.Lock()
	defer mutex.Unlock()
	isDataModified = false
}

func SaveToFile(backupFileName string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Create(backupFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(models.Backup{
		Service: func() string {
			if len(dataAll.Service) == 0 {
				return os.Getenv("SERVICE_NAME")
			}
			return dataAll.Service
		}(),
		Version: func() string {
			if len(dataAll.Version) == 0 {
				return os.Getenv("VERSION")
			}
			return dataAll.Version
		}(),
		Data: dataStore,
	})
}
