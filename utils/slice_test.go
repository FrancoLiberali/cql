package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/utils"
)

var (
	testResult3 = 33.04
	testResult4 = 0.11
)

var findFirstTests = []struct {
	ss         []float64
	expression func(value float64) bool
	expected   *float64
}{
	{
		nil,
		func(value float64) bool { return value == 1.5 },
		nil,
	},
	{
		[]float64{},
		func(value float64) bool { return value == 0.1 },
		nil,
	},
	{
		[]float64{0.0, 1.5, 3.2},
		func(value float64) bool { return value == 9.99 },
		nil,
	},
	{
		[]float64{5.4, 6.98, 4.987, 33.04},
		func(value float64) bool { return value == 33.04 },
		&testResult3,
	},
	{
		[]float64{9.0, 0.11, 150.44, 33.04},
		func(value float64) bool { return value == 0.11 },
		&testResult4,
	},
}

func TestFindFirst(t *testing.T) {
	for _, test := range findFirstTests {
		t.Run("", func(t *testing.T) {
			result := utils.FindFirst(test.ss, test.expression)
			if result == nil {
				assert.Nil(t, test.expected)
			} else {
				assert.Equal(t, *test.expected, *result)
			}
		})
	}
}
