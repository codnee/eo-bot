package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a temporary file for the backup
	tempFile, err := os.CreateTemp("", "eo-bot-backup-*.sqlite")
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	// Clean up the temp file when done
	defer os.Remove(tempPath)

	// Use SQLite VACUUM INTO command to create a backup via the existing connection
	err = db.Exec("VACUUM INTO ?", tempPath).Error
	if err != nil {
		log.Printf("Error creating SQLite backup: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Open the backup file for serving
	backupFile, err := os.Open(tempPath)
	if err != nil {
		log.Printf("Error opening backup file for serving: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer backupFile.Close()

	// Get file info for Content-Length
	fileInfo, err := backupFile.Stat()
	if err != nil {
		log.Printf("Error getting backup file info: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Type", "application/x-sqlite3")
	w.Header().Set("Content-Disposition", "attachment; filename=eo-bot-backup.sqlite")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Serve the backup file
	http.ServeContent(w, r, "eo-bot-backup.sqlite", time.Now(), backupFile)
}
