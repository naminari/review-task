package repo

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reviewtask/models"
	"strconv"
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

func (r *Repository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, is_active, team_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.DB.QueryRow(query, user.Username, user.IsActive, user.TeamID).Scan(&user.ID, &user.CreatedAt)
}

func (r *Repository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, is_active, team_id, created_at FROM users WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.IsActive, &user.TeamID, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *Repository) GetActiveUsersByTeam(teamID int, excludeUserIDs []int) ([]models.User, error) {
	query := `
    SELECT id, username, is_active, team_id, created_at 
    FROM users 
    WHERE team_id = $1 AND is_active = true
    ORDER BY id
  `

	rows, err := r.DB.Query(query, teamID)
	if err != nil {
		return nil, fmt.Errorf("get active users by team: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.IsActive, &user.TeamID, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		if !contains(excludeUserIDs, user.ID) {
			users = append(users, user)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (r *Repository) CreateTeam(team *models.Team) error {
	query := `INSERT INTO teams (name) VALUES ($1) RETURNING id, created_at`
	return r.DB.QueryRow(query, team.Name).Scan(&team.ID, &team.CreatedAt)
}

func (r *Repository) GetTeamByID(id int) (*models.Team, error) {
	team := &models.Team{}
	query := `SELECT id, name, created_at FROM teams WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&team.ID, &team.Name, &team.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("team not found: %w", err)
		}
		return nil, fmt.Errorf("get team by id: %w", err)
	}
	return team, nil
}

func (r *Repository) CreatePR(pr *models.PullRequest) error {
	query := `
    INSERT INTO pull_requests (title, author_id, status, reviewers) 
    VALUES ($1, $2, $3, $4) 
    RETURNING id, created_at, updated_at
  `
	reviewersStr := intSliceToString(pr.Reviewers)
	return r.DB.QueryRow(query, pr.Title, pr.AuthorID, pr.Status, reviewersStr).Scan(&pr.ID, &pr.CreatedAt, &pr.UpdatedAt)
}

func (r *Repository) GetPRByID(id int) (*models.PullRequest, error) {
	pr := &models.PullRequest{}
	query := `
    SELECT id, title, author_id, status, reviewers, created_at, updated_at 
    FROM pull_requests 
    WHERE id = $1
  `

	var reviewersStr string
	err := r.DB.QueryRow(query, id).Scan(
		&pr.ID,
		&pr.Title,
		&pr.AuthorID,
		&pr.Status,
		&reviewersStr,
		&pr.CreatedAt,
		&pr.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("PR not found: %w", err)
		}
		return nil, fmt.Errorf("get PR by id: %w", err)
	}

	pr.Reviewers = stringToIntSlice(reviewersStr)
	return pr, nil
}

func (r *Repository) UpdatePR(pr *models.PullRequest) error {
	query := `
    UPDATE pull_requests 
    SET title = $1, status = $2, reviewers = $3, updated_at = NOW() 
    WHERE id = $4 
    RETURNING updated_at
  `
	reviewersStr := intSliceToString(pr.Reviewers)
	return r.DB.QueryRow(query, pr.Title, pr.Status, reviewersStr, pr.ID).Scan(&pr.UpdatedAt)
}

// GetRandomReviewers - выбирает случайных ревьюеров
func (r *Repository) GetRandomReviewers(users []models.User, count int) []int {
	if len(users) == 0 || count <= 0 {
		return []int{}
	}

	if count > len(users) {
		count = len(users)
	}

	shuffled := make([]models.User, len(users))
	copy(shuffled, users)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	result := make([]int, count)
	for i := 0; i < count; i++ {
		result[i] = shuffled[i].ID
	}

	return result
}

func intSliceToString(arr []int) string {
	if len(arr) == 0 {
		return "{}"
	}
	strArr := make([]string, len(arr))
	for i, v := range arr {
		strArr[i] = strconv.Itoa(v)
	}
	return "{" + strings.Join(strArr, ",") + "}"
}

func stringToIntSlice(s string) []int {
	if s == "" || s == "{}" {
		return []int{}
	}
	clean := strings.Trim(s, "{}")
	if clean == "" {
		return []int{}
	}
	parts := strings.Split(clean, ",")
	result := make([]int, len(parts))
	for i, part := range parts {
		result[i], _ = strconv.Atoi(strings.TrimSpace(part))
	}
	return result
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (r *Repository) GetPRsByReviewer(userID int) ([]models.PullRequest, error) {
	query := `
    SELECT id, title, author_id, status, reviewers, created_at, updated_at 
    FROM pull_requests 
    WHERE $1 = ANY(reviewers)
    ORDER BY created_at DESC
  `

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("get PRs by reviewer: %w", err)
	}
	defer rows.Close()

	var prs []models.PullRequest
	for rows.Next() {
		var pr models.PullRequest
		var reviewersStr string

		err := rows.Scan(
			&pr.ID,
			&pr.Title,
			&pr.AuthorID,
			&pr.Status,
			&reviewersStr,
			&pr.CreatedAt,
			&pr.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan PR: %w", err)
		}

		pr.Reviewers = stringToIntSlice(reviewersStr)
		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return prs, nil
}
