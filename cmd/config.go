package main

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/max-durnea/ByteBucket/internal/database"
)

type apiConfig struct {
	db          *database.Queries
	port        string
	tokenSecret string
	platform    string
	s3Client    *s3.Client
	s3Region    string
	s3Bucket    string
}
