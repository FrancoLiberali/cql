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
		20*time.Second,
		"the duration should be equals",
	)
	assert.Equal(
		t,
		IntToSecond(-5),
		-5*time.Second,
		"the duration should be equals",
	)
	assert.Equal(
		t,
		IntToSecond(3600),
		time.Hour,
		"the duration should be equals",
	)
}
