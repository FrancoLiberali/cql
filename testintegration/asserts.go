package testintegration

import (
	"log"

	"github.com/stretchr/testify/suite"
	is "gotest.tools/assert/cmp"
)

func EqualList[T any](ts *suite.Suite, expectedList, actualList []T) {
	expectedLen := len(expectedList)
	equalLen := ts.Len(actualList, expectedLen)

	if equalLen {
		for i := 0; i < expectedLen; i++ {
			j := 0
			for ; j < expectedLen; j++ {
				if is.DeepEqual(
					actualList[j],
					expectedList[i],
				)().Success() {
					break
				}
			}

			if j == expectedLen {
				for _, element := range actualList {
					log.Println(element)
				}

				ts.FailNow("Lists not equal", "element %v not in list %v", expectedList[i], actualList)
			}
		}
	}
}
