package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestModels(t *testing.T) {
	t.Run("Team model", func(t *testing.T) {
		team := Team{
			TeamName: "test-team",
			Members: []TeamMember{
				{UserID: "u1", Username: "User1", IsActive: true},
				{UserID: "u2", Username: "User2", IsActive: false},
			},
		}

		assert.Equal(t, "test-team", team.TeamName)
		assert.Len(t, team.Members, 2)
		assert.Equal(t, "u1", team.Members[0].UserID)
		assert.False(t, team.Members[1].IsActive)
	})

	t.Run("PullRequest model", func(t *testing.T) {
		now := time.Now()
		pr := PullRequest{
			PullRequestID:     "pr-001",
			PullRequestName:   "Test PR",
			AuthorID:          "author1",
			Status:            StatusOpen,
			AssignedReviewers: []string{"reviewer1", "reviewer2"},
			CreatedAt:         now,
		}

		assert.Equal(t, "pr-001", pr.PullRequestID)
		assert.Equal(t, StatusOpen, pr.Status)
		assert.Len(t, pr.AssignedReviewers, 2)
		assert.Equal(t, now, pr.CreatedAt)
	})

	t.Run("PRStatus constants", func(t *testing.T) {
		assert.Equal(t, PRStatus("OPEN"), StatusOpen)
		assert.Equal(t, PRStatus("MERGED"), StatusMerged)
	})
}
