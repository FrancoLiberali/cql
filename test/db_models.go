package test

import (
	"log"

	"github.com/elliotchance/pie/v2"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/models"
)

var ListOfTables = []any{
	models.Product{},
	models.ProductNoTimestamps{},
	models.Company{},
	models.CompanyNoTimestamps{},
	models.Seller{},
	models.SellerNoTimestamps{},
	models.Sale{},
	models.SaleNoTimestamps{},
	models.Country{},
	models.City{},
	models.Employee{},
	models.Person{},
	models.Bicycle{},
	models.Brand{},
	models.Phone{},
	models.PhoneNoTimestamps{},
	models.ParentParent{},
	models.Parent1{},
	models.Parent2{},
	models.Child{},
}

func CleanDB(db *cql.DB) {
	CleanDBTables(db, pie.Reverse(ListOfTables))
}

func CleanDBTables(db *cql.DB, listOfTables []any) {
	// clean database to ensure independency between tests
	for _, table := range listOfTables {
		err := db.GormDB.Unscoped().Where("1 = 1").Delete(table).Error
		if err != nil {
			log.Fatalln("could not clean database: ", err)
		}
	}
}
