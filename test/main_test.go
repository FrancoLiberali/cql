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

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/logger"
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
	case condition.SQLite:
		dialector = sqlite.Open(cql.CreateSQLiteDSN(host))
	case condition.MySQL:
		dialector = mysql.Open(cql.CreateMySQLDSN(host, username, password, dbName, port))
	case condition.SQLServer:
		dialector = sqlserver.Open(cql.CreateSQLServerDSN(host, username, password, dbName, port))
	default:
		dialector = postgres.Open(cql.CreatePostgreSQLDSN(host, username, password, sslMode, dbName, port))
	}

	return OpenWithRetry(
		dialector,
		logger.Default.ToLogMode(logger.Info),
		10, time.Duration(5)*time.Second,
	)
}

func getDBDialector() condition.Dialector {
	return condition.Dialector(os.Getenv(dbTypeEnvKey))
}
