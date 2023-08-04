package orm

import (
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/sql"
)

const deletedAtField = "DeletedAt"

var (
	IDFieldID        = FieldIdentifier{Field: "ID"}
	CreatedAtFieldID = FieldIdentifier{Field: "CreatedAt"}
	UpdatedAtFieldID = FieldIdentifier{Field: "UpdatedAt"}
	DeletedAtFieldID = FieldIdentifier{Field: deletedAtField}
)

type Condition[T Model] interface {
	// Applies the condition to the "query"
	// using the "tableName" as name for the table holding
	// the data for object of type T
	ApplyTo(query *Query, table Table) error

	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T],
	// since if no method receives by parameter a type T,
	// any other Condition[T2] would also be considered a Condition[T].
	InterfaceVerificationMethod(T)
}

// Conditions that can be used in a where clause
// (or in a on of a join)
type WhereCondition[T Model] interface {
	Condition[T]

	// Get the sql string and values to use in the query
	GetSQL(query *Query, table Table) (string, []any, error)

	// Returns true if the DeletedAt column if affected by the condition
	// If no condition affects the DeletedAt, the verification that it's null will be added automatically
	AffectsDeletedAt() bool
}

// Condition that contains a internal condition.
// Example: NOT (internal condition)
type ContainerCondition[T Model] struct {
	ConnectionCondition WhereCondition[T]
	Prefix              sql.Operator
}

func (condition ContainerCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition ContainerCondition[T]) ApplyTo(query *Query, table Table) error {
	return ApplyWhereCondition[T](condition, query, table)
}

func (condition ContainerCondition[T]) GetSQL(query *Query, table Table) (string, []any, error) {
	sqlString, values, err := condition.ConnectionCondition.GetSQL(query, table)
	if err != nil {
		return "", nil, err
	}

	sqlString = condition.Prefix.String() + " (" + sqlString + ")"

	return sqlString, values, nil
}

func (condition ContainerCondition[T]) AffectsDeletedAt() bool {
	return condition.ConnectionCondition.AffectsDeletedAt()
}

// Condition that contains a internal condition.
// Example: NOT (internal condition)
func NewContainerCondition[T Model](prefix sql.Operator, conditions ...WhereCondition[T]) WhereCondition[T] {
	if len(conditions) == 0 {
		return NewInvalidCondition[T](emptyConditionsError[T](prefix))
	}

	return ContainerCondition[T]{
		Prefix:              prefix,
		ConnectionCondition: And(conditions...),
	}
}

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
type ConnectionCondition[T Model] struct {
	Connector  sql.Operator
	Conditions []WhereCondition[T]
}

func (condition ConnectionCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition ConnectionCondition[T]) ApplyTo(query *Query, table Table) error {
	return ApplyWhereCondition[T](condition, query, table)
}

func (condition ConnectionCondition[T]) GetSQL(query *Query, table Table) (string, []any, error) {
	sqlStrings := []string{}
	values := []any{}

	for _, internalCondition := range condition.Conditions {
		internalSQLString, internalValues, err := internalCondition.GetSQL(query, table)
		if err != nil {
			return "", nil, err
		}

		sqlStrings = append(sqlStrings, internalSQLString)

		values = append(values, internalValues...)
	}

	return strings.Join(
		sqlStrings,
		" "+condition.Connector.String()+" ",
	), values, nil
}

func (condition ConnectionCondition[T]) AffectsDeletedAt() bool {
	return pie.Any(condition.Conditions, func(internalCondition WhereCondition[T]) bool {
		return internalCondition.AffectsDeletedAt()
	})
}

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
func NewConnectionCondition[T Model](connector sql.Operator, conditions ...WhereCondition[T]) WhereCondition[T] {
	return ConnectionCondition[T]{
		Connector:  connector,
		Conditions: conditions,
	}
}

type FieldIdentifier struct {
	Column       string
	Field        string
	ColumnPrefix string
}

