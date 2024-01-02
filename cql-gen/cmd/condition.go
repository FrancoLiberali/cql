package cmd

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"

	"github.com/FrancoLiberali/cql/cql-gen/cmd/log"
)

const (
	// cql/condition
	conditionPath             = cqlPath + "/condition"
	cqlCondition              = "Condition"
	cqlJoinCondition          = "JoinCondition"
	cqlNewJoinCondition       = "NewJoinCondition"
	cqlNewPreloadCondition    = "NewPreloadCondition"
	cqlField                  = "Field"
	cqlNewField               = "NewField"
	cqlUpdatableField         = "UpdatableField"
	cqlNewUpdatableField      = "NewUpdatableField"
	cqlNullableField          = "NullableField"
	cqlNewNullableField       = "NewNullableField"
	cqlBoolField              = "BoolField"
	cqlNewBoolField           = "NewBoolField"
	cqlNullableBoolField      = "NullableBoolField"
	cqlNewNullableBoolField   = "NewNullableBoolField"
	cqlStringField            = "StringField"
	cqlNewStringField         = "NewStringField"
	cqlNullableStringField    = "NullableStringField"
	cqlNewNullableStringField = "NewNullableStringField"
	cqlCollection             = "Collection"
	cqlNewCollection          = "NewCollection"
	// cql/model
	modelPath = cqlPath + "/model"
	uIntID    = "UIntID"
	uuid      = "UUID"
	uuidModel = "UUIDModel"
	uIntModel = "UIntModel"
)

const preloadMethod = "preload"

type Condition struct {
	FieldName             string
	FieldType             *jen.Statement
	FieldDefinition       *jen.Statement
	FieldIsCollection     bool
	ConditionMethod       *jen.Statement
	PreloadRelationMethod *jen.Statement
	param                 *JenParam
	destPkg               string
	modelType             string
}

func NewCondition(destPkg string, objectType Type, field Field) *Condition {
	condition := &Condition{
		FieldName: field.CompleteName(),
		param:     NewJenParam(),
		destPkg:   destPkg,
		modelType: getConditionsModelType(objectType.Name()),
	}
	condition.generate(objectType, field)

	return condition
}

// Generate the condition between the object and the field
func (condition *Condition) generate(objectType Type, field Field) {
	switch fieldType := field.GetType().(type) {
	case *types.Basic:
		// the field is a basic type (string, int, etc)
		// adapt param to that type and generate a WhereCondition
		condition.param.ToBasicKind(fieldType)
		condition.createField(
			objectType,
			field,
		)
	case *types.Named:
		// the field is a named type (user defined structs)
		condition.generateForNamedType(
			objectType,
			field,
		)
	case *types.Pointer:
		// the field is a pointer
		condition.generate(
			objectType,
			field.ChangeType(fieldType.Elem(), true),
		)
	case *types.Slice:
		// the field is a slice
		// generate code for the type of the elements of the slice
		condition.generateForSlice(
			objectType,
			field.ChangeType(fieldType.Elem(), false),
		)
	default:
		log.Logger.Debugf("struct field type not handled: %T", fieldType)
	}
}

// Generate condition between the object and the field when the field is a slice
func (condition *Condition) generateForSlice(objectType Type, field Field) {
	switch elemType := field.GetType().(type) {
	case *types.Basic:
		// slice of basic types ([]string, []int, etc.)
		// the only one supported directly by gorm is []byte
		// but the others can be used after configuration in some dbs
		condition.param.ToSlice()
		condition.generate(
			objectType,
			field,
		)
	case *types.Named:
		// slice of named types (user defined types)
		_, err := field.Type.CQLModelStruct()
		if err == nil {
			// field is a CQL Model
			condition.createCollection(objectType, field)
		}
	case *types.Pointer:
		// slice of pointers, generate code for a slice of the pointed type
		condition.generateForSlice(
			objectType,
			field.ChangeType(elemType.Elem(), false),
		)
	default:
		log.Logger.Debugf("struct field list elem type not handled: %T", elemType)
	}
}

