package main

import (
	"log"
	"reviewtask/database"
	"reviewtask/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация БД
	db := database.InitDB()
	defer db.DB().Close()

	// Инициализация репозиториев и сервисов
	app := handlers.NewApp(db)

	// Тестовые данные
	database.InitTestData(app.Repo)

	// Роутер
	r := setupRouter(app)

	// Запуск
	port := "8080"
	log.Printf("Server starting on :%s...", port)
	log.Fatal(r.Run(":" + port))
}

func setupRouter(app *handlers.App) *gin.Engine {
	r := gin.Default()

	// Health checks
	r.GET("/health", handlers.HealthHandler(app.Repo.DB()))
	r.GET("/tables", handlers.TablesHandler(app.Repo.DB()))

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/pull-requests", app.CreatePRHandler)
		api.GET("/pull-requests/:id", app.GetPRHandler)
		api.POST("/pull-requests/:id/merge", app.MergePRHandler)
	}

	return r
}
