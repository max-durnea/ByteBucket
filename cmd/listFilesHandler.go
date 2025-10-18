package main

import (
    "net/http"
    "github.com/max-durnea/ByteBucket/internal/auth"
    "github.com/google/uuid"
    "fmt"
)

func (cfg *apiConfig) listFilesHandler(w http.ResponseWriter, r *http.Request) {
    rawUser := r.Context().Value(auth.UserIDKey)
    userIDStr, ok := rawUser.(string)
    if !ok || userIDStr == "" {
        respondWithError(w, http.StatusUnauthorized, "User not authorized")
        return
    }
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid user in context")
        return
    }

    files, err := cfg.db.GetFilesByUser(r.Context(), userID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
        return
    }

    // map DB files to response objects
    type fileResp struct{
        ID string `json:"id"`
        FileName string `json:"file_name"`
        MimeType string `json:"mime_type"`
        Key string `json:"key"`
        CreatedAt string `json:"created_at"`
    }
    out := make([]fileResp, 0, len(files))
    for _, f := range files {
        out = append(out, fileResp{
            ID: f.ID.String(),
            FileName: f.FileName,
            MimeType: f.MimeType,
            Key: f.ObjectKey,
            CreatedAt: f.CreatedAt.String(),
        })
    }

    respondWithJson(w, http.StatusOK, out)
}
