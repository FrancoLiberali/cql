package services

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
)

var ServicesModule = fx.Module(
	"services",
	fx.Provide(userservice.NewUserService),
	fx.Provide(sessionservice.NewSessionService),
)
