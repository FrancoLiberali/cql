package orm

import (
	"fmt"
	"net"
	"strconv"
)

func CreatePostgreSQLDSN(host, username, password, sslmode, dbname string, port int) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d sslmode=%s dbname=%s",
		username, password, host, port, sslmode, dbname,
	)
}

func CreateMySQLDSN(host, username, password, dbname string, port int) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, net.JoinHostPort(host, strconv.Itoa(port)), dbname,
	)
}

func CreateSQLiteDSN(path string) string {
	return fmt.Sprintf("sqlite:%s", path)
}

func CreateSQLServerDSN(host, username, password, dbname string, port int) string {
	return fmt.Sprintf(
		"sqlserver://%s:%s@%s?database=%s",
		username,
		password,
		net.JoinHostPort(host, strconv.Itoa(port)),
		dbname,
	)
}
