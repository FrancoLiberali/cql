package cmd

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"

	"github.com/FrancoLiberali/cql/cql-gen/cmd/log"
)

const (
	// cql/condition
	conditionPath           = cqlPath + "/condition"
	cqlCondition            = "Condition"
	cqlJoinCondition        = "JoinCondition"
	cqlNewJoinCondition     = "NewJoinCondition"
	cqlNewCollectionPreload = "NewCollectionPreloadCondition"
	cqlNewPreloadCondition  = "NewPreloadCondition"
	cqlField                = "Field"
	cqlUpdatableField       = "UpdatableField"
	cqlNullableField        = "NullableField"
	cqlBoolField            = "BoolField"
	cqlNullableBoolField    = "NullableBoolField"
	cqlStringField          = "StringField"
	cqlNullableStringField  = "NullableStringField"
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
		// adapt param to slice and
		// generate code for the type of the elements of the slice
		condition.param.ToSlice()
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
		condition.generate(
			objectType,
			field,
		)
	case *types.Named:
		// slice of named types (user defined types)
		_, err := field.Type.CQLModelStruct()
		if err == nil {
			// field is a CQL Model
			condition.generateCollectionPreload(objectType, field)
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
	fieldValues := jen.Dict{
		jen.Id("Name"): jen.Lit(field.Name),
	}

	columnName := field.getColumnName()

	if columnName != "" {
		fieldValues[jen.Id("Column")] = jen.Lit(columnName)
	}

	columnPrefix := field.ColumnPrefix
	if columnPrefix != "" {
		fieldValues[jen.Id("ColumnPrefix")] = jen.Lit(columnPrefix)
	}

	objectTypeQual := jen.Qual(
		getRelativePackagePath(condition.destPkg, objectType),
		objectType.Name(),
	)

	fieldQual := jen.Qual(
		conditionPath, cqlField,
	).Types(
		objectTypeQual,
		condition.param.GenericType(),
	)

	if field.IsUpdatable() {
		fieldValues = jen.Dict{
			jen.Id(cqlField): fieldQual.Clone().Values(fieldValues),
		}
		fieldQual = jen.Qual(
			conditionPath, cqlUpdatableField,
		).Types(
			objectTypeQual,
			condition.param.GenericType(),
		)

		if field.IsNullable() {
			fieldValues = jen.Dict{
				jen.Id(cqlUpdatableField): fieldQual.Clone().Values(fieldValues),
			}
			fieldQual = jen.Qual(
				conditionPath, cqlNullableField,
			).Types(
				objectTypeQual,
				condition.param.GenericType(),
			)
		}
	}

	if condition.param.isString {
		fieldQual, fieldValues = condition.transformIntoSpecificField(field, objectTypeQual, fieldQual, fieldValues, cqlNullableStringField, cqlStringField)
	} else if condition.param.isBool {
		fieldQual, fieldValues = condition.transformIntoSpecificField(field, objectTypeQual, fieldQual, fieldValues, cqlNullableBoolField, cqlBoolField)
	}

	condition.FieldType = jen.Id(condition.FieldName).Add(
		fieldQual,
	)

	condition.FieldDefinition = fieldQual.Clone().Values(fieldValues)
}

// Transforms the fieldQual and the fieldValues into a specific type (string, bool) depending if the type is nullable or not
func (condition *Condition) transformIntoSpecificField(
	field Field, objectTypeQual *jen.Statement,
	fieldQual *jen.Statement, fieldValues jen.Dict,
	nullableType string, notNullableType string,
) (*jen.Statement, jen.Dict) {
	if field.IsNullable() {
		fieldValues = jen.Dict{
			jen.Id(cqlNullableField): fieldQual.Clone().Values(fieldValues),
		}
		fieldQual = jen.Qual(
			conditionPath, nullableType,
		).Types(objectTypeQual)
	} else {
		fieldValues = jen.Dict{
			jen.Id(cqlUpdatableField): fieldQual.Clone().Values(fieldValues),
		}
		fieldQual = jen.Qual(
			conditionPath, notNullableType,
		).Types(objectTypeQual)
	}

	return fieldQual, fieldValues
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

func (condition *Condition) generateCollectionPreload(objectType Type, field Field) {
	t1 := jen.Qual(
		getRelativePackagePath(condition.destPkg, objectType),
		objectType.Name(),
	)

	t2 := jen.Qual(
		getRelativePackagePath(condition.destPkg, field.Type),
		field.TypeName(),
	)

	ormT1Condition := jen.Qual(
		conditionPath, cqlCondition,
	).Types(t1)
	ormT2IJoinCondition := jen.Qual(
		conditionPath, cqlJoinCondition,
	).Types(t2)
	ormNewCollectionPreload := jen.Qual(
		conditionPath, cqlNewCollectionPreload,
	).Types(
		t1, t2,
	)

	condition.PreloadRelationMethod = createMethod(condition.modelType, "Preload"+field.Name).Params(
		jen.Id("nestedPreloads").Op("...").Add(ormT2IJoinCondition),
	).Add(
		ormT1Condition,
	).Block(
		jen.Return(
			ormNewCollectionPreload.Call(
				jen.Lit(field.Name),
				jen.Id("nestedPreloads"),
			),
		),
	)
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
