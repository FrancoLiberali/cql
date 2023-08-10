package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePostgreDSN(t *testing.T) {
	assert.Equal(
		t,
		"user=username password=password host=192.168.2.5 port=1225 sslmode=disable dbname=badaas_db",
		CreatePostgreSQLDSN(
			"192.168.2.5",
			"username",
			"password",
			"disable",
			"badaas_db",
			1225,
		),
	)
}

func TestCreateMySQLDSN(t *testing.T) {
	assert.Equal(
		t,
		"username:password@tcp(192.168.2.5:1225)/badaas_db?charset=utf8mb4&parseTime=True&loc=Local",
		CreateMySQLDSN(
			"192.168.2.5",
			"username",
			"password",
			"badaas_db",
			1225,
		),
	)
}

func TestCreateSQLiteDSN(t *testing.T) {
	assert.Equal(
		t,
		"sqlite:/dir/file",
		CreateSQLiteDSN(
			"/dir/file",
		),
	)
}

func TestCreateSQLServerDSN(t *testing.T) {
	assert.Equal(
		t,
		"sqlserver://username:password@192.168.2.5:1225?database=badaas_db",
		CreateSQLServerDSN(
			"192.168.2.5",
			"username",
			"password",
			"badaas_db",
			1225,
		),
	)
}
