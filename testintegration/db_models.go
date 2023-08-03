package testintegration

import (
	"log"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/models"
)

var ListOfTables = []any{
	models.Product{},
	models.Company{},
	models.Seller{},
	models.Sale{},
	models.Country{},
	models.City{},
	models.Employee{},
	models.Person{},
	models.Bicycle{},
	models.Brand{},
	models.Phone{},
	models.ParentParent{},
	models.Parent1{},
	models.Parent2{},
	models.Child{},
}

func GetModels() orm.GetModelsResult {
	return orm.GetModelsResult{
		Models: ListOfTables,
	}
}

func CleanDB(db *gorm.DB) {
	CleanDBTables(db, pie.Reverse(ListOfTables))
}

func CleanDBTables(db *gorm.DB, listOfTables []any) {
	// clean database to ensure independency between tests
	for _, table := range listOfTables {
		err := db.Unscoped().Where("1 = 1").Delete(table).Error
		if err != nil {
			log.Fatalln("could not clean database: ", err)
		}
	}
}
