package services

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/persistence/gormfx"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
)

var AuthServiceModule = fx.Module(
	"authService",
	// models
	fx.Provide(getAuthModels),
	// repositories
	fx.Provide(repository.NewCRUD[models.Session, model.UUID]),
	fx.Provide(repository.NewCRUD[models.User, model.UUID]),

	// services
	fx.Provide(userservice.NewUserService),
	fx.Provide(sessionservice.NewSessionService),
)

func getAuthModels() gormfx.GetModelsResult {
	return gormfx.GetModelsResult{
		Models: []any{
			models.Session{},
			models.User{},
		},
	}
}
