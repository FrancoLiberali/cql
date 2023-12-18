//go:build cockroachdb
// +build cockroachdb

package testintegration

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"

	"github.com/ditrit/badaas"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/testintegration/models"
)

func TestBaDaaS(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	viper.Set("config_path", path.Join(basePath, "int_test_config.yml"))

	tGlobal = t

	badaas.BaDaaS.AddModules(
		badaas.AuthModule,
	).Provide(
		// provide test suites
		GetModels,
		repository.NewCRUD[models.Product, model.UUID],
		// create test suites
		NewCRUDRepositoryIntTestSuite,
		NewAuthServiceIntTestSuite,
	).Invoke(
		// run tests
		runTestSuites,
	).Start()
}

func runTestSuites(
	tsCRUDRepository *CRUDRepositoryIntTestSuite,
	tsAuthService *AuthServiceIntTestSuite,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsCRUDRepository)
	suite.Run(tGlobal, tsAuthService)

	shutdowner.Shutdown()
}
