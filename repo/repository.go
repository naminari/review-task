package repo

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reviewtask/models"
	"time"

	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	rand.Seed(time.Now().UnixNano())
	return &Repository{db: db}
}
func (r *Repository) DB() *sql.DB {
	return r.db
}

// UserRepository методы для работы с пользователями
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetActiveUsersByTeam(teamID int, excludeUserIDs []int) ([]models.User, error)
	DeactivateUser(userID int) error
}

func (r *Repository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, is_active, team_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, user.Username, user.IsActive, user.TeamID).Scan(&user.ID, &user.CreatedAt)
}

func (r *Repository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, is_active, team_id, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.IsActive, &user.TeamID, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

// GetActiveUsersByTeam - получает активных пользователей команды, исключая указанных
func (r *Repository) GetActiveUsersByTeam(teamID int, excludeUserIDs []int) ([]models.User, error) {
	query := `
    SELECT id, username, is_active, team_id, created_at 
    FROM users 
    WHERE team_id = $1 AND is_active = true AND id != ALL($2)
    ORDER BY id
  `

	rows, err := r.db.Query(query, teamID, pq.Array(excludeUserIDs))
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
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

// TeamRepository методы для работы с командами
type TeamRepository interface {
	CreateTeam(team *models.Team) error
	GetTeamByID(id int) (*models.Team, error)
	GetTeamByName(name string) (*models.Team, error)
}

func (r *Repository) CreateTeam(team *models.Team) error {
	query := `INSERT INTO teams (name) VALUES ($1) RETURNING id, created_at`
	return r.db.QueryRow(query, team.Name).Scan(&team.ID, &team.CreatedAt)
}

func (r *Repository) GetTeamByID(id int) (*models.Team, error) {
	team := &models.Team{}
	query := `SELECT id, name, created_at FROM teams WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&team.ID, &team.Name, &team.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("team not found: %w", err)
		}
		return nil, fmt.Errorf("get team by id: %w", err)
	}
	return team, nil
}

func (r *Repository) GetTeamByName(name string) (*models.Team, error) {
	team := &models.Team{}
	query := `SELECT id, name, created_at FROM teams WHERE name = $1`
	err := r.db.QueryRow(query, name).Scan(&team.ID, &team.Name, &team.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("team not found: %w", err)
		}
		return nil, fmt.Errorf("get team by name: %w", err)
	}
	return team, nil
}

// PRRepository методы для работы с Pull Request'ами
type PRRepository interface {
	CreatePR(pr *models.PullRequest) error
	GetPRByID(id int) (*models.PullRequest, error)
	UpdatePR(pr *models.PullRequest) error
	GetPRsByReviewer(userID int) ([]models.PullRequest, error)
}

func (r *Repository) CreatePR(pr *models.PullRequest) error {
	query := `
    INSERT INTO pull_requests (title, author_id, status, reviewers) 
    VALUES ($1, $2, $3, $4) 
    RETURNING id, created_at, updated_at
  `
	return r.db.QueryRow(
		query,
		pr.Title,
		pr.AuthorID,
		pr.Status,
		pr.Reviewers,
	).Scan(&pr.ID, &pr.CreatedAt, &pr.UpdatedAt)
}

func (r *Repository) GetPRByID(id int) (*models.PullRequest, error) {
	pr := &models.PullRequest{}
	query := `
    SELECT id, title, author_id, status, reviewers, created_at, updated_at 
    FROM pull_requests 
    WHERE id = $1
  `
	err := r.db.QueryRow(query, id).Scan(
		&pr.ID,
		&pr.Title,
		&pr.AuthorID,
		&pr.Status,
		&pr.Reviewers,
		&pr.CreatedAt,
		&pr.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("PR not found: %w", err)
		}
		return nil, fmt.Errorf("get PR by id: %w", err)
	}
	return pr, nil
}

func (r *Repository) UpdatePR(pr *models.PullRequest) error {
	query := `
    UPDATE pull_requests 
    SET title = $1, status = $2, reviewers = $3, updated_at = NOW() 
    WHERE id = $4 
    RETURNING updated_at
  `
	return r.db.QueryRow(
		query,
		pr.Title,
		pr.Status,
		pr.Reviewers,
		pr.ID,
	).Scan(&pr.UpdatedAt)
}

// GetRandomReviewers - выбирает случайных ревьюеров из списка пользователей
func (r *Repository) GetRandomReviewers(users []models.User, count int) []int {
	if len(users) == 0 || count <= 0 {
		return []int{}
	}

	if count > len(users) {
		count = len(users)
	}

	// Перемешиваем массив
	shuffled := make([]models.User, len(users))
	copy(shuffled, users)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Берем первых count элементов
	result := make([]int, count)
	for i := 0; i < count; i++ {
		result[i] = shuffled[i].ID
	}

	return result
}
