package orm

import (
	"fmt"
	"log"
	"reflect"

	"go.uber.org/fx"

	"github.com/ditrit/badaas/orm/model"
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

func GetCRUDServiceModule[T model.Model]() fx.Option {
	entity := *new(T)

	moduleName := fmt.Sprintf(
		"%TCRUDServiceModule",
		entity,
	)

	kind := getModelKind(entity)
	switch kind {
	case kindUUIDModel:
		return fx.Module(
			moduleName,
			// repository
			fx.Provide(NewCRUDRepository[T, model.UUID]),
			// service
			fx.Provide(NewCRUDService[T, model.UUID]),
		)
	case kindUIntModel:
		return fx.Module(
			moduleName,
			// repository
			fx.Provide(NewCRUDRepository[T, model.UIntID]),
			// service
			fx.Provide(NewCRUDService[T, model.UIntID]),
		)
	case kindNotModel:
		log.Printf("type %T is not a BaDaaS model\n", entity)
		return fx.Invoke(failNotBaDaaSModel())
	}

	return nil
}

func failNotBaDaaSModel() error {
	return fmt.Errorf("type is not a BaDaaS model")
}

type modelKind uint

const (
	kindUUIDModel modelKind = iota
	kindUIntModel
	kindNotModel
)

func getModelKind(entity model.Model) modelKind {
	entityType := getEntityType(entity)

	_, isUUIDModel := entityType.FieldByName("UUIDModel")
	if isUUIDModel {
		return kindUUIDModel
	}

	_, isUIntModel := entityType.FieldByName("UIntModel")
	if isUIntModel {
		return kindUIntModel
	}

	return kindNotModel
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
