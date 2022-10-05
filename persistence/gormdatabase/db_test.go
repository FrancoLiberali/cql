package gormdatabase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDsnFromconf(t *testing.T) {
	assert.NotEmpty(t, createDsnFromConf(), "no dsn should be empty")
}

func TestCreateDsn(t *testing.T) {
	assert.Equal(t,
		"user=username password=password host=192.168.2.5 port=1225 sslmode=disable dbname=badaas_db",
		createDsn(
			"192.168.2.5",
			"username",
			"password",
			"disable",
			"badaas_db",
			1225,
		),
		"no dsn should be empty",
	)
}
