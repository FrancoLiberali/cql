package pagination

import "github.com/ditrit/badaas/persistence/models"

// A page hold ressources and data regarding the pagination
type Page[T models.Tabler] struct {
	Ressources []*T `json:"ressources"`
	// max d'element par page
	Limit uint `json:"limit"`
	// page courante
	Offset uint `json:"offset"`
	// total d'element en base
	Total uint `json:"total"`
	// total de pages
	TotalPages      uint `json:"totalPages"`
	HasNextPage     bool `json:"hasNextpage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	IsFirstPage     bool `json:"isFirstPage"`
	IsLastPage      bool `json:"isLastPage"`
	HasContent      bool `json:"hasContent"`
}

// Create a new page
func NewPage[T models.Tabler](records []*T, offset, size, nbElemTotal uint) *Page[T] {
	nbPage := nbElemTotal / size
	p := Page[T]{
		Ressources: records,
		Limit:      size,
		Offset:     offset,
		Total:      nbElemTotal,
		TotalPages: nbPage,

		HasNextPage:     nbElemTotal > (offset+1)*size,
		HasPreviousPage: offset >= 1,
		IsFirstPage:     offset == 0,
		IsLastPage:      offset == (nbPage - 1),
		HasContent:      len(records) != 0,
	}
	return &p
}
