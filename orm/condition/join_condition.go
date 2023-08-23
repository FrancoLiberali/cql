package condition

import (
	"fmt"

	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
)

// Condition that joins T with any other model
type JoinCondition[T model.Model] interface {
	Condition[T]

	// Returns true if this condition or any nested condition makes a preload
	makesPreload() bool

	// Returns true if the condition of nay nested condition applies a filter (has where conditions)
	makesFilter() bool
}

// Condition that joins T with any other model
func NewJoinCondition[T1, T2 model.Model](
	conditions []Condition[T2],
	relationField string,
	t1Field string,
	t1PreloadCondition Condition[T1],
	t2Field string,
) JoinCondition[T1] {
	return joinConditionImpl[T1, T2]{
		Conditions:         conditions,
		RelationField:      relationField,
		T1Field:            t1Field,
		T1PreloadCondition: t1PreloadCondition,
		T2Field:            t2Field,
	}
}

// Implementation of join condition
type joinConditionImpl[T1, T2 model.Model] struct {
	T1Field       string
	T2Field       string
	RelationField string
	Conditions    []Condition[T2]
	// condition to preload T1 in case T2 any nested object is preloaded by user
	T1PreloadCondition Condition[T1]
}

func (condition joinConditionImpl[T1, T2]) InterfaceVerificationMethod(_ T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns true if this condition or any nested condition makes a preload
func (condition joinConditionImpl[T1, T2]) makesPreload() bool {
	_, joinConditions, t2PreloadCondition := divideConditionsByType(condition.Conditions)

	return t2PreloadCondition != nil || pie.Any(joinConditions, func(cond JoinCondition[T2]) bool {
		return cond.makesPreload()
	})
}

// Returns true if the condition of nay nested condition applies a filter (has where conditions)
//
//nolint:unused // is used
func (condition joinConditionImpl[T1, T2]) makesFilter() bool {
	whereConditions, joinConditions, _ := divideConditionsByType(condition.Conditions)

	return len(whereConditions) != 0 || pie.Any(joinConditions, func(cond JoinCondition[T2]) bool {
		return cond.makesFilter()
	})
}

// Applies a join between the tables of T1 and T2
// previousTableName is the name of the table of T1
// It also applies the nested conditions
func (condition joinConditionImpl[T1, T2]) ApplyTo(query *query.GormQuery, t1Table query.Table) error {
	whereConditions, joinConditions, t2PreloadCondition := divideConditionsByType(condition.Conditions)

	t2Model := *new(T2)

	// get the sql to do the join with T2
	t2Table, err := t1Table.DeliverTable(query, t2Model, condition.RelationField)
	if err != nil {
		return err
	}

	makesPreload := condition.makesPreload()
	joinQuery := condition.getSQLJoin(
		query,
		t1Table,
		t2Table,
		len(whereConditions) == 0 && makesPreload,
	)

	query.AddConcernedModel(
		t2Model,
		t2Table,
	)

	// apply WhereConditions to the join in the "on" clause
	connectionCondition := And(whereConditions...)

	onQuery, onValues, err := connectionCondition.GetSQL(query, t2Table)
	if err != nil {
		return err
	}

	if onQuery != "" {
		joinQuery += " AND " + onQuery
	}

	if !connectionCondition.AffectsDeletedAt() {
		joinQuery += fmt.Sprintf(
			" AND %s.deleted_at IS NULL",
			t2Table.Alias,
		)
	}

	// add the join to the query
	query.Joins(joinQuery, onValues...)

	// apply T1 preload condition
	// if this condition has a T2 preload condition
	// or any nested join condition has a preload condition
	// and this is not first level (T1 is the type of the repository)
	// because T1 is always loaded in that case
	if makesPreload && !t1Table.IsInitial() {
		err = condition.T1PreloadCondition.ApplyTo(query, t1Table)
		if err != nil {
			return err
		}
	}

	// apply T2 preload condition
	if t2PreloadCondition != nil {
		err = t2PreloadCondition.ApplyTo(query, t2Table)
		if err != nil {
			return err
		}
	}

	// apply nested joins
	for _, joinCondition := range joinConditions {
		err = joinCondition.ApplyTo(query, t2Table)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns the SQL string to do a join between T1 and T2
// taking into account that the ID attribute necessary to do it
// can be either in T1's or T2's table.
func (condition joinConditionImpl[T1, T2]) getSQLJoin(
	query *query.GormQuery,
	t1Table query.Table,
	t2Table query.Table,
	isLeftJoin bool,
) string {
	joinString := "INNER JOIN"
	if isLeftJoin {
		joinString = "LEFT JOIN"
	}

	return fmt.Sprintf(
		`%[6]s %[1]s %[2]s ON %[2]s.%[3]s = %[4]s.%[5]s
		`,
		t2Table.Name,
		t2Table.Alias,
		query.ColumnName(t2Table, condition.T2Field),
		t1Table.Alias,
		query.ColumnName(t1Table, condition.T1Field),
		joinString,
	)
}

// Divides a list of conditions by its type: WhereConditions and JoinConditions
func divideConditionsByType[T model.Model](
	conditions []Condition[T],
) (whereConditions []WhereCondition[T], joinConditions []JoinCondition[T], preload *preloadCondition[T]) {
	for _, condition := range conditions {
		possibleWhereCondition, ok := condition.(WhereCondition[T])
		if ok {
			whereConditions = append(whereConditions, possibleWhereCondition)
			continue
		}

		possiblePreloadCondition, ok := condition.(preloadCondition[T])
		if ok {
			preload = &possiblePreloadCondition
			continue
		}

		possibleJoinCondition, ok := condition.(JoinCondition[T])
		if ok {
			joinConditions = append(joinConditions, possibleJoinCondition)
			continue
		}
	}

	return
}
