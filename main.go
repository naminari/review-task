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

	r.GET("/health", handlers.HealthHandler(app.Repo.DB))
	r.GET("/tables", handlers.TablesHandler(app.Repo.DB))

	api := r.Group("/api/v1")
	{
		api.POST("/pull-requests", app.CreatePRHandler)
		api.GET("/pull-requests/:id", app.GetPRHandler)
		api.POST("/pull-requests/:id/merge", app.MergePRHandler)
	}

	return r
}
