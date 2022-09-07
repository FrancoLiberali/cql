package configuration

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestIntToSecond(t *testing.T) {
	testCases := []struct {
		desc     string
		seconds  int
		duration time.Duration
	}{
		{
			desc:     "20 secondes",
			seconds:  20,
			duration: time.Duration(20 * time.Second),
		}, {
			desc:     "-5 seconds",
			seconds:  -5,
			duration: time.Duration(-5 * time.Second),
		}, {
			desc:     "3600	seconds",
			seconds:  3600,
			duration: time.Duration(time.Hour),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(
				t,
				intToSecond(tC.seconds),
				tC.duration,
				"the duration should be equals",
			)
		})
	}
}
