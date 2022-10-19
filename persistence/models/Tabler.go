package models

// The interface "type" need to implement to be considered models
type Tabler interface {
	// pluralized name
	TableName() string
}
