package condition

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

type Condition[T model.Model] interface {
	// Applies the condition to the "query"
	// using the table holding
	// the data for object of type T
	applyTo(query *GormQuery, table Table) error

	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T],
	// since if no method receives by parameter a type T,
	// any other Condition[T2] would also be considered a Condition[T].
	interfaceVerificationMethod(t T)
}

// Create a GormQuery to which the conditions are applied
func ApplyConditions[T model.Model](db *gorm.DB, conditions []Condition[T]) (*GormQuery, error) {
	model := *new(T)

	initialTable, err := NewTable(db, model)
	if err != nil {
		return nil, err
	}

	query := NewGormQuery(db, model, initialTable)

	for _, condition := range conditions {
		err := condition.applyTo(query, initialTable)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}
