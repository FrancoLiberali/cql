package testintegration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/logger"
	"github.com/ditrit/badaas/orm/query"
	"github.com/ditrit/badaas/persistence/database"
	"github.com/ditrit/badaas/persistence/gormfx"
)

var tGlobal *testing.T

const dbTypeEnvKey = "DB"

const (
	username = "badaas"
	password = "badaas_password2023"
	host     = "localhost"
	port     = 5000
	sslMode  = "disable"
	dbName   = "badaas_db"
)

func TestBaDaaSORM(t *testing.T) {
	tGlobal = t

	fx.New(
		// connect to db
		fx.Provide(NewDBConnection),
		fx.Provide(GetModels),
		gormfx.AutoMigrate,

		// create test suites
		fx.Provide(NewQueryIntTestSuite),
		fx.Provide(NewWhereConditionsIntTestSuite),
		fx.Provide(NewJoinConditionsIntTestSuite),
		fx.Provide(NewPreloadConditionsIntTestSuite),
		fx.Provide(NewOperatorsIntTestSuite),

		// run tests
		fx.Invoke(runORMTestSuites),
	).Run()
}

func runORMTestSuites(
	tsQuery *QueryIntTestSuite,
	tsWhereConditions *WhereConditionsIntTestSuite,
	tsJoinConditions *JoinConditionsIntTestSuite,
	tsPreloadConditions *PreloadConditionsIntTestSuite,
	tsOperators *OperatorsIntTestSuite,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsQuery)
	suite.Run(tGlobal, tsWhereConditions)
	suite.Run(tGlobal, tsJoinConditions)
	suite.Run(tGlobal, tsPreloadConditions)
	suite.Run(tGlobal, tsOperators)

	shutdowner.Shutdown()
}

func NewDBConnection() (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch getDBDialector() {
	case query.Postgres:
		dialector = postgres.Open(orm.CreatePostgreSQLDSN(host, username, password, sslMode, dbName, port))
	case query.SQLite:
		dialector = sqlite.Open(orm.CreateSQLiteDSN(host))
	case query.MySQL:
		dialector = mysql.Open(orm.CreateMySQLDSN(host, username, password, dbName, port))
	case query.SQLServer:
		dialector = sqlserver.Open(orm.CreateSQLServerDSN(host, username, password, dbName, port))
	default:
		return nil, fmt.Errorf("unknown db %s", getDBDialector())
	}

	return database.OpenWithRetry(
		dialector,
		logger.Default.ToLogMode(logger.Info),
		10, time.Duration(5)*time.Second,
	)
}

func getDBDialector() query.Dialector {
	return query.Dialector(os.Getenv(dbTypeEnvKey))
}