// Generate condition between object and field when the field is a named type (user defined struct)
func (condition *Condition) generateForNamedType(objectType Type, field Field) {
	_, err := field.Type.CQLModelStruct()

	switch {
	case err == nil:
		// field is a cql model
		condition.generateForCQLModel(objectType, field)
	case field.Type.IsSQLNullableType():
		// field is a sql nullable type (sql.NullBool, sql.NullInt, etc.)
		condition.param.SQLToBasicType(field.Type)
		condition.createField(
			objectType,
			field,
		)
	case field.Type.IsGormCustomType() || field.TypeString() == "time.Time" || field.IsModelID():
		// field is a Gorm Custom type (implements Scanner and Valuer interfaces)
		// or a named type supported by gorm (time.Time)
		// or a cql id (uuid or uintid)
		condition.param.ToCustomType(condition.destPkg, field.Type)
		condition.createField(
			objectType,
			field,
		)
	default:
		log.Logger.Debugf("struct field type not handled: %s", field.TypeString())
	}
}

// Generate condition between object and field when the field is a CQL Model
func (condition *Condition) generateForCQLModel(objectType Type, field Field) {
	_, err := objectType.GetFK(field)
	if err == nil {
		// has the fk -> belongsTo relation
		condition.generateJoinWithFK(
			objectType,
			field,
		)
	} else {
		// has not the fk -> hasOne relation
		condition.generateJoinWithoutFK(
			objectType,
			field,
		)
	}
}

func createMethod(typeName, methodName string) *jen.Statement {
	return jen.Func().Params(
		jen.Id(typeName).Id(typeName),
	).Id(methodName)
}

// create a variable containing the definition of the field identifier
// to use it in the where condition and in the preload condition
func (condition *Condition) createField(objectType Type, field Field) {
	fieldName := jen.Lit(field.Name)
	fieldColumn := jen.Lit("")
	fieldColumnPrefix := jen.Lit("")

	columnName := field.getColumnName()

	if columnName != "" {
		fieldColumn = jen.Lit(columnName)
	}

	columnPrefix := field.ColumnPrefix
	if columnPrefix != "" {
		fieldColumnPrefix = jen.Lit(columnPrefix)
	}

	objectTypeQual := jen.Qual(
		getRelativePackagePath(condition.destPkg, objectType),
		objectType.Name(),
	)

	var fieldQual *jen.Statement
	var newFieldQual *jen.Statement

	if condition.param.isString {
		fieldQual, newFieldQual = condition.specificField(field, objectTypeQual, cqlNullableStringField, cqlNewNullableStringField, cqlStringField, cqlNewStringField)
	} else if condition.param.isBool {
		fieldQual, newFieldQual = condition.specificField(field, objectTypeQual, cqlNullableBoolField, cqlNewNullableBoolField, cqlBoolField, cqlNewBoolField)
	} else {
		if field.IsNullable() {
			fieldQual = jen.Qual(conditionPath, cqlNullableField)
			newFieldQual = jen.Qual(conditionPath, cqlNewNullableField)
		} else if field.IsUpdatable() {
			fieldQual = jen.Qual(conditionPath, cqlUpdatableField)
			newFieldQual = jen.Qual(conditionPath, cqlNewUpdatableField)
		} else {
			fieldQual = jen.Qual(conditionPath, cqlField)
			newFieldQual = jen.Qual(conditionPath, cqlNewField)
		}

		fieldQual = fieldQual.Types(
			objectTypeQual,
			condition.param.GenericType(),
		)
		newFieldQual = newFieldQual.Types(
			objectTypeQual,
			condition.param.GenericType(),
		)
	}

	condition.FieldType = jen.Id(condition.FieldName).Add(
		fieldQual,
	)

	condition.FieldDefinition = newFieldQual.Call(fieldName, fieldColumn, fieldColumnPrefix)
}

