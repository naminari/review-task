package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryHelpers(t *testing.T) {
	t.Run("containsString helper", func(t *testing.T) {
		assert.True(t, containsString([]string{"a", "b", "c"}, "b"))
		assert.False(t, containsString([]string{"a", "b", "c"}, "d"))
		assert.False(t, containsString([]string{}, "a"))
	})

	t.Run("Repository creation", func(t *testing.T) {
		// Просто проверяем что структура создается
		repo := &Repository{}
		assert.NotNil(t, repo)
	})
}
