package repo

import (
	"database/sql"

	"math/rand"
	"reviewtask/models"
	"strings"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	rand.Seed(time.Now().UnixNano())
	return &Repository{DB: db}
}

// Team methods
func (r *Repository) CreateTeam(team *models.Team) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO teams (team_name) VALUES ($1)", team.TeamName)
	if err != nil {
		return err
	}

	for _, member := range team.Members {
		_, err = tx.Exec(
			"INSERT INTO users (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4)",
			member.UserID, member.Username, team.TeamName, member.IsActive,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetTeam(teamName string) (*models.Team, error) {
	team := &models.Team{TeamName: teamName}

	rows, err := r.DB.Query(
		"SELECT user_id, username, is_active FROM users WHERE team_name = $1",
		teamName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member models.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, err
		}
		team.Members = append(team.Members, member)
	}

	return team, nil
}

// User methods
func (r *Repository) SetUserActive(userID string, isActive bool) error {
	_, err := r.DB.Exec(
		"UPDATE users SET is_active = $1 WHERE user_id = $2",
		isActive, userID,
	)
	return err
}

func (r *Repository) GetUser(userID string) (*models.User, error) {
	user := &models.User{}
	err := r.DB.QueryRow(
		"SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1",
		userID,
	).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// PR methods - обновляем для работы со строковыми ID
func (r *Repository) CreatePR(pr *models.PullRequest) error {
	query := `
    INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers) 
    VALUES ($1, $2, $3, $4, $5) 
    RETURNING created_at
  `
	reviewersStr := strings.Join(pr.AssignedReviewers, ",")
	return r.DB.QueryRow(
		query,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
		reviewersStr,
	).Scan(&pr.CreatedAt)
}

func (r *Repository) GetPR(pullRequestID string) (*models.PullRequest, error) {
	pr := &models.PullRequest{}
	var reviewersStr string
	var mergedAt sql.NullTime

	query := `
    SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
    FROM pull_requests 
    WHERE pull_request_id = $1
  `

	err := r.DB.QueryRow(query, pullRequestID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&reviewersStr,
		&pr.CreatedAt,
		&mergedAt,
	)
	if err != nil {
		return nil, err
	}

	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	pr.AssignedReviewers = strings.Split(reviewersStr, ",")
	if len(pr.AssignedReviewers) == 1 && pr.AssignedReviewers[0] == "" {
		pr.AssignedReviewers = []string{}
	}

	return pr, nil
}

func (r *Repository) MergePR(pullRequestID string) error {
	_, err := r.DB.Exec(
		"UPDATE pull_requests SET status = 'MERGED', merged_at = NOW() WHERE pull_request_id = $1",
		pullRequestID,
	)
	return err
}

func (r *Repository) UpdatePRReviewers(pr *models.PullRequest) error {
	reviewersStr := strings.Join(pr.AssignedReviewers, ",")
	_, err := r.DB.Exec(
		"UPDATE pull_requests SET assigned_reviewers = $1 WHERE pull_request_id = $2",
		reviewersStr, pr.PullRequestID,
	)
	return err
}

func (r *Repository) GetPRsByReviewer(userID string) ([]models.PullRequestShort, error) {
	query := `
    SELECT pull_request_id, pull_request_name, author_id, status
    FROM pull_requests 
    WHERE $1 = ANY(STRING_TO_ARRAY(assigned_reviewers, ','))
    ORDER BY created_at DESC
  `

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

// Business logic
func (r *Repository) GetActiveUsersByTeam(teamName string, excludeUserIDs []string) ([]models.User, error) {
	query := `
    SELECT user_id, username, team_name, is_active 
    FROM users 
    WHERE team_name = $1 AND is_active = true
    ORDER BY user_id
  `

	rows, err := r.DB.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, err
		}

		// Фильтруем исключенных пользователей
		if !containsString(excludeUserIDs, user.UserID) {
			users = append(users, user)
		}
	}

	return users, nil
}

func (r *Repository) GetRandomReviewers(users []models.User, count int) []string {
	if len(users) == 0 || count <= 0 {
		return []string{}
	}

	if count > len(users) {
		count = len(users)
	}

	shuffled := make([]models.User, len(users))
	copy(shuffled, users)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = shuffled[i].UserID
	}

	return result
}

// Вспомогательные функции
func containsString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
