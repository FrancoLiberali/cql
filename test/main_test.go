package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/cql"
	"github.com/ditrit/badaas/orm/logger"
)

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
	db, err := NewDBConnection()
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(ListOfTables...)
	if err != nil {
		log.Fatalln(err)
	}

	suite.Run(t, NewQueryIntTestSuite(db))
	suite.Run(t, NewWhereConditionsIntTestSuite(db))
	suite.Run(t, NewJoinConditionsIntTestSuite(db))
	suite.Run(t, NewPreloadConditionsIntTestSuite(db))
	suite.Run(t, NewOperatorsIntTestSuite(db))
	suite.Run(t, NewUpdateIntTestSuite(db))
	suite.Run(t, NewDeleteIntTestSuite(db))
}

func NewDBConnection() (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch getDBDialector() {
	case cql.SQLite:
		dialector = sqlite.Open(orm.CreateSQLiteDSN(host))
	case cql.MySQL:
		dialector = mysql.Open(orm.CreateMySQLDSN(host, username, password, dbName, port))
	case cql.SQLServer:
		dialector = sqlserver.Open(orm.CreateSQLServerDSN(host, username, password, dbName, port))
	default:
		dialector = postgres.Open(orm.CreatePostgreSQLDSN(host, username, password, sslMode, dbName, port))
	}

	return OpenWithRetry(
		dialector,
		logger.Default.ToLogMode(logger.Info),
		10, time.Duration(5)*time.Second,
	)
}

func getDBDialector() cql.Dialector {
	return cql.Dialector(os.Getenv(dbTypeEnvKey))
}