// Returns the name of the column in which the field is saved in the table
func (fieldID FieldIdentifier) ColumnName(query *Query, table Table) string {
	columnName := fieldID.Column
	if columnName == "" {
		columnName = query.ColumnName(table, fieldID.Field)
	}

	// add column prefix and table name once we know the column name
	return fieldID.ColumnPrefix + columnName
}

// Returns the SQL to get the value of the field in the table
func (fieldID FieldIdentifier) ColumnSQL(query *Query, table Table) string {
	return table.Alias + "." + fieldID.ColumnName(query, table)
}

// Condition used to the preload the attributes of a model
type PreloadCondition[T Model] struct {
	Fields []FieldIdentifier
}

func (condition PreloadCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition PreloadCondition[T]) ApplyTo(query *Query, table Table) error {
	for _, fieldID := range condition.Fields {
		query.AddSelect(table, fieldID)
	}

	return nil
}

// Condition used to the preload the attributes of a model
func NewPreloadCondition[T Model](fields ...FieldIdentifier) PreloadCondition[T] {
	return PreloadCondition[T]{
		Fields: append(
			fields,
			// base model fields
			IDFieldID,
			CreatedAtFieldID,
			UpdatedAtFieldID,
			DeletedAtFieldID,
		),
	}
}

// Condition used to the preload a collection of models of a model
type CollectionPreloadCondition[T1 Model, T2 Model] struct {
	CollectionField string
	NestedPreloads  []IJoinCondition[T2]
}

func (condition CollectionPreloadCondition[T1, T2]) InterfaceVerificationMethod(_ T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T1]
}

func (condition CollectionPreloadCondition[T1, T2]) ApplyTo(query *Query, _ Table) error {
	if len(condition.NestedPreloads) == 0 {
		query.Preload(condition.CollectionField)
		return nil
	}

	query.Preload(
		condition.CollectionField,
		func(db *gorm.DB) *gorm.DB {
			preloadsAsCondition := pie.Map(
				condition.NestedPreloads,
				func(joinCondition IJoinCondition[T2]) Condition[T2] {
					return joinCondition
				},
			)

			preloadInternalQuery, err := NewQuery(db, preloadsAsCondition)
			if err != nil {
				_ = db.AddError(err)
				return db
			}

			return preloadInternalQuery.gormDB
		},
	)

	return nil
}

// Condition used to the preload a collection of models of a model
func NewCollectionPreloadCondition[T1 Model, T2 Model](collectionField string, nestedPreloads []IJoinCondition[T2]) Condition[T1] {
	if pie.Any(nestedPreloads, func(nestedPreload IJoinCondition[T2]) bool {
		return !nestedPreload.makesPreload() || nestedPreload.makesFilter()
	}) {
		return NewInvalidCondition[T1](onlyPreloadsAllowedError[T1](collectionField))
	}

	return CollectionPreloadCondition[T1, T2]{
		CollectionField: collectionField,
		NestedPreloads:  nestedPreloads,
	}
}

// Condition that verifies the value of a field,
// using the Operator
type FieldCondition[TObject Model, TAtribute any] struct {
	FieldIdentifier FieldIdentifier
	Operator        Operator[TAtribute]
}

func (condition FieldCondition[TObject, TAtribute]) InterfaceVerificationMethod(_ TObject) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
func (condition FieldCondition[TObject, TAtribute]) ApplyTo(query *Query, table Table) error {
	return ApplyWhereCondition[TObject](condition, query, table)
}

func ApplyWhereCondition[T Model](condition WhereCondition[T], query *Query, table Table) error {
	sql, values, err := condition.GetSQL(query, table)
	if err != nil {
		return err
	}

	if condition.AffectsDeletedAt() {
		query.Unscoped()
	}

	query.Where(
		sql,
		values...,
	)

	return nil
}

func (condition FieldCondition[TObject, TAtribute]) AffectsDeletedAt() bool {
	return condition.FieldIdentifier.Field == deletedAtField
}

