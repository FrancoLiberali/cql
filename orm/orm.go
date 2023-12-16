package orm

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

func GetCRUD[T Model, ID ModelID](db *gorm.DB) (CRUDService[T, ID], CRUDRepository[T, ID]) {
	repository := NewCRUDRepository[T, ID]()
	return NewCRUDService(db, repository), repository
}

func autoMigrate(modelsLists [][]any, db *gorm.DB) error {
	if len(modelsLists) > 0 {
		allModels := pie.Flat(modelsLists)
		return db.AutoMigrate(allModels...)
	}

	return nil
}
