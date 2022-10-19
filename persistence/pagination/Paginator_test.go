package pagination_test

import (
	"testing"

	"github.com/ditrit/badaas/persistence/pagination"
	"github.com/stretchr/testify/assert"
)

func TestPaginator(t *testing.T) {
	paginator := pagination.NewPaginator(uint(0), uint(12))
	assert.NotNil(t, paginator)
	assert.Equal(t, uint(12), paginator.Limit())

	paginator = pagination.NewPaginator(uint(2), uint(12))
	assert.NotNil(t, paginator)
	assert.Equal(t, uint(12), paginator.Limit())

	paginator = pagination.NewPaginator(uint(2), uint(0))
	assert.NotNil(t, paginator)
	assert.Equal(t, uint(1), paginator.Limit())
}
