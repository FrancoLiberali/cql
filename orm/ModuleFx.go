package orm

import (
	"fmt"
	"log"
	"reflect"

	"go.uber.org/fx"
)

type GetModelsResult struct {
	fx.Out

	Models []any `group:"modelsTables"`
}

var AutoMigrate = fx.Module(
	"AutoMigrate",
	fx.Invoke(
		fx.Annotate(
			autoMigrate,
			fx.ParamTags(`group:"modelsTables"`),
		),
	),
)

func GetCRUDServiceModule[T any]() fx.Option {
	entity := *new(T)

	moduleName := fmt.Sprintf(
		"%TCRUDServiceModule",
		entity,
	)

	kind := getModelKind(entity)
	switch kind {
	case KindUUIDModel:
		return fx.Module(
			moduleName,
			// repository
			fx.Provide(NewCRUDRepository[T, UUID]),
			// service
			fx.Provide(NewCRUDService[T, UUID]),
		)
	case KindUIntModel:
		return fx.Module(
			moduleName,
			// repository
			fx.Provide(NewCRUDRepository[T, uint]),
			// service
			fx.Provide(NewCRUDService[T, uint]),
		)
	default:
		log.Printf("type %T is not a BaDaaS model\n", entity)
		return fx.Invoke(failNotBaDaaSModel())
	}
}

func failNotBaDaaSModel() error {
	return fmt.Errorf("type is not a BaDaaS model")
}

type modelKind uint

const (
	KindUUIDModel modelKind = iota
	KindUIntModel
	KindNotModel
)

func getModelKind(entity any) modelKind {
	entityType := getEntityType(entity)

	_, isUUIDModel := entityType.FieldByName("UUIDModel")
	if isUUIDModel {
		return KindUUIDModel
	}

	_, isUIntModel := entityType.FieldByName("UIntModel")
	if isUIntModel {
		return KindUIntModel
	}

	return KindNotModel
}

// Get the reflect.Type of any entity or pointer to entity
func getEntityType(entity any) reflect.Type {
	entityType := reflect.TypeOf(entity)

	// entityType will be a pointer if the relation can be nullable
	if entityType.Kind() == reflect.Pointer {
		entityType = entityType.Elem()
	}

	return entityType
}
