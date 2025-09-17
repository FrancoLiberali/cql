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

	"github.com/FrancoLiberali/cql/logger"
	"github.com/FrancoLiberali/cql/sql"
)

const dbTypeEnvKey = "DB"

const (
	username = "cql"
	password = "cql_password2023"
	host     = "localhost"
	port     = 5001
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
	suite.Run(t, NewGroupByIntTestSuite(db))
	suite.Run(t, NewSelectIntTestSuite(db))
	suite.Run(t, NewFunctionsIntTestSuite(db))
	suite.Run(t, NewInsertIntTestSuite(db))
}

func NewDBConnection() (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch getDBDialector() {
	case sql.Postgres:
		dialector = postgres.Open(
			fmt.Sprintf(
				"user=%s password=%s host=%s port=%d sslmode=%s dbname=%s",
				username, password, host, port, sslMode, dbName,
			),
		)
	case sql.MySQL:
		dialector = mysql.Open(
			fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				username, password, net.JoinHostPort(host, strconv.Itoa(port)), dbName,
			),
		)
	case sql.SQLServer:
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
		dialector = sqlite.Open("sqlite:" + host)
	}

	return OpenWithRetry(
		dialector,
		logger.Default.ToLogMode(logger.Info),
		10, time.Duration(5)*time.Second,
	)
}

func getDBDialector() sql.Dialector {
	dialector := os.Getenv(dbTypeEnvKey)
	if dialector != "" {
		return sql.Dialector(os.Getenv(dbTypeEnvKey))
	}

	return sql.SQLite
}
