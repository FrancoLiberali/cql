package main

import (
	"log"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/persistence/models"
)

var ListOfTables = []any{
	models.Session{},
	models.User{},
}

func CleanDB(db *gorm.DB) {
	// clean database to ensure independency between tests
	for _, table := range ListOfTables {
		err := db.Unscoped().Where("1 = 1").Delete(table).Error
		if err != nil {
			log.Fatalln("could not clean database: ", err)
		}
	}
}
