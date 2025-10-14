package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/max-durnea/ByteBucket/internal/database"
	"net/http"
	"os"
	"time"
)

var apiCfg = apiConfig{}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("ERROR: Could not open database: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)
	apiCfg.db = dbQueries
	apiCfg.port = os.Getenv("PORT")
	apiCfg.tokenSecret = os.Getenv("TOKEN_SECRET")
	apiCfg.platform = os.Getenv("PLATFORM")
	apiCfg.s3Bucket = os.Getenv("S3_BUCKET")
	apiCfg.s3Region = os.Getenv("AWS_REGION")

	cfgAWS, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(apiCfg.s3Region),
	)

	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(cfgAWS)
	apiCfg.s3Client = s3Client

	mux := http.NewServeMux()
	server := &http.Server{}

	server.Handler = mux
	server.Addr = fmt.Sprintf("localhost:%v", apiCfg.port)
	fmt.Printf("Server running on port %v\n", apiCfg.port)

	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginUserHandler)
	mux.Handle("POST /api/files", apiCfg.JwtMiddleware(http.HandlerFunc(apiCfg.uploadFileHandler)))
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshTokenHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.Handle("GET /api/files/",apiCfg.JwtMiddleware(http.HandlerFunc(apiCfg.downloadFileHandler)))
	apiCfg.StartRefreshTokenCleanup(15 * time.Minute)
	server.ListenAndServe()
}
