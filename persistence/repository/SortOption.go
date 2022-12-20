package repository

type SortOption interface {
	Column() string
	Desc() bool
}

// SortOption constructor
func NewSortOption(column string, desc bool) SortOption {
	return &sortOption{column, desc}
}

// Sorting option for the repository
type sortOption struct {
	column string
	desc   bool
}

// return the column name to  sort on
func (sortOption *sortOption) Column() string {
	return sortOption.column
}

// return true for descending sort and false for ascending
func (sortOption *sortOption) Desc() bool {
	return sortOption.desc
}
