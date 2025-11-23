package handlers

import (
	"reviewtask/repo"
	"reviewtask/service"
)

type App struct {
	Repo    *repo.Repository
	Service *service.ReviewService
}

func NewApp(db *repo.Repository) *App {
	repo := db
	reviewService := service.NewReviewService(repo)
	return &App{Repo: repo, Service: reviewService}
}
