package main

import (
    _ "github.com/lib/pq"
    "database/sql"
    "github.com/max-durnea/ByteBucket/internal/database"
    "fmt"
    "net/http"
    "github.com/joho/godotenv"
    "os"
)

var apiCfg = apiConfig{}

func main() {
    
    godotenv.Load()
    dbURL := os.Getenv("DATABASE_URL")

    db, err := sql.Open("postgres",dbURL)
    if err != nil {
        fmt.Printf("ERROR: Could not open database: %v",err)
        os.Exit(1)
    }
    defer db.Close()
    dbQueries := database.New(db)
    apiCfg.db = dbQueries
    apiCfg.port = os.Getenv("PORT")
    mux := http.NewServeMux()
    server := &http.Server{}

    server.Handler = mux
    server.Addr = fmt.Sprintf("localhost:%v",apiCfg.port)
    fmt.Printf("Server running on port %v\n",apiCfg.port)

    mux.HandleFunc("POST /api/users",apiCfg.createUserHandler)
    server.ListenAndServe()
}
