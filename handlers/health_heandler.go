package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

func TablesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"table_count": count})
	}
}
