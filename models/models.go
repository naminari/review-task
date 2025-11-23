package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	TeamID    int       `json:"team_id" db:"team_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Team struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

// PullRequest - сущность PR
type PullRequest struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	Status    PRStatus  `json:"status" db:"status"`
	Reviewers UserIDs   `json:"reviewers" db:"reviewers"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserIDs - кастомный тип для массива ID пользователей
type UserIDs []int

// Scan - реализация интерфейса Scanner для работы с PostgreSQL массивами
func (u *UserIDs) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, u)
}

// Value - реализация интерфейса Valuer
func (u UserIDs) Value() (driver.Value, error) {
	return json.Marshal(u)
}
