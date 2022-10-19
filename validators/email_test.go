package validator_test

import (
	"testing"

	validator "github.com/ditrit/badaas/validators"
	"github.com/stretchr/testify/assert"
)

func TestValidEmail(t *testing.T) {
	assert.False(t, validator.ValidEmail("bob.bobemail.com"))
	assert.False(t, validator.ValidEmail("bob.bob@"))

	assert.True(t, validator.ValidEmail("bob.bob@email.com"))
	assert.True(t, validator.ValidEmail("bob.bob%@email.com"))
}
