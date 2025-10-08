package main

import(
	"net/http"
	"encoding/json"
	_"github.com/max-durnea/ByteBucket/internal/database"
	"github.com/max-durnea/ByteBucket/internal/auth"
	"github.com/google/uuid"
	"fmt"
	"time"
)



func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request){
	type params struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	data := params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	userDb, err := cfg.db.GetUserByEmail(r.Context(),data.Email)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	err = auth.CheckPasswordHash(data.Password,userDb.PasswordHash)
	if err != nil {
		respondWithError(w,401,fmt.Sprintf("%v",err))
		return
	}
	//For now simply respond with OK
	jwtToken,err := auth.MakeJWT(userDb.ID,cfg.tokenSecret,15*time.Minute)
	if err != nil {
		respondWithError(w,401,fmt.Sprintf("%v",err))
		return
	}
	refreshToken := auth.MakeRefreshToken();
	refreshTokenParams := CreateRefreshTokenParams{
		Token: refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: userDb.ID,
		ExpiresAt: 24 * time.Hour,
	}
	_,err := cfg.db.CreateRefreshToken(r.Context(),refreshTokenParams)
	response:= struct{
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Username string `json:"username"`
		JWTtoken string `json:"jwt_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID: userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email: userDb.Email,
		Username: userDb.Username,
		JWTtoken: jwtToken,
		RefreshToken: refreshToken,
	}
	respondWithJson(w,200,response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

}