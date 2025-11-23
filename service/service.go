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

func (s *ReviewService) AssignReviewers(authorID int, teamID int) ([]int, error) {
	excludeIDs := []int{authorID}

	availableUsers, err := s.repo.GetActiveUsersByTeam(teamID, excludeIDs)
	if err != nil {
		return nil, fmt.Errorf("get available reviewers: %w", err)
	}

	reviewers := s.repo.GetRandomReviewers(availableUsers, 2)

	return reviewers, nil
}

func (s *ReviewService) ValidatePRForUpdate(pr *models.PullRequest) error {
	if pr.Status == models.StatusMerged {
		return fmt.Errorf("cannot update merged PR")
	}
	return nil
}
