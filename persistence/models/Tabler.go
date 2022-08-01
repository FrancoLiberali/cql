package models

var ListOfTables = []any{
	User{},
	Session{},
}

// The interface "type" need to implement to be considered models
type Tabler interface {
	// pluralized name
	TableName() string
}
