package gormdatabase

import (
	"testing"

	configurationmocks "github.com/ditrit/badaas/mocks/configuration"
	"github.com/stretchr/testify/assert"
)

func TestCreateDsnFromconf(t *testing.T) {
	conf := configurationmocks.NewDatabaseConfiguration(t)
	conf.On("GetPort").Return(1225)
	conf.On("GetHost").Return("192.168.2.5")
	conf.On("GetDBName").Return("badaas_db")
	conf.On("GetUsername").Return("username")
	conf.On("GetPassword").Return("password")
	conf.On("GetSSLMode").Return("disable")
	assert.Equal(t, "user=username password=password host=192.168.2.5 port=1225 sslmode=disable dbname=badaas_db",
		createDsnFromConf(conf))
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
