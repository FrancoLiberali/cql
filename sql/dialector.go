package sql

type Dialector string

const (
	Postgres  Dialector = "postgres"
	MySQL     Dialector = "mysql"
	SQLite    Dialector = "sqlite"
	SQLServer Dialector = "sqlserver"
)
