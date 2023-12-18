package condition

import "github.com/ditrit/badaas/orm/model"

type DynamicCondition[T model.Model] interface {
	WhereCondition[T]

	// Allows to choose which number of join use
	// for the operation in position "operationNumber"
	// when the value is a field and its model is joined more than once.
	// Does nothing if the operationNumber is bigger than the amount of operations.
	SelectJoin(operationNumber, joinNumber uint) DynamicCondition[T]
}
