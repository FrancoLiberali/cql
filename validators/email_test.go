package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	validator "github.com/ditrit/badaas/validators"
)

func TestValidEmail(t *testing.T) {
	mail, err := validator.ValidEmail("bob.bobemail.com")
	assert.Error(t, err)
	assert.Equal(t, "", mail)

	mail, err = validator.ValidEmail("bob.bob@")
	assert.Error(t, err)
	assert.Equal(t, "", mail)

	mail, err = validator.ValidEmail("bob.bob@email.com")
	assert.NoError(t, err)
	assert.Equal(t, "bob.bob@email.com", mail)

	mail, err = validator.ValidEmail("Gopher <from@example.com>")
	assert.NoError(t, err)
	assert.Equal(t, "from@example.com", mail)

	mail, err = validator.ValidEmail("bob.bob%@email.com")
	assert.NoError(t, err)
	assert.Equal(t, "bob.bob%@email.com", mail)
}
