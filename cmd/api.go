package main
import (
	"github.com/max-durnea/ByteBucket/internal/database"
)
type apiConfig struct{
	db *database.Queries
	port string
}
