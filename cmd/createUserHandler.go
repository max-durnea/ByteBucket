package main
import (
	"github.com/max-durnea/ByteBucket/internal/database"
	"github.com/google/uuid"
	"database/sql"
	"encoding/json"
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
}