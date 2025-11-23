package main

import (
	"log"
	"reviewtask/database"
	"reviewtask/handlers"

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

	// Health checks
	r.GET("/health", handlers.HealthHandler(app.Repo.DB))
	r.GET("/tables", handlers.TablesHandler(app.Repo.DB))

	// Teams endpoints
	r.POST("/team/add", app.CreateTeamHandler)
	r.GET("/team/get", app.GetTeamHandler)

	// Users endpoints
	r.POST("/users/setIsActive", app.SetUserActiveHandler)
	r.GET("/users/getReview", app.GetUserReviewHandler)

	// Pull Request endpoints
	r.POST("/pullRequest/create", app.CreatePRHandler)
	r.POST("/pullRequest/merge", app.MergePRHandler)
	r.POST("/pullRequest/reassign", app.ReassignReviewerHandler)

	return r
}
