package orm

import (
	"fmt"
)

func methodError(err error, method string) error {
	return fmt.Errorf("%w; method: %s", err, method)
}
