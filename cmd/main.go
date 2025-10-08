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
    fmt.Println("Hello, ByteBucket!")
    godotenv.Load()
    dbURL := os.Getenv("DATABASE_URL")
    db, err := sql.Open("postgres",dbURL)
    if err != nil {
        fmt.Printf("ERROR: Could not open database: %v",err)
        os.Exit(1)
    }
    dbQueries := database.New(db)
    apiCfg.db = dbQueries

    mux := http.NewServeMux()
    server := &http.Server{}
    server.Handler = mux
    
}
