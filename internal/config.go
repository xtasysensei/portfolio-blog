package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Connstr, dbname string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("POSTGRES_SERVER")
	dbname = os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	Connstr = fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
}
func ConnectDB(w http.ResponseWriter) *sql.DB {
	db, err := sql.Open("postgres", Connstr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to database!")
	w.WriteHeader(http.StatusOK)
	return db
}
