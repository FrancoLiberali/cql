package conditions

import (
	"errors"
	"go/types"

	"github.com/elliotchance/pie/v2"
)

// badaas/orm/baseModels.go
var modelIDs = []string{
	badaasORMPath + "." + uIntID,
	badaasORMPath + "." + uuid,
}

type Field struct {
	Name         string
	NamePrefix   string
	Type         Type
	Embedded     bool
	Tags         GormTags
	ColumnPrefix string
}

func (field Field) IsModelID() bool {
	return pie.Contains(modelIDs, field.TypeString())
}

// Get the name of the column where the data for a field will be saved
func (field Field) getColumnName() string {
	columnTag, isPresent := field.Tags[columnTagName]
	if isPresent {
		// field has a gorm column tag, so the name of the column will be that tag
		return columnTag
	}

	return ""
}

// Get name of the attribute of the object that is a foreign key to the field's object
func (field Field) getFKAttribute() string {
	foreignKeyTag, isPresent := field.Tags[foreignKeyTagName]
	if isPresent {
		// field has a foreign key tag, so the name will be that tag
		return foreignKeyTag
	}

	// gorm default
	return field.Name + "ID"
}

// Get name of the attribute of the field's object that is references by the foreign key
func (field Field) getFKReferencesAttribute() string {
	referencesTag, isPresent := field.Tags[referencesTagName]
	if isPresent {
		// field has a references tag, so the name will be that tag
		return referencesTag
	}

	// gorm default
	return "ID"
}

// Get name of the attribute of field's object that is a foreign key to the object
func (field Field) getRelatedTypeFKAttribute(structName string) string {
	foreignKeyTag, isPresent := field.Tags[foreignKeyTagName]
	if isPresent {
		// field has a foreign key tag, so the name will that tag
		return foreignKeyTag
	}

	// gorm default
	return structName + "ID"
}

func (field Field) GetType() types.Type {
	return field.Type.Type
}

// Get field's type full string (pkg + name)
func (field Field) TypeString() string {
	return field.Type.String()
}

// Get field's type name
func (field Field) TypeName() string {
	return field.Type.Name()
}

// Create a new field with the same name and tags but a different type
func (field Field) ChangeType(newType types.Type) Field {
	return Field{
		Name: field.Name,
		Type: Type{newType},
		Tags: field.Tags,
	}
}

// Get fields of a Badaas model
// Returns error is objectType is not a Badaas model
func getFields(objectType Type) ([]Field, error) {
	// The underlying type has to be a struct and a Badaas Model
	// (ignore const, var, func, etc.)
	structType, err := objectType.BadaasModelStruct()
	if err != nil {
		return nil, err
	}

	return getStructFields(structType)
}

// Get fields of a struct
// Returns errors if the struct has not fields
func getStructFields(structType *types.Struct) ([]Field, error) {
	numFields := structType.NumFields()
	if numFields == 0 {
		return nil, errors.New("struct has 0 fields")
	}

	fields := []Field{}

	// Iterate over struct fields
	for i := 0; i < numFields; i++ {
		fieldObject := structType.Field(i)
		gormTags := getGormTags(structType.Tag(i))
		fields = append(fields, Field{
			Name:     fieldObject.Name(),
			Type:     Type{fieldObject.Type()},
			Embedded: fieldObject.Embedded() || gormTags.hasEmbedded(),
			Tags:     gormTags,
		})
	}

	return fields, nil
}
