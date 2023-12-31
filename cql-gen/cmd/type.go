package cmd

import (
	"errors"
	"fmt"
	"go/types"
	"regexp"
	"strings"

	"github.com/elliotchance/pie/v2"

	"github.com/FrancoLiberali/cql/cql-gen/cmd/utils"
)

var (
	// cql/model/models.go
	cqlBaseModels = []string{
		modelPath + "." + uuidModel,
		modelPath + "." + uIntModel,
	}

	// database/sql
	nullString       = "database/sql.NullString"
	nullInt64        = "database/sql.NullInt64"
	nullInt32        = "database/sql.NullInt32"
	nullInt16        = "database/sql.NullInt16"
	nullFloat64      = "database/sql.NullFloat64"
	nullByte         = "database/sql.NullByte"
	nullBool         = "database/sql.NullBool"
	nullTime         = "database/sql.NullTime"
	deletedAt        = "gorm.io/gorm.DeletedAt"
	sqlNullableTypes = []string{
		nullString, nullInt64, nullInt32, nullInt16, nullFloat64,
		nullByte, nullBool, nullTime, deletedAt,
	}
)

var ErrFkNotInTypeFields = errors.New("fk not in type's fields")

type Type struct {
	types.Type

	wasPointer bool
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

// Get the struct under type if it is a cql model
// Returns error if the type is not a cql model
func (t Type) CQLModelStruct() (*types.Struct, error) {
	structType, ok := t.Underlying().(*types.Struct)
	if !ok || !isCQLModel(structType) {
		return nil, fmt.Errorf("type %s is not a cql Model", t.String())
	}

	return structType, nil
}

// Returns true if the type is a cql model
func isCQLModel(structType *types.Struct) bool {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)

		if field.Embedded() && isBaseModel(field.Type().String()) {
			return true
		}
	}

	return false
}

func isBaseModel(fieldName string) bool {
	return pie.Contains(cqlBaseModels, fieldName)
}

// Returns the fk field of the type to the "field"'s object
// (another field that references that object)
func (t Type) GetFK(field Field) (*Field, error) {
	objectFields, err := getFields(t)
	if err != nil {
		return nil, err
	}

	fk := utils.FindFirst(objectFields, func(otherField Field) bool {
		return strings.EqualFold(otherField.Name, field.getFKAttribute())
	})

	if fk == nil {
		return nil, ErrFkNotInTypeFields
	}

	return fk, nil
}

var (
	scanMethod  = regexp.MustCompile(`func \(\*.*\)\.Scan\([a-zA-Z0-9_-]* (interface\{\}|any)\) error$`)
	valueMethod = regexp.MustCompile(`func \(.*\)\.Value\(\) \(database/sql/driver\.Value\, error\)$`)
)

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

// Returns true if the type is a sql nullable type (sql.NullBool, sql.NullInt, etc.)
func (t Type) IsSQLNullableType() bool {
	return pie.Contains(sqlNullableTypes, t.String())
}

// Returns true if the type used to be a pointer (pointers are transformed into the pointed type)
func (t Type) WasPointer() bool {
	return t.wasPointer
}
