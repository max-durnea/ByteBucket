package main

import(
	"net/http"
	"encoding/json"
	_"github.com/max-durnea/ByteBucket/internal/database"
	"github.com/max-durnea/ByteBucket/internal/auth"
	"fmt"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("GOOD"))

}