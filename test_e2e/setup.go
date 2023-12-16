package main

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/testintegration"
)

func CleanDB(db *gorm.DB) {
	testintegration.CleanDBTables(db,
		[]any{
			models.Session{},
			models.User{},
		},
	)
}
