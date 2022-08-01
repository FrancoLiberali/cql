package pagination

// Handle pagination
type Paginator interface {
	Offset() uint
	Limit() uint
}

type paginatorImpl struct {
	offset uint
	limit  uint
}

// Constructor of Paginator
func NewPaginator(page, limit uint) Paginator {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 1
	}
	return &paginatorImpl{
		offset: page,
		limit:  limit,
	}
}

// Return the page number
func (p *paginatorImpl) Offset() uint {
	return p.offset
}

// Return the max number of records for one page
func (p *paginatorImpl) Limit() uint {
	return p.limit
}
