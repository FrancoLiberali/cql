package testintegration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/logger"
	"github.com/ditrit/badaas/persistence/database"
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
		// connect to db
		fx.Provide(NewDBConnection),
		fx.Provide(GetModels),
		orm.AutoMigrate,

		// create crud services for models
		orm.GetCRUDServiceModule[models.Seller](),
		orm.GetCRUDServiceModule[models.Company](),
		orm.GetCRUDServiceModule[models.Product](),
		orm.GetCRUDServiceModule[models.Sale](),
		orm.GetCRUDServiceModule[models.City](),
		orm.GetCRUDServiceModule[models.Country](),
		orm.GetCRUDServiceModule[models.Employee](),
		orm.GetCRUDServiceModule[models.Bicycle](),
		orm.GetCRUDServiceModule[models.Phone](),
		orm.GetCRUDServiceModule[models.Brand](),
		orm.GetCRUDServiceModule[models.Child](),

		// create test suites
		fx.Provide(NewCRUDRepositoryIntTestSuite),
		fx.Provide(NewWhereConditionsIntTestSuite),
		fx.Provide(NewJoinConditionsIntTestSuite),
		fx.Provide(NewPreloadConditionsIntTestSuite),
		fx.Provide(NewOperatorsIntTestSuite),

		// run tests
		fx.Invoke(runORMTestSuites),
	).Run()
}

func runORMTestSuites(
	tsCRUDRepository *CRUDRepositoryIntTestSuite,
	tsWhereConditions *WhereConditionsIntTestSuite,
	tsJoinConditions *JoinConditionsIntTestSuite,
	tsPreloadConditions *PreloadConditionsIntTestSuite,
	tsOperators *OperatorsIntTestSuite,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsCRUDRepository)
	suite.Run(tGlobal, tsWhereConditions)
	suite.Run(tGlobal, tsJoinConditions)
	suite.Run(tGlobal, tsPreloadConditions)
	suite.Run(tGlobal, tsOperators)

	shutdowner.Shutdown()
}

func NewDBConnection() (*gorm.DB, error) {
	return database.OpenWithRetry(
		postgres.Open(orm.CreateDSN(host, username, password, sslMode, dbName, port)),
		logger.Default.ToLogMode(logger.Info),
		10, time.Duration(5)*time.Second,
	)
}
