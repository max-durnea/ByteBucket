package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/max-durnea/ByteBucket/internal/auth"
	"github.com/max-durnea/ByteBucket/internal/database"
	"github.com/google/uuid"
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

	type req struct {
		FileName string `json:"file_name"`
		MimeType string `json:"mime_type"`
	}
	var data req
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	key := fmt.Sprintf("%s/%s", userID, data.FileName)

	presigner := s3.NewPresignClient(cfg.s3Client)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(cfg.s3Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(data.MimeType),
	}

	presignedURL, err := presigner.PresignPutObject(context.TODO(), input, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create signed URL")
		return
	}
	parsed_user_id, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create signed URL")
		return
	}
	createFileParams := database.CreateFileParams{
		ID: uuid.New(),
		UserID: parsed_user_id,
		ObjectKey: key,
		FileName: data.FileName,
		MimeType: data.MimeType,
		CreatedAt: time.Now(),
	}
	_, err = cfg.db.CreateFile(r.Context(),createFileParams)
	if err != nil {
		respondWithError(w,http.StatusInternalServerError,"Failed to upload file")
		return
	}
	// Return the URL and key to the client
	respondWithJson(w, http.StatusOK, map[string]string{
		"upload_url": presignedURL.URL,
		"key":        key,
		"mime_type":  data.MimeType,
	})

}
