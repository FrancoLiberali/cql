package cmd

import (
	"errors"
	"go/types"

	"github.com/elliotchance/pie/v2"
)

// cql/model/models.go
var (
	modelIDs = []string{
		modelPath + "." + uIntID,
		modelPath + "." + uuid,
	}
	baseModelFields = []string{
		"ID", "CreatedAt", "UpdatedAt", "DeletedAt",
	}
)

type Field struct {
	Name         string
	NamePrefix   string
	Type         Type
	Embedded     bool
	Tags         GormTags
	ColumnPrefix string

	// GoAnonymous is true for Go-style embedded struct fields (the field
	// has no name, only a type — e.g. `model.UUIDModelWithTimestamps` in a
	// struct literal). Inner fields of such embeds are PROMOTED, so they
	// are accessible directly as dest.X. Distinguished from `Embedded`,
	// which is set for ANY embedding including named gorm-tagged fields
	// like `GormEmbedded ToBeGormEmbedded `gorm:"embedded"`` where the
	// inner field is accessed as dest.GormEmbedded.X.
	GoAnonymous bool

	// AccessPath is the chain of Go field names from the model root to
	// this leaf. Top-level fields have AccessPath == [Field.Name].
	// Fields under a Go-anonymous embed inherit the parent's path
	// unchanged (promotion). Fields under a named gorm-embedded struct
	// prepend the parent's Go name. The scanner generator uses this to
	// emit dest.A.B.C-style assignments — without it, gorm-embedded
	// columns silently land on the wrong field.
	AccessPath []string
}

func (field Field) CompleteName() string {
	return field.NamePrefix + field.Name
}

func (field Field) IsModelID() bool {
	return pie.Contains(modelIDs, field.TypeString())
}

func (field Field) IsUpdatable() bool {
	return !pie.Contains(baseModelFields, field.Name)
}

func (field Field) IsNullable() bool {
	return field.IsUpdatable() && (field.Type.IsSQLNullableType() || field.Type.WasPointer()) && !field.Tags.hasNotNull()
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
//
// TODO: this assumes the child's FK is named after the PARENT TYPE
// (<ParentType>ID), but gorm derives a belongsTo FK from the RELATION FIELD
// name (<field>ID). They only coincide when the child's belongsTo field is
// named exactly like the parent type. For role-named relations — e.g.
// `Author *User` (FK AuthorID, not UserID) — the has-many loader, collection
// and join FK references disagree with the conditions field the generator
// emits, producing code that doesn't compile (see the aligned-on-purpose
// hasmanyuint / hasmanywithpointers fixtures, and cql-gen/TODO.md). Fix:
// resolve the FK from the child's belongsTo field to this parent (its
// getFKAttribute), falling back to <ParentType>ID only when no such field
// exists.
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
func (field Field) ChangeType(newType types.Type, fromPointer bool) Field {
	return Field{
		Name: field.Name,
		Type: Type{Type: newType, wasPointer: fromPointer},
		Tags: field.Tags,
	}
}

// Get fields of a cql model
// Returns error is objectType is not a cql model
func getFields(objectType Type) ([]Field, error) {
	// The underlying type has to be a struct and a cql Model
	// (ignore const, var, func, etc.)
	structType, err := objectType.CQLModelStruct()
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
	for i := range numFields {
		fieldObject := structType.Field(i)
		gormTags := getGormTags(structType.Tag(i))
		fields = append(fields, Field{
			Name:        fieldObject.Name(),
			Type:        Type{Type: fieldObject.Type()},
			Embedded:    fieldObject.Embedded() || gormTags.hasEmbedded(),
			GoAnonymous: fieldObject.Embedded(),
			Tags:        gormTags,
		})
	}

	return fields, nil
}
