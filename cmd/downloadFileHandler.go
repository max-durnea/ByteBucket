package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"time"
	"github.com/max-durnea/ByteBucket/internal/auth"
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

		// Ensure requester owns the file
		rawUser := r.Context().Value(auth.UserIDKey)
		userIDStr, ok := rawUser.(string)
		if !ok || userIDStr == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}
		if fileDb.UserID != userID {
			// Do not reveal existence - return 404
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		
		// Generate a pre-signed URL for downloading
		presigner := s3.NewPresignClient(cfg.s3Client)
		
		input := &s3.GetObjectInput{
			Bucket: aws.String(cfg.s3Bucket),
			Key:    aws.String(fileDb.ObjectKey),
		}
		presignedURL, err := presigner.PresignGetObject(r.Context(), input, s3.WithPresignExpires(15*time.Minute))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create signed URL: %v", err), http.StatusInternalServerError)
			return
		}
		
		signedURL := presignedURL.URL
		// Respond with the signed URL
		respondWithJson(w, http.StatusOK, map[string]string{"url": signedURL})
		return
	}

	// If prefix does not match
	http.Error(w, "Invalid URL", http.StatusBadRequest)
}
