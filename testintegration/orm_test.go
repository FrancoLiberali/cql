package testintegration

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/models"
)

var tGlobal *testing.T

const (
	username = "root"
	password = "postgres"
	host     = "localhost"
	port     = 26257
	sslMode  = "disable"
	dbName   = "badaas_db"
)

func TestBaDaaSORM(t *testing.T) {
	tGlobal = t

	fx.New(
		// logger
		fx.Provide(NewLoggerConfiguration),
		logger.LoggerModule,

		// connect to db
		fx.Provide(NewGormDBConnection),

		// activate badaas-orm
		fx.Provide(GetModels),
		orm.AutoMigrate,

		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		// create crud services for models
		orm.GetCRUDServiceModule[models.Seller](),
		orm.GetCRUDServiceModule[models.Product](),
		orm.GetCRUDServiceModule[models.Sale](),
		orm.GetCRUDServiceModule[models.City](),
		orm.GetCRUDServiceModule[models.Country](),
		orm.GetCRUDServiceModule[models.Employee](),
		orm.GetCRUDServiceModule[models.Bicycle](),
		orm.GetCRUDServiceModule[models.Phone](),
		orm.GetCRUDServiceModule[models.Brand](),

		// create test suites
		fx.Provide(NewCRUDRepositoryIntTestSuite),
		fx.Provide(NewWhereConditionsIntTestSuite),
		fx.Provide(NewJoinConditionsIntTestSuite),
		fx.Provide(NewOperatorsIntTestSuite),

		// run tests
		fx.Invoke(runORMTestSuites),
	).Run()
}

func runORMTestSuites(
	tsCRUDRepository *CRUDRepositoryIntTestSuite,
	tsWhereConditions *WhereConditionsIntTestSuite,
	tsJoinConditions *JoinConditionsIntTestSuite,
	tsOperators *OperatorsIntTestSuite,
	db *gorm.DB,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsCRUDRepository)
	suite.Run(tGlobal, tsWhereConditions)
	suite.Run(tGlobal, tsJoinConditions)
	suite.Run(tGlobal, tsOperators)

	shutdowner.Shutdown()
}

func NewLoggerConfiguration() configuration.LoggerConfiguration {
	viper.Set(configuration.LoggerModeKey, "dev")
	return configuration.NewLoggerConfiguration()
}

func NewGormDBConnection(logger *zap.Logger) (*gorm.DB, error) {
	return orm.ConnectToDialector(
		logger,
		orm.CreateDialector(host, username, password, sslMode, dbName, port),
		10, time.Duration(5)*time.Second,
	)
}