func (condition FieldCondition[TObject, TAtribute]) GetSQL(query *Query, table Table) (string, []any, error) {
	sqlString, values, err := condition.Operator.ToSQL(
		condition.FieldIdentifier.ColumnSQL(query, table),
	)
	if err != nil {
		return "", nil, conditionOperatorError[TObject](err, condition)
	}

	return sqlString, values, nil
}

// Interface of a join condition that joins T with any other model
type IJoinCondition[T Model] interface {
	Condition[T]

	// Returns true if this condition or any nested condition makes a preload
	makesPreload() bool

	// Returns true if the condition of nay nested condition applies a filter (has where conditions)
	makesFilter() bool
}

// Condition that joins with other table
type JoinCondition[T1 Model, T2 Model] struct {
	T1Field       string
	T2Field       string
	RelationField string
	Conditions    []Condition[T2]
	// condition to preload T1 in case T2 any nested object is preloaded by user
	T1PreloadCondition PreloadCondition[T1]
}

func (condition JoinCondition[T1, T2]) InterfaceVerificationMethod(_ T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns true if this condition or any nested condition makes a preload
func (condition JoinCondition[T1, T2]) makesPreload() bool {
	_, joinConditions, t2PreloadCondition := divideConditionsByType(condition.Conditions)

	return t2PreloadCondition != nil || pie.Any(joinConditions, func(cond IJoinCondition[T2]) bool {
		return cond.makesPreload()
	})
}

// Returns true if the condition of nay nested condition applies a filter (has where conditions)
//
//nolint:unused // is used
func (condition JoinCondition[T1, T2]) makesFilter() bool {
	whereConditions, joinConditions, _ := divideConditionsByType(condition.Conditions)

	return len(whereConditions) != 0 || pie.Any(joinConditions, func(cond IJoinCondition[T2]) bool {
		return cond.makesFilter()
	})
}

// Applies a join between the tables of T1 and T2
// previousTableName is the name of the table of T1
// It also applies the nested conditions
func (condition JoinCondition[T1, T2]) ApplyTo(query *Query, t1Table Table) error {
	whereConditions, joinConditions, t2PreloadCondition := divideConditionsByType(condition.Conditions)

	// get the sql to do the join with T2
	t2Table, err := t1Table.DeliverTable(query, *new(T2), condition.RelationField)
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
func (condition JoinCondition[T1, T2]) getSQLJoin(
	query *Query,
	t1Table Table,
	t2Table Table,
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
func divideConditionsByType[T Model](
	conditions []Condition[T],
) (whereConditions []WhereCondition[T], joinConditions []IJoinCondition[T], preloadCondition *PreloadCondition[T]) {
	for _, condition := range conditions {
		possibleWhereCondition, ok := condition.(WhereCondition[T])
		if ok {
			whereConditions = append(whereConditions, possibleWhereCondition)
			continue
		}

		possiblePreloadCondition, ok := condition.(PreloadCondition[T])
		if ok {
			preloadCondition = &possiblePreloadCondition
			continue
		}

		possibleJoinCondition, ok := condition.(IJoinCondition[T])
		if ok {
			joinConditions = append(joinConditions, possibleJoinCondition)
			continue
		}
	}

	return
}

// Condition used to returns an error when the query is executed
type InvalidCondition[T any] struct {
	Err error
}

func (condition InvalidCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition InvalidCondition[T]) ApplyTo(_ *Query, _ Table) error {
	return condition.Err
}

func (condition InvalidCondition[T]) GetSQL(_ *Query, _ Table) (string, []any, error) {
	return "", nil, condition.Err
}

func (condition InvalidCondition[T]) AffectsDeletedAt() bool {
	return false
}

// Condition used to returns an error when the query is executed
func NewInvalidCondition[T any](err error) InvalidCondition[T] {
	return InvalidCondition[T]{
		Err: err,
	}
}

// Logical Operators
// ref: https://www.postgresql.org/docs/current/functions-logical.html

func And[T Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition(sql.And, conditions...)
}

func Or[T Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition(sql.Or, conditions...)
}

func Not[T Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewContainerCondition(sql.Not, conditions...)
}
