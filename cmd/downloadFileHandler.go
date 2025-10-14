package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (cfg *apiConfig) downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	prefix := "/api/files/"

	// Check if the path starts with the prefix
	if strings.HasPrefix(path, prefix) {
		fileIDStr := strings.TrimPrefix(path, prefix)

		// Optional: parse UUID if your IDs are UUIDs
		id, err := uuid.Parse(fileIDStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid file ID: %v", err), http.StatusBadRequest)
			return
		}

		// Fetch the file from DB
		fileDb, err := cfg.db.GetFileById(r.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf("File not found: %v", err), http.StatusNotFound)
			return
		}

		// Respond with JSON (or generate signed URL)
		respondWithJson(w, http.StatusOK, fileDb)
		return
	}

	// If prefix does not match
	http.Error(w, "Invalid URL", http.StatusBadRequest)
}