func (condition *Condition) specificField(
	field Field, objectTypeQual *jen.Statement,
	nullableType, newNullableType string,
	notNullableType, newNotNullableType string,
) (*jen.Statement, *jen.Statement) {
	var fieldQual *jen.Statement
	var newFieldQual *jen.Statement

	if field.IsNullable() {
		fieldQual = jen.Qual(conditionPath, nullableType)
		newFieldQual = jen.Qual(conditionPath, newNullableType)
	} else {
		fieldQual = jen.Qual(conditionPath, notNullableType)
		newFieldQual = jen.Qual(conditionPath, newNotNullableType)
	}

	fieldQual = fieldQual.Types(objectTypeQual)
	newFieldQual = newFieldQual.Types(objectTypeQual)

	return fieldQual, newFieldQual
}

// Generate a JoinCondition between the object and field's object
// when object has a foreign key to the field's object
func (condition *Condition) generateJoinWithFK(objectType Type, field Field) {
	condition.generateJoin(
		objectType,
		field,
		field.getFKAttribute(),
		field.getFKReferencesAttribute(),
	)
}

// Generate a JoinCondition between the object and field's object
// when object has not a foreign key to the field's object
// (so the field's object has it)
func (condition *Condition) generateJoinWithoutFK(objectType Type, field Field) {
	condition.generateJoin(
		objectType,
		field,
		field.getFKReferencesAttribute(),
		field.getRelatedTypeFKAttribute(objectType.Name()),
	)
}

// Generate a JoinCondition
func (condition *Condition) generateJoin(objectType Type, field Field, t1Field, t2Field string) {
	t1 := jen.Qual(
		getRelativePackagePath(condition.destPkg, objectType),
		objectType.Name(),
	)

	t2 := jen.Qual(
		getRelativePackagePath(condition.destPkg, field.Type),
		field.TypeName(),
	)

	conditionName := getConditionName(field)
	log.Logger.Debugf("Generated %q", conditionName)

	ormT1IJoinCondition := jen.Qual(
		conditionPath, cqlJoinCondition,
	).Types(t1)
	ormT2Condition := jen.Qual(
		conditionPath, cqlCondition,
	).Types(t2)
	ormJoinCondition := jen.Qual(
		conditionPath, cqlNewJoinCondition,
	).Types(
		t1, t2,
	)

	condition.ConditionMethod = createMethod(condition.modelType, conditionName).Params(
		jen.Id("conditions").Op("...").Add(ormT2Condition),
	).Add(
		ormT1IJoinCondition,
	).Block(
		jen.Return(
			ormJoinCondition.Call(
				jen.Id("conditions"),
				jen.Lit(field.Name),
				jen.Lit(t1Field),
				jen.Id(condition.modelType).Dot(preloadMethod).Call(),
				jen.Lit(t2Field),
				jen.Id(field.Type.Name()).Dot(preloadMethod).Call(),
			),
		),
	)
}

func (condition *Condition) createCollection(objectType Type, field Field) {
	t1 := jen.Qual(
		getRelativePackagePath(condition.destPkg, objectType),
		objectType.Name(),
	)

	t2 := jen.Qual(
		getRelativePackagePath(condition.destPkg, field.Type),
		field.TypeName(),
	)

	condition.FieldType = jen.Id(condition.FieldName).Add(
		jen.Qual(
			conditionPath, cqlCollection,
		).Types(
			t1,
			t2,
		),
	)

	condition.FieldDefinition = jen.Qual(
		conditionPath, cqlNewCollection,
	).Types(
		t1,
		t2,
	).Call(
		jen.Lit(field.Name),
		jen.Lit(field.getFKReferencesAttribute()),
		jen.Lit(field.getRelatedTypeFKAttribute(objectType.Name())),
	)

	condition.FieldIsCollection = true
}

// Generate condition names
func getConditionName(field Field) string {
	return strcase.ToPascal(field.NamePrefix) + strcase.ToPascal(field.Name)
}

// Avoid importing the same package as the destination one
func getRelativePackagePath(destPkg string, typeV Type) string {
	srcPkg := typeV.Pkg()
	if srcPkg == nil || srcPkg.Name() == destPkg {
		return ""
	}

	return srcPkg.Path()
}
