package configs

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConfDb() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env")
	}

	connStr := os.Getenv("POSTGRESQL_CONNECTION_URI")
	if connStr == "" {
		log.Fatal("POSTGRESQL_CONNECTION_URI is not set in environment variables")
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db
	log.Println("Database connection established")
}
