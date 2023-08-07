package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/utils/validators"
)

func TestValidEmail(t *testing.T) {
	mail, err := validators.ValidEmail("bob.bobemail.com")
	assert.Error(t, err)
	assert.Equal(t, "", mail)

	mail, err = validators.ValidEmail("bob.bob@")
	assert.Error(t, err)
	assert.Equal(t, "", mail)

	mail, err = validators.ValidEmail("bob.bob@email.com")
	assert.NoError(t, err)
	assert.Equal(t, "bob.bob@email.com", mail)

	mail, err = validators.ValidEmail("Gopher <from@example.com>")
	assert.NoError(t, err)
	assert.Equal(t, "from@example.com", mail)

	mail, err = validators.ValidEmail("bob.bob%@email.com")
	assert.NoError(t, err)
	assert.Equal(t, "bob.bob%@email.com", mail)
}
