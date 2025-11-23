package main

import (
	"log"
	"reviewtask/database"
	"reviewtask/handlers"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.InitDB()
	defer db.DB.Close()

	app := handlers.NewApp(db)

	database.InitTestData(app.Repo)

	r := setupRouter(app)

	port := "8080"
	log.Printf("Server starting on :%s...", port)
	log.Fatal(r.Run(":" + port))
}

func setupRouter(app *handlers.App) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.HealthHandler(app.Repo.DB))
	r.GET("/tables", handlers.TablesHandler(app.Repo.DB))

	api := r.Group("/api/v1")
	{
		api.POST("/pull-requests", app.CreatePRHandler)
		api.GET("/pull-requests/:id", app.GetPRHandler)
		api.POST("/pull-requests/:id/merge", app.MergePRHandler)
		api.GET("/users/:id/pull-requests", app.GetPRsByUserHandler)

		api.PUT("/pull-requests/:id/reviewers", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))

			var req struct {
				OldReviewerID int `json:"old_reviewer_id"`
				NewReviewerID int `json:"new_reviewer_id"`
			}

			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "bad request"})
				return
			}

			pr, err := app.Repo.GetPRByID(id)
			if err != nil {
				c.JSON(404, gin.H{"error": "pr not found"})
				return
			}

			if pr.Status == "MERGED" {
				c.JSON(400, gin.H{"error": "cannot change reviewers on merged PR"})
				return
			}

			found := false
			for i, reviewerID := range pr.Reviewers {
				if reviewerID == req.OldReviewerID {
					pr.Reviewers[i] = req.NewReviewerID
					found = true
					break
				}
			}

			if !found {
				c.JSON(400, gin.H{"error": "old reviewer not found in PR"})
				return
			}

			if err := app.Repo.UpdatePR(pr); err != nil {
				c.JSON(500, gin.H{"error": "failed to update PR"})
				return
			}

			c.JSON(200, pr)
		})

	}

	return r
}
