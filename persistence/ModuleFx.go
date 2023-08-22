package persistence

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/persistence/database"
	"github.com/ditrit/badaas/persistence/gormfx"
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
	fx.Provide(database.SetupDatabaseConnection),
	// auto-migrate
	gormfx.AutoMigrate,
)
