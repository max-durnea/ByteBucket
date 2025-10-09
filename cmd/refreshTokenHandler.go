package main

import (
	"encoding/json"
	"github.com/max-durnea/ByteBucket/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Token string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	var data params
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w, 400, "Could not decode refresh token")
		return
	}
	tokenDB, err := cfg.db.GetRefreshToken(r.Context(), data.Token)
	if err != nil {
		respondWithError(w, 400, "Bad Token")
		return
	}
	if time.Now().After(tokenDB.ExpiresAt) {
		respondWithError(w, 400, "Token Expired")
		return
	}
	if tokenDB.RevokedAt.Valid {
		respondWithError(w, 400, "Token revoked")
		return
	}
	//If refresh token is valid, extract the user_id from it and create a jwt_token based on this
	userDB, err := cfg.db.GetUserFromRefreshToken(r.Context(), data.Token)
	if err != nil {
		respondWithError(w, 400, "User does not exist")
		return
	}
	//create the jwt
	jwt_token, err := auth.MakeJWT(userDB.ID, cfg.tokenSecret, 15*time.Minute)
	if err != nil {
		respondWithError(w, 400, "Could not generate JWT token")
		return
	}
	respondWithJson(w, 200, struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: jwt_token,
	})

}
