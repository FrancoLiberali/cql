package utils

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestIntToSecond(t *testing.T) {
	assert.Equal(
		t,
		IntToSecond(20),
		time.Duration(20*time.Second),
		"the duration should be equals",
	)
	assert.Equal(
		t,
		IntToSecond(-5),
		time.Duration(-5*time.Second),
		"the duration should be equals",
	)
	assert.Equal(
		t,
		IntToSecond(3600),
		time.Duration(time.Hour),
		"the duration should be equals",
	)
}
