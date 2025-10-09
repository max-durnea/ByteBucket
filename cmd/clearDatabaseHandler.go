package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "")
		return
	}
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("%v", err))
		return
	}
	err = cfg.db.ResetRefreshTokens(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("%v", err))
		return
	}
	w.WriteHeader(200)
}
