package conditions

import (
	"fmt"
	"go/types"
	"regexp"
	"strings"

	"github.com/elliotchance/pie/v2"
)

// badaas/orm/baseModels.go
var badaasORMBaseModels = []string{"github.com/ditrit/badaas/orm.UUIDModel", "github.com/ditrit/badaas/orm.UIntModel", "gorm.io/gorm.Model"}

type Type struct {
	types.Type
}

// Get the name of the type depending of the internal type
func (t Type) Name() string {
	switch typeTyped := t.Type.(type) {
	case *types.Named:
		return typeTyped.Obj().Name()
	default:
		return pie.Last(strings.Split(t.String(), "."))
	}
}

// Get the package of the type depending of the internal type
func (t Type) Pkg() *types.Package {
	switch typeTyped := t.Type.(type) {
	case *types.Named:
		return typeTyped.Obj().Pkg()
	default:
		return nil
	}
}

// Get the struct under type if it is a Badaas model
// Returns error if the type is not a Badaas model
func (t Type) BadaasModelStruct() (*types.Struct, error) {
	structType, ok := t.Underlying().(*types.Struct)
	if !ok || !isBadaasModel(structType) {
		return nil, fmt.Errorf("type %s is not a Badaas Model", t.String())
	}

	return structType, nil
}

// Returns true if the type is a Badaas model
func isBadaasModel(structType *types.Struct) bool {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)

		if field.Embedded() && pie.Contains(badaasORMBaseModels, field.Type().String()) {
			return true
		}
	}

	return false
}

// Returns true is the type has a foreign key to the field's object
// (another field that references that object)
func (t Type) HasFK(field Field) (bool, error) {
	objectFields, err := getFields(t, "")
	if err != nil {
		return false, err
	}
	return pie.Any(objectFields, func(otherField Field) bool {
		return otherField.Name == field.getFKAttribute()
	}), nil
}

var scanMethod = regexp.MustCompile(`func \(\*.*\)\.Scan\([a-zA-Z0-9_-]* interface\{\}\) error$`)
var valueMethod = regexp.MustCompile(`func \(.*\)\.Value\(\) \(database/sql/driver\.Value\, error\)$`)

// Returns true if the type is a Gorm Custom type (https://gorm.io/docs/data_types.html)
func (t Type) IsGormCustomType() bool {
	typeNamed, isNamedType := t.Type.(*types.Named)
	if !isNamedType {
		return false
	}

	hasScanMethod := false
	hasValueMethod := false
	for i := 0; i < typeNamed.NumMethods(); i++ {
		methodSignature := typeNamed.Method(i).String()

		if !hasScanMethod && scanMethod.MatchString(methodSignature) {
			hasScanMethod = true
		} else if !hasValueMethod && valueMethod.MatchString(methodSignature) {
			hasValueMethod = true
		}
	}

	return hasScanMethod && hasValueMethod
}
