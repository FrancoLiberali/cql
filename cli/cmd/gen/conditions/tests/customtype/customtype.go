package customtype

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/FrancoLiberali/cql/orm/model"
)

type MultiString []string

func (s *MultiString) Scan(src interface{}) error {
	switch typedSrc := src.(type) {
	case string:
		*s = strings.Split(typedSrc, ",")
		return nil
	case []byte:
		str := string(typedSrc)
		*s = strings.Split(str, ",")

		return nil
	default:
		return fmt.Errorf("failed to scan multistring field - source is not a string, is %T", src)
	}
}

func (s MultiString) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	return strings.Join(s, ","), nil
}

func (MultiString) GormDataType() string {
	return "text"
}

type CustomType struct {
	model.UUIDModel

	Custom MultiString
}
