package cmd

import (
	"errors"
	"go/types"

	"github.com/elliotchance/pie/v2"
)

// Generate conditions for a embedded field using the "generator"
// it will generate a condition for each of the field of the embedded field's type
func generateForEmbeddedField[T any](file *File, field Field, generator CodeGenerator[T]) []T {
	embeddedStructType, ok := field.Type.Underlying().(*types.Struct)
	if !ok {
		panic(errors.New("unreachable! embedded objects are always structs"))
	}

	fields, err := getStructFields(embeddedStructType)
	if err != nil {
		// embedded field's type has not fields
		return []T{}
	}

	if !isBaseModel(field.TypeString()) {
		fields = pie.Map(fields, func(embeddedField Field) Field {
			embeddedField.ColumnPrefix = field.Tags.getEmbeddedPrefix()
			embeddedField.NamePrefix = field.Name

			return embeddedField
		})
	}

	return generator.ForEachField(file, fields)
}
