package registry

import (
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
)

// The list of database table to migrate if necessary
var listOfDatabaseTables = []any{
	&models.User{},
}

// Create a registry using gorm
func createGormRegistry() (*Registry, error) {
	gormDatabase, err := gormdatabase.InitializeDBFromConf()
	if err != nil {
		return nil, err
	}
	err = gormdatabase.AutoMigrate(gormDatabase, listOfDatabaseTables...)
	if err != nil {
		return nil, err
	}
	return &Registry{
			UserRepository: repository.NewCRUDRepository[models.User](gormDatabase),
		},
		nil
}
