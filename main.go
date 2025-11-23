package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		fmt.Println("Loading .env file...")
		// В продакшн лучше использовать github.com/joho/godotenv
		// godotenv.Load()
	}

	db, err := sql.Open("postgres", getDBConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := waitForDB(db); err != nil {
		log.Fatal("Database not ready:", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		var result int
		err := db.QueryRow("SELECT 1").Scan(&result)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "DOWN",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"db":     "connected",
		})
	})

	r.GET("/tables", func(c *gin.Context) {
		rows, err := db.Query(`
      SELECT table_name 
      FROM information_schema.tables 
      WHERE table_schema = 'public'
    `)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			tables = append(tables, table)
		}

		c.JSON(http.StatusOK, gin.H{"tables": tables})
	})

	port := getEnv("APP_PORT", "8080")
	log.Printf("Server starting on :%s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getDBConnectionString() string {
	host := getEnv("DB_HOST", "localhost")
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
	for i := 0; i < 30; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}
		log.Printf("Waiting for database... (attempt %d/30)", i+1)
	}
	return fmt.Errorf("database connection timeout")
}
