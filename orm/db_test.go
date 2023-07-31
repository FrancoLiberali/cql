package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDSN(t *testing.T) {
	assert.Equal(
		t,
		"user=username password=password host=192.168.2.5 port=1225 sslmode=disable dbname=badaas_db",
		CreateDSN(
			"192.168.2.5",
			"username",
			"password",
			"disable",
			"badaas_db",
			1225,
		),
	)
}
