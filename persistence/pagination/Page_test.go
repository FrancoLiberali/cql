package pagination_test

import (
	"testing"

	"github.com/ditrit/badaas/persistence/pagination"
	"github.com/stretchr/testify/assert"
)

type Whatever struct {
	DumbData int
}

func (Whatever) TableName() string {
	return "whatevers"
}

var (
	// test fixture
	ressources = []*Whatever{
		{10},
		{11},
		{12},
		{13},
		{14},
		{15},
		{16},
		{17},
		{18},
		{19},
	}
)

func TestNewPage(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		1,  // page 1
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.ElementsMatch(t, ressources, p.Ressources)
	assert.Equal(t, uint(10), p.Limit)
	assert.Equal(t, uint(1), p.Offset)
	assert.Equal(t, uint(5), p.TotalPages)
}

func TestPageHasNextPageFalse(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		4,  // page 4: last page
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.False(t, p.HasNextPage)
}

func TestPageHasNextPageTrue(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		1,  // page 1
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.True(t, p.HasNextPage)
}

func TestPageIsLastPageFalse(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		1,  // page 1
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.False(t, p.IsLastPage)
}

func TestPageIsLastPageTrue(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		4,  // page 4: last page
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.True(t, p.IsLastPage)
}

func TestPageHasPreviousPageFalse(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		0,  // page 1
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.False(t, p.HasPreviousPage)
}

func TestPageHasPreviousPageTrue(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		1,  // page 1
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.True(t, p.HasPreviousPage)
}

func TestPageIsFirstPageFalse(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		1,  // page 1
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.False(t, p.IsFirstPage)
}

func TestPageIsFirstPageTrue(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		0,  // page 0: first page
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.True(t, p.IsFirstPage)
}

func TestPageHasContentFalse(t *testing.T) {
	p := pagination.NewPage(
		[]*Whatever{}, // no content
		0,             // page 1
		10,            // 10 elems per page
		50,            // 50 elem in total
	)
	assert.False(t, p.HasPreviousPage)
}

func TestPageHasContentTrue(t *testing.T) {
	p := pagination.NewPage(
		ressources,
		1,  // page 1
		10, // 10 elems per page
		50, // 50 elem in total
	)
	assert.True(t, p.HasContent)
}
