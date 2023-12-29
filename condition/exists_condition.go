package condition

import (
	"fmt"

	"github.com/FrancoLiberali/cql/model"
)

// Condition that generates a WHERE EXISTS
type existsCondition[T1 model.Model, T2 model.Model] struct {
	Conditions    []WhereCondition[T2]
	RelationField string
	T1Field       string
	T2Field       string
}

//nolint:unused // is used
func (condition existsCondition[T1, T2]) interfaceVerificationMethod(_ T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

//nolint:unused // is used
func (condition existsCondition[T1, T2]) applyTo(query *GormQuery, table Table) error {
	return ApplyWhereCondition[T1](condition, query, table)
}

//nolint:unused // is used
func (condition existsCondition[T1, T2]) getSQL(query *GormQuery, t1Table Table) (string, []any, error) {
	connectionCondition := And(condition.Conditions...)

	t2Table, err := t1Table.DeliverTable(query, *new(T2), condition.RelationField)
	if err != nil {
		return "", nil, err
	}

	sql, values, err := connectionCondition.getSQL(query, t2Table)
	if err != nil {
		return "", nil, err
	}

	deletedAtSQL := ""
	if !connectionCondition.affectsDeletedAt() {
		deletedAtSQL = fmt.Sprintf(
			"AND %s.deleted_at IS NULL",
			t2Table.Alias,
		)
	}

	return fmt.Sprintf(
		"EXISTS (SELECT(1) FROM %s %s WHERE %s AND %s %s)",
		t2Table.Name,
		t2Table.Alias,
		getSQLJoin(query, t1Table, condition.T1Field, t2Table, condition.T2Field),
		sql,
		deletedAtSQL,
	), values, nil
}

//nolint:unused // is used
func (condition existsCondition[T1, T2]) affectsDeletedAt() bool {
	return false
}
