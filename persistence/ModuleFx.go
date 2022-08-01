package persistence

import (
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

// PersistanceModule for fx
//
// Provides:
//
// - The database connection
//
// - The repositories
var PersistanceModule = fx.Module(
	"persistence",
	// Database connection
	fx.Provide(gormdatabase.CreateDatabaseConnectionFromConfiguration),

	//repositories
	fx.Provide(repository.NewCRUDRepository[models.Session, uuid.UUID]),
	fx.Provide(repository.NewCRUDRepository[models.User, uuid.UUID]),
)
