package test

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/logger"
)

const dbTypeEnvKey = "DB"

const (
	username = "cql"
	password = "cql_password2023"
	host     = "localhost"
	port     = 5000
	sslMode  = "disable"
	dbName   = "cql_db"
)

func TestCQL(t *testing.T) {
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
		dialector = sqlite.Open(fmt.Sprintf("sqlite:%s", host))
	case condition.MySQL:
		dialector = mysql.Open(
			fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				username, password, net.JoinHostPort(host, strconv.Itoa(port)), dbName,
			),
		)
	case condition.SQLServer:
		dialector = sqlserver.Open(
			fmt.Sprintf(
				"sqlserver://%s:%s@%s?database=%s",
				username,
				password,
				net.JoinHostPort(host, strconv.Itoa(port)),
				dbName,
			),
		)
	default:
		dialector = postgres.Open(
			fmt.Sprintf(
				"user=%s password=%s host=%s port=%d sslmode=%s dbname=%s",
				username, password, host, port, sslMode, dbName,
			),
		)
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
