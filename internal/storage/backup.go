package storage

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
