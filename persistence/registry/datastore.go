package registry

// Describe a type of data storage
type DataStore int

const (
	_ DataStore = iota

	// A Datastore using gorm
	//
	// please see gorm.io
	GormDataStore
)
