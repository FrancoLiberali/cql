package configuration

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestIntToSecond(t *testing.T) {
	assert.Equal(
		t,
		intToSecond(20),
		time.Duration(20*time.Second),
		"the duration should be equals",
	)
	assert.Equal(
		t,
		intToSecond(-5),
		time.Duration(-5*time.Second),
		"the duration should be equals",
	)
	assert.Equal(
		t,
		intToSecond(3600),
		time.Duration(time.Hour),
		"the duration should be equals",
	)
}
