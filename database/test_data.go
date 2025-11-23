package database

import (
	"log"
	"reviewtask/models"
	"reviewtask/repo"
)

func InitTestData(repo *repo.Repository) {
	// Создаем тестовую команду backend
	backendTeam := &models.Team{
		TeamName: "backend",
		Members: []models.TeamMember{
			{UserID: "u1", Username: "alice", IsActive: true},
			{UserID: "u2", Username: "bob", IsActive: true},
			{UserID: "u3", Username: "charlie", IsActive: true},
			{UserID: "u4", Username: "diana", IsActive: true},
		},
	}

	if err := repo.CreateTeam(backendTeam); err != nil {
		log.Printf("Backend team might already exist: %v", err)
	}

	// Создаем тестовую команду frontend
	frontendTeam := &models.Team{
		TeamName: "frontend",
		Members: []models.TeamMember{
			{UserID: "u5", Username: "eve", IsActive: true},
			{UserID: "u6", Username: "frank", IsActive: true},
		},
	}

	if err := repo.CreateTeam(frontendTeam); err != nil {
		log.Printf("Frontend team might already exist: %v", err)
	}

	log.Println("Test data initialized with string IDs")
}
