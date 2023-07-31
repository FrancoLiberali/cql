package persistence

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/persistence/gormdatabase"
)

// PersistanceModule for fx
//
// Provides:
//
// - The database connection
// - badaas-orm auto-migration
var PersistanceModule = fx.Module(
	"persistence",
	// Database connection
	fx.Provide(gormdatabase.SetupDatabaseConnection),
	// auto-migrate
	orm.AutoMigrate,
)
