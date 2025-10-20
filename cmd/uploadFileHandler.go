package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/max-durnea/ByteBucket/internal/auth"
	"github.com/max-durnea/ByteBucket/internal/database"
)

// user_id can be taken from the context
func (cfg *apiConfig) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	rawUser := r.Context().Value(auth.UserIDKey)
	userID, ok := rawUser.(string)
	if !ok || userID == "" {
		log.Printf("uploadFileHandler: userID missing or wrong type in context: %#v\n", rawUser)
		respondWithError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	// Expect a multipart/form-data request with a "file" field
	// Parse up to 32MB of memory, rest will be stored in tmp files by the server
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not parse multipart form")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Missing file field")
		return
	}
	defer file.Close()

	fileName := header.Filename
	// create temp file
	tmpPath := ""
	tmpFile, err := os.CreateTemp(os.TempDir(), "bb-upload-*")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create temp file")
		return
	}
	tmpPath = tmpFile.Name()
	defer func() {
		tmpFile.Close()
		// best-effort cleanup
		os.Remove(tmpPath)
	}()

	// copy uploaded content to temp file
	if _, err := io.Copy(tmpFile, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save uploaded file")
		return
	}

	// reopen temp file for reading
	tmpRead, err := os.Open(tmpPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to open temp file")
		return
	}
	defer tmpRead.Close()

	// determine mime type: prefer the header's content-type if present
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	key := fmt.Sprintf("%s/%s", userID, filepath.Base(fileName))

	// upload to S3 from server side
	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(cfg.s3Bucket),
		Key:         aws.String(key),
		Body:        tmpRead,
		ContentType: aws.String(mimeType),
	}
	_, err = cfg.s3Client.PutObject(context.TODO(), putInput)
	if err != nil {
		log.Printf("s3 PutObject error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to upload to S3")
		return
	}

	parsed_user_id, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid user id")
		return
	}
	createFileParams := database.CreateFileParams{
		ID:        uuid.New(),
		UserID:    parsed_user_id,
		ObjectKey: key,
		FileName:  fileName,
		MimeType:  mimeType,
		CreatedAt: time.Now(),
	}
	_, err = cfg.db.CreateFile(r.Context(), createFileParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to record file in database")
		return
	}

	// Return the key and success to the client (no presigned URLs)
	respondWithJson(w, http.StatusOK, map[string]string{
		"key":       key,
		"file_name": fileName,
		"mime_type": mimeType,
	})

}
