package service

import (
	"fmt"
	"reviewtask/models"
	"reviewtask/repo"
)

type ReviewService struct {
	repo *repo.Repository
}

func NewReviewService(repo *repo.Repository) *ReviewService {
	return &ReviewService{repo: repo}
}

// AssignReviewers - основная логика назначения ревьюеров
func (s *ReviewService) AssignReviewers(authorID int, teamID int) ([]int, error) {
	// Исключаем автора из списка возможных ревьюеров
	excludeIDs := []int{authorID}

	// Получаем активных пользователей команды
	availableUsers, err := s.repo.GetActiveUsersByTeam(teamID, excludeIDs)
	if err != nil {
		return nil, fmt.Errorf("get available reviewers: %w", err)
	}

	// Выбираем случайных ревьюеров (до 2)
	reviewers := s.repo.GetRandomReviewers(availableUsers, 2)

	return reviewers, nil
}

// ValidatePRForUpdate - проверяет можно ли обновлять PR
func (s *ReviewService) ValidatePRForUpdate(pr *models.PullRequest) error {
	if pr.Status == models.StatusMerged {
		return fmt.Errorf("cannot update merged PR")
	}
	return nil
}
