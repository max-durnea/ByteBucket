package main
import (
	"github.com/max-durnea/ByteBucket/internal/database"
	"github.com/max-durnea/ByteBucket/internal/auth"
	"github.com/google/uuid"
	"encoding/json"
	"time"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request){
	type params struct{
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	data := params{}
	err := decoder.Decode(&data)
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	hashed_password,err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	createUserParams := database.CreateUserParams{
		ID: uuid.New(),
		Username: data.Username,
		Email: data.Email,
		PasswordHash: hashed_password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user,err:=cfg.db.CreateUser(r.Context(),createUserParams)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	respondWithJson(w,201,user)

}