package repository_test

import (
	"testing"

	"github.com/ditrit/badaas/persistence/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewSortOption(t *testing.T) {
	sortOption := repository.NewSortOption("a", true)
	assert.Equal(t, "a", sortOption.Column())
	assert.True(t, sortOption.Desc())
}
