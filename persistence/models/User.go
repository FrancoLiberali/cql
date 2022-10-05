package models

// Represents a user
type User struct {
	BaseModel
	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`

	// password hash
	Password []byte `gorm:"not null"`
}

// Return the pluralized table name
//
// Satisfie the Tabler interface
func (User) TableName() string {
	return "users"
}
