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

func (s *ReviewService) AssignReviewers(authorID string) ([]string, error) {
	author, err := s.repo.GetUser(authorID)
	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	excludeIDs := []string{authorID}

	availableUsers, err := s.repo.GetActiveUsersByTeam(author.TeamName, excludeIDs)
	if err != nil {
		return nil, fmt.Errorf("get available reviewers: %w", err)
	}

	reviewers := s.repo.GetRandomReviewers(availableUsers, 2)

	return reviewers, nil
}

func (s *ReviewService) CreatePRWithReviewers(prID, prName, authorID string) (*models.PullRequest, error) {
	existingPR, _ := s.repo.GetPR(prID)
	if existingPR != nil {
		return nil, fmt.Errorf("PR id already exists")
	}

	_, err := s.repo.GetUser(authorID)
	if err != nil {
		return nil, fmt.Errorf("author not found")
	}

	reviewers, err := s.AssignReviewers(authorID)
	if err != nil {
		return nil, err
	}

	pr := &models.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            models.StatusOpen,
		AssignedReviewers: reviewers,
	}

	if err := s.repo.CreatePR(pr); err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	return pr, nil
}

func (s *ReviewService) MergePR(prID string) (*models.PullRequest, error) {
	pr, err := s.repo.GetPR(prID)
	if err != nil {
		return nil, fmt.Errorf("PR not found")
	}

	if pr.Status == models.StatusMerged {
		return pr, nil
	}

	if err := s.repo.MergePR(prID); err != nil {
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}

	return s.repo.GetPR(prID)
}

func (s *ReviewService) GetUserReviewPRs(userID string) ([]models.PullRequestShort, error) {
	if _, err := s.repo.GetUser(userID); err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return s.repo.GetPRsByReviewer(userID)
}

func (s *ReviewService) ReassignReviewer(pullRequestID string, oldUserID string) (string, error) {
	pr, err := s.repo.GetPR(pullRequestID)
	if err != nil {
		return "", fmt.Errorf("PR not found: %w", err)
	}

	if pr.Status == models.StatusMerged {
		return "", fmt.Errorf("cannot reassign on merged PR")
	}

	if !containsString(pr.AssignedReviewers, oldUserID) {
		return "", fmt.Errorf("reviewer is not assigned to this PR")
	}

	oldReviewer, err := s.repo.GetUser(oldUserID)
	if err != nil {
		return "", fmt.Errorf("old reviewer not found")
	}

	availableUsers, err := s.repo.GetActiveUsersByTeam(oldReviewer.TeamName, pr.AssignedReviewers)
	if err != nil {
		return "", fmt.Errorf("get available candidates: %w", err)
	}

	if len(availableUsers) == 0 {
		return "", fmt.Errorf("no active replacement candidate in team")
	}

	newReviewer := s.repo.GetRandomReviewers(availableUsers, 1)[0]

	for i, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldUserID {
			pr.AssignedReviewers[i] = newReviewer
			break
		}
	}

	if err := s.repo.UpdatePRReviewers(pr); err != nil {
		return "", fmt.Errorf("failed to update PR: %w", err)
	}

	return newReviewer, nil
}

// Team management methods
func (s *ReviewService) CreateTeam(team *models.Team) error {
	existingTeam, _ := s.repo.GetTeam(team.TeamName)
	if existingTeam != nil {
		return fmt.Errorf("team_name already exists")
	}

	return s.repo.CreateTeam(team)
}

func (s *ReviewService) SetUserActive(userID string, isActive bool) (*models.User, error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := s.repo.SetUserActive(userID, isActive); err != nil {
		return nil, err
	}

	user.IsActive = isActive
	return user, nil
}

// Вспомогательная функция
func containsString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
