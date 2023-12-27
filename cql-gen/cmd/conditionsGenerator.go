package cmd

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"

	"github.com/FrancoLiberali/cql/cql-gen/cmd/log"
)

type ConditionsGenerator struct {
	object     types.Object
	objectType Type
}

func NewConditionsGenerator(object types.Object) *ConditionsGenerator {
	return &ConditionsGenerator{
		object:     object,
		objectType: Type{Type: object.Type()},
	}
}

// Add conditions for an object in the file
func (cg ConditionsGenerator) Into(file *File) error {
	fields, err := getFields(cg.objectType)
	if err != nil {
		return err
	}

	log.Logger.Infof("Generating conditions for type %q in %s", cg.object.Name(), file.name)

	// Add one condition for each field of the object
	conditions := cg.ForEachField(file, fields)

	objectName := cg.object.Name()
	objectQual := jen.Qual(
		getRelativePackagePath(file.destPkg, cg.objectType),
		cg.objectType.Name(),
	)

	fieldIdentifiers := []jen.Code{}

	conditionsModelType := getConditionsModelType(objectName)
	conditionsModelAttributesDef := []jen.Code{}
	conditionsModelAttributesIns := jen.Dict{}

	for _, condition := range conditions {
		if condition.ConditionMethod != nil {
			file.Add(condition.ConditionMethod)
		}

		// add all field names to the list of fields of the preload condition
		if condition.FieldDefinition != nil {
			conditionsModelAttributesDef = append(conditionsModelAttributesDef, condition.FieldType)
			conditionsModelAttributesIns[jen.Id(condition.FieldName)] = condition.FieldDefinition
			fieldIdentifiers = append(
				fieldIdentifiers,
				jen.Id(conditionsModelType).Dot(condition.FieldName),
			)
		}

		// add the preload to the list of all possible preloads
		if condition.PreloadRelationMethod != nil {
			file.Add(condition.PreloadRelationMethod)
		}
	}

	addConditionsModelDefinition(file, conditionsModelType, conditionsModelAttributesDef)
	addConditionsModelInstantiation(file, objectName, conditionsModelType, conditionsModelAttributesIns)
	addPreloadMethod(file, objectName, objectQual, conditionsModelType, fieldIdentifiers)

	return nil
}

func addPreloadMethod(file *File, objectName string, objectQual *jen.Statement, conditionsModelType string, fieldIdentifiers []jen.Code) {
	file.Add(
		jen.Comment(fmt.Sprintf("Preload allows preloading the %s when doing a query", objectName)),
		createMethod(conditionsModelType, preloadMethod).Params().Add(
			jen.Qual(
				conditionPath, cqlCondition,
			).Types(
				objectQual,
			),
		).Block(
			jen.Return(
				jen.Qual(
					conditionPath, cqlNewPreloadCondition,
				).Types(
					objectQual,
				).Call(fieldIdentifiers...),
			),
		),
	)
}

func addConditionsModelInstantiation(file *File, objectName, conditionsModelType string, conditionsModelAttributes jen.Dict) {
	file.Add(
		jen.Var().Id(
			objectName,
		).Op("=").Add(
			jen.Id(conditionsModelType).Values(
				conditionsModelAttributes,
			),
		),
	)
}

func addConditionsModelDefinition(file *File, conditionsModelType string, conditionsModelAttributes []jen.Code) {
	file.Add(
		jen.Type().Id(
			conditionsModelType,
		).Struct(
			conditionsModelAttributes...,
		),
	)
}

func getConditionsModelType(objectName string) string {
	return strcase.ToCamel(objectName) + "Conditions"
}

// Generate the conditions for each of the object's fields
func (cg ConditionsGenerator) ForEachField(file *File, fields []Field) []Condition {
	conditions := []Condition{}

	for _, field := range fields {
		log.Logger.Debugf("Generating condition for field %q", field.Name)

		if field.Embedded {
			conditions = append(
				conditions,
				generateForEmbeddedField[Condition](
					file,
					field,
					cg,
				)...,
			)
		} else {
			conditions = append(conditions, *NewCondition(
				file.destPkg, cg.objectType, field,
			))
		}
	}

	return conditions
}
