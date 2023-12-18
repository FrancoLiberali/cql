package condition

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
)

type Condition[T model.Model] interface {
	// Applies the condition to the "query"
	// using the table holding
	// the data for object of type T
	ApplyTo(*query.Query, query.Table) error

	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T],
	// since if no method receives by parameter a type T,
	// any other Condition[T2] would also be considered a Condition[T].
	InterfaceVerificationMethod(T)
}

// Create a Query to which the conditions are applied
func ApplyConditions[T model.Model](db *gorm.DB, conditions []Condition[T]) (*query.Query, error) {
	model := *new(T)

	initialTable, err := query.NewTable(db, model)
	if err != nil {
		return nil, err
	}

	query := query.NewQuery(db, model, initialTable)

	for _, condition := range conditions {
		err := condition.ApplyTo(query, initialTable)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}
