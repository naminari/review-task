package database

import (
	"log"
	"reviewtask/models"
	"reviewtask/repo"
)

func InitTestData(repo *repo.Repository) {
	team := &models.Team{Name: "backend-team"}
	if err := repo.CreateTeam(team); err != nil {
		log.Printf("Team might already exist: %v", err)
	}

	users := []struct {
		username string
		teamID   int
	}{
		{"alice", 1},
		{"bob", 1},
		{"charlie", 1},
		{"diana", 1},
	}

	for _, u := range users {
		user := &models.User{
			Username: u.username,
			IsActive: true,
			TeamID:   u.teamID,
		}
		if err := repo.CreateUser(user); err != nil {
			log.Printf("Failed to create user %s: %v", u.username, err)
		}
	}

	log.Println("Test data initialized")
}
