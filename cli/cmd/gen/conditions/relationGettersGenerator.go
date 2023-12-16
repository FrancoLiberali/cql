package conditions

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"

	"github.com/ditrit/badaas-cli/cmd/log"
)

const (
	badaasORMVerifyStructLoaded        = "VerifyStructLoaded"
	badaasORMVerifyPointerLoaded       = "VerifyPointerLoaded"
	badaasORMVerifyPointerWithIDLoaded = "VerifyPointerWithIDLoaded"
	badaasORMVerifyCollectionLoaded    = "VerifyCollectionLoaded"
)

type RelationGettersGenerator struct {
	object     types.Object
	objectType Type
}

func NewRelationGettersGenerator(object types.Object) *RelationGettersGenerator {
	return &RelationGettersGenerator{
		object:     object,
		objectType: Type{object.Type()},
	}
}

// Add conditions for an object in the file
func (generator RelationGettersGenerator) Into(file *File) error {
	fields, err := getFields(generator.objectType)
	if err != nil {
		return err
	}

	log.Logger.Infof("Generating relation getters for type %q in %s", generator.object.Name(), file.name)

	file.Add(generator.ForEachField(file, fields)...)

	return nil
}

func (generator RelationGettersGenerator) ForEachField(file *File, fields []Field) []jen.Code {
	relationGetters := []jen.Code{}

	for _, field := range fields {
		if field.Embedded {
			relationGetters = append(
				relationGetters,
				generateForEmbeddedField[jen.Code](
					file,
					field,
					generator,
				)...,
			)
		} else {
			getterForField := generator.generateForField(field)
			if getterForField != nil {
				relationGetters = append(relationGetters, getterForField)
			}
		}
	}

	return relationGetters
}

func (generator RelationGettersGenerator) generateForField(field Field) jen.Code {
	switch fieldType := field.GetType().(type) {
	case *types.Named:
		// the field is a named type (user defined structs)
		_, err := field.Type.BadaasModelStruct()
		if err == nil {
			log.Logger.Debugf("Generating relation getter for type %q and field %s", generator.object.Name(), field.Name)
			// field is a badaas Model
			return generator.verifyStruct(field)
		}
	case *types.Pointer:
		// the field is a pointer
		return generator.generateForPointer(field.ChangeType(fieldType.Elem()))
	default:
		log.Logger.Debugf("struct field type not handled: %T", fieldType)
	}

	return nil
}

func (generator RelationGettersGenerator) generateForPointer(field Field) jen.Code {
	switch fieldType := field.GetType().(type) {
	case *types.Named:
		_, err := field.Type.BadaasModelStruct()
		if err == nil {
			// field is a pointer to Badaas Model
			fk, err := generator.objectType.GetFK(field)
			if err != nil {
				log.Logger.Debugf("unhandled: field is a pointer and object not has the fk: %s", field.Type)
				return nil
			}

			log.Logger.Debugf("Generating relation getter for type %q and field %s", generator.object.Name(), field.Name)

			switch fk.GetType().(type) {
			case *types.Named:
				if fk.IsModelID() {
					return generator.verifyPointerWithID(field)
				}
			case *types.Pointer:
				// the fk is a pointer
				return generator.verifyPointer(field)
			}
		}
	case *types.Slice:
		return generator.generateForSlicePointer(
			field.ChangeType(fieldType.Elem()),
			nil,
		)
	}

	return nil
}

func (generator RelationGettersGenerator) generateForSlicePointer(field Field, fieldTypePrefix *jen.Statement) jen.Code {
	switch fieldType := field.GetType().(type) {
	case *types.Named:
		_, err := field.Type.BadaasModelStruct()
		if err == nil {
			// field is a pointer to a slice of badaas Model
			return generator.verifyCollection(field, fieldTypePrefix)
		}
	case *types.Pointer:
		return generator.generateForSlicePointer(
			field.ChangeType(fieldType.Elem()),
			jen.Op("*"),
		)
	}

	return nil
}

func getGetterName(field Field) string {
	return "Get" + strcase.ToPascal(field.Name)
}

func (generator RelationGettersGenerator) verifyStruct(field Field) *jen.Statement {
	return generator.verifyCommon(
		field,
		badaasORMVerifyStructLoaded,
		jen.Op("*"),
		nil,
		jen.Op("&").Id("m").Op(".").Id(field.Name),
	)
}

func (generator RelationGettersGenerator) verifyPointer(field Field) *jen.Statement {
	return generator.verifyPointerCommon(field, badaasORMVerifyPointerLoaded)
}

func (generator RelationGettersGenerator) verifyPointerWithID(field Field) *jen.Statement {
	return generator.verifyPointerCommon(field, badaasORMVerifyPointerWithIDLoaded)
}

func (generator RelationGettersGenerator) verifyCollection(field Field, fieldTypePrefix *jen.Statement) jen.Code {
	return generator.verifyCommon(
		field,
		badaasORMVerifyCollectionLoaded,
		jen.Index(),
		fieldTypePrefix,
		jen.Id("m").Op(".").Id(field.Name),
	)
}

func (generator RelationGettersGenerator) verifyPointerCommon(field Field, verifyFunc string) *jen.Statement {
	return generator.verifyCommon(
		field,
		verifyFunc,
		jen.Op("*"),
		nil,
		jen.Id("m").Op(".").Id(field.Name+"ID"),
		jen.Id("m").Op(".").Id(field.Name),
	)
}

func (generator RelationGettersGenerator) verifyCommon(
	field Field,
	verifyFunc string,
	returnType *jen.Statement,
	fieldTypePrefix *jen.Statement,
	callParams ...jen.Code,
) *jen.Statement {
	fieldType := jen.Qual(
		getRelativePackagePath(
			generator.object.Pkg().Name(),
			field.Type,
		),
		field.TypeName(),
	)

	if fieldTypePrefix != nil {
		fieldType = fieldTypePrefix.Add(fieldType)
	}

	return jen.Func().Parens(
		jen.Id("m").Id(generator.object.Name()),
	).Id(getGetterName(field)).Params().Add(
		jen.Parens(
			jen.List(
				returnType.Add(fieldType),
				jen.Id("error"),
			),
		),
	).Block(
		jen.Return(
			jen.Qual(
				badaasORMPath,
				verifyFunc,
			).Types(
				fieldType,
			).Call(
				callParams...,
			),
		),
	)
}
