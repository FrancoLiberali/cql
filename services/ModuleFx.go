package services

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
)

var AuthServiceModule = fx.Module(
	"authService",
	// models
	fx.Provide(getAuthModels),
	// repositories
	fx.Provide(orm.NewCRUDRepository[models.Session, orm.UUID]),
	fx.Provide(orm.NewCRUDRepository[models.User, orm.UUID]),

	// services
	fx.Provide(userservice.NewUserService),
	fx.Provide(sessionservice.NewSessionService),
)

func getAuthModels() orm.GetModelsResult {
	return orm.GetModelsResult{
		Models: []any{
			models.Session{},
			models.User{},
		},
	}
}
