package basicauth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
)

func TestSaltAndHashPassword(t *testing.T) {
	password := "password"
	hash := basicauth.SaltAndHashPassword(password)
	assert.NotEqual(t, string(hash), password, "the password and it's hash shouln't be equals")
}

func TestCheckUserPassword(t *testing.T) {
	password := "voila"
	hash := basicauth.SaltAndHashPassword(password)
	assert.True(t, basicauth.CheckUserPassword(hash, password), "the password and it's hash should match")
	assert.False(t, basicauth.CheckUserPassword(hash, "wrong password"), "the password and it's hash should match")
}
