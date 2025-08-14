package conditions

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"

	"github.com/ditrit/badaas-cli/cmd/log"
)

const (
	// badaas/orm/cql
	cqlPath                       = badaasORMPath + "/cql"
	badaasORMCondition            = "Condition"
	badaasORMJoinCondition        = "JoinCondition"
	badaasORMNewJoinCondition     = "NewJoinCondition"
	badaasORMNewCollectionPreload = "NewCollectionPreloadCondition"
	badaasORMNewPreloadCondition  = "NewPreloadCondition"
	badaasORMField                = "Field"
	badaasORMBoolField            = "BoolField"
	badaasORMStringField          = "StringField"
	// badaas/orm/model
	modelPath = badaasORMPath + "/model"
	uIntID    = "UIntID"
	uuid      = "UUID"
	uuidModel = "UUIDModel"
	uIntModel = "UIntModel"
)

const preloadMethod = "Preload"

type Condition struct {
	FieldName             string
	FieldType             *jen.Statement
	FieldDefinition       *jen.Statement
	ConditionMethod       *jen.Statement
	PreloadRelationName   string
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
			field.ChangeType(fieldType.Elem()),
		)
	case *types.Slice:
		// the field is a slice
		// adapt param to slice and
		// generate code for the type of the elements of the slice
		condition.param.ToSlice()
		condition.generateForSlice(
			objectType,
			field.ChangeType(fieldType.Elem()),
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
		_, err := field.Type.BadaasModelStruct()
		if err == nil {
			// field is a Badaas Model
			condition.generateCollectionPreload(objectType, field)
		}
	case *types.Pointer:
		// slice of pointers, generate code for a slice of the pointed type
		condition.generateForSlice(
			objectType,
			field.ChangeType(elemType.Elem()),
		)
	default:
		log.Logger.Debugf("struct field list elem type not handled: %T", elemType)
	}
}

// Generate condition between object and field when the field is a named type (user defined struct)
func (condition *Condition) generateForNamedType(objectType Type, field Field) {
	_, err := field.Type.BadaasModelStruct()

	switch {
	case err == nil:
		// field is a badaas model
		condition.generateForBadaasModel(objectType, field)
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
		// or a badaas-orm id (uuid or uintid)
		condition.param.ToCustomType(condition.destPkg, field.Type)
		condition.createField(
			objectType,
			field,
		)
	default:
		log.Logger.Debugf("struct field type not handled: %s", field.TypeString())
	}
}

// Generate condition between object and field when the field is a Badaas Model
func (condition *Condition) generateForBadaasModel(objectType Type, field Field) {
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
		cqlPath, badaasORMField,
	).Types(
		objectTypeQual,
		condition.param.GenericType(),
	)
	if condition.param.isString {
		fieldValues = jen.Dict{
			jen.Id("Field"): fieldQual.Clone().Values(fieldValues),
		}
		fieldQual = jen.Qual(
			cqlPath, badaasORMStringField,
		).Types(objectTypeQual)
	} else if condition.param.isBool {
		fieldValues = jen.Dict{
			jen.Id("Field"): fieldQual.Clone().Values(fieldValues),
		}
		fieldQual = jen.Qual(
			cqlPath, badaasORMBoolField,
		).Types(objectTypeQual)
	}

	condition.FieldType = jen.Id(condition.FieldName).Add(
		fieldQual,
	)

	condition.FieldDefinition = fieldQual.Clone().Values(fieldValues)
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
		cqlPath, badaasORMJoinCondition,
	).Types(t1)
	ormT2Condition := jen.Qual(
		cqlPath, badaasORMCondition,
	).Types(t2)
	ormJoinCondition := jen.Qual(
		cqlPath, badaasORMNewJoinCondition,
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
			),
		),
	)

	// preload for the relation
	condition.setPreloadRelationName(field)

	condition.PreloadRelationMethod = createMethod(condition.modelType, condition.PreloadRelationName).Params().Add(
		ormT1IJoinCondition,
	).Block(
		jen.Return(jen.Id(condition.modelType).Dot(conditionName).Call(
			jen.Id(field.TypeName()).Dot(preloadMethod).Call(),
		)),
	)
}

func (condition *Condition) setPreloadRelationName(field Field) {
	condition.PreloadRelationName = "Preload" + field.Name
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
		cqlPath, badaasORMCondition,
	).Types(t1)
	ormT2IJoinCondition := jen.Qual(
		cqlPath, badaasORMJoinCondition,
	).Types(t2)
	ormNewCollectionPreload := jen.Qual(
		cqlPath, badaasORMNewCollectionPreload,
	).Types(
		t1, t2,
	)

	condition.setPreloadRelationName(field)

	condition.PreloadRelationMethod = createMethod(condition.modelType, condition.PreloadRelationName).Params(
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
