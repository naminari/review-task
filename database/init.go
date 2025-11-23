package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"reviewtask/repo"

	_ "github.com/lib/pq"
)

func InitDB() *repo.Repository {
	db, err := sql.Open("postgres", getDBConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := waitForDB(db); err != nil {
		log.Fatal("Database not ready:", err)
	}

	return repo.NewRepository(db)
}

func getDBConnectionString() string {
	host := getEnv("DB_HOST", "db")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "review_service")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func waitForDB(db *sql.DB) error {
	for i := 0; i < 10; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}
		log.Printf("Waiting for database... (attempt %d/10)", i+1)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("database connection timeout")
}
