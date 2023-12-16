package orm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

const DeletedAtField = "DeletedAt"

var ErrEmptyConditions = errors.New("condition must have at least one inner condition")

type Condition[T any] interface {
	// Applies the condition to the "query"
	// using the "tableName" as name for the table holding
	// the data for object of type T
	ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error)

	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T],
	// since if no method receives by parameter a type T,
	// any other Condition[T2] would also be considered a Condition[T].
	interfaceVerificationMethod(T)
}

// Conditions that can be used in a where clause
// (or in a on of a join)
type WhereCondition[T any] interface {
	Condition[T]

	// Get the sql string and values to use in the query
	GetSQL(query *gorm.DB, tableName string) (string, []any, error)

	// Returns true if the DeletedAt column if affected by the condition
	// If no condition affects the DeletedAt, the verification that it's null will be added automatically
	affectsDeletedAt() bool
}

// Condition that contains a internal condition.
// Example: NOT (internal condition)
type ContainerCondition[T any] struct {
	ConnectionCondition WhereCondition[T]
	Prefix              string
}

//nolint:unused // see inside
func (condition ContainerCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition ContainerCondition[T]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	return applyWhereCondition[T](condition, query, tableName)
}

func (condition ContainerCondition[T]) GetSQL(query *gorm.DB, tableName string) (string, []any, error) {
	sqlString, values, err := condition.ConnectionCondition.GetSQL(query, tableName)
	if err != nil {
		return "", nil, err
	}

	sqlString = condition.Prefix + " (" + sqlString + ")"

	return sqlString, values, nil
}

//nolint:unused // is used
func (condition ContainerCondition[T]) affectsDeletedAt() bool {
	return condition.ConnectionCondition.affectsDeletedAt()
}

// Condition that contains a internal condition.
// Example: NOT (internal condition)
func NewContainerCondition[T any](prefix string, conditions ...WhereCondition[T]) WhereCondition[T] {
	if len(conditions) == 0 {
		return NewInvalidCondition[T](ErrEmptyConditions)
	}

	return ContainerCondition[T]{
		Prefix:              prefix,
		ConnectionCondition: And(conditions...),
	}
}

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
type ConnectionCondition[T any] struct {
	Connector  string
	Conditions []WhereCondition[T]
}

//nolint:unused // see inside
func (condition ConnectionCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition ConnectionCondition[T]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	return applyWhereCondition[T](condition, query, tableName)
}

func (condition ConnectionCondition[T]) GetSQL(query *gorm.DB, tableName string) (string, []any, error) {
	sqlStrings := []string{}
	values := []any{}

	for _, internalCondition := range condition.Conditions {
		internalSQLString, internalValues, err := internalCondition.GetSQL(query, tableName)
		if err != nil {
			return "", nil, err
		}

		sqlStrings = append(sqlStrings, internalSQLString)

		values = append(values, internalValues...)
	}

	return strings.Join(sqlStrings, " "+condition.Connector+" "), values, nil
}

//nolint:unused // is used
func (condition ConnectionCondition[T]) affectsDeletedAt() bool {
	return pie.Any(condition.Conditions, func(internalCondition WhereCondition[T]) bool {
		return internalCondition.affectsDeletedAt()
	})
}

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
func NewConnectionCondition[T any](connector string, conditions ...WhereCondition[T]) WhereCondition[T] {
	return ConnectionCondition[T]{
		Connector:  connector,
		Conditions: conditions,
	}
}

// Condition that verifies the value of a field,
// using the Operator
type FieldCondition[TObject any, TAtribute any] struct {
	Field        string
	Column       string
	ColumnPrefix string
	Operator     Operator[TAtribute]
}

//nolint:unused // see inside
func (condition FieldCondition[TObject, TAtribute]) interfaceVerificationMethod(_ TObject) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
func (condition FieldCondition[TObject, TAtribute]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	return applyWhereCondition[TObject](condition, query, tableName)
}

func applyWhereCondition[T any](condition WhereCondition[T], query *gorm.DB, tableName string) (*gorm.DB, error) {
	sql, values, err := condition.GetSQL(query, tableName)
	if err != nil {
		return nil, err
	}

	if condition.affectsDeletedAt() {
		query = query.Unscoped()
	}

	return query.Where(
		sql,
		values...,
	), nil
}

//nolint:unused // is used
func (condition FieldCondition[TObject, TAtribute]) affectsDeletedAt() bool {
	return condition.Field == DeletedAtField
}

func (condition FieldCondition[TObject, TAtribute]) GetSQL(query *gorm.DB, tableName string) (string, []any, error) {
	columnName := condition.Column
	if columnName == "" {
		columnName = query.NamingStrategy.ColumnName(tableName, condition.Field)
	}

	// add column prefix and table name once we know the column name
	columnName = tableName + "." + condition.ColumnPrefix + columnName

	return condition.Operator.ToSQL(columnName)
}

// Condition that joins with other table
type JoinCondition[T1 any, T2 any] struct {
	T1Field    string
	T2Field    string
	Conditions []Condition[T2]
}

func (condition JoinCondition[T1, T2]) interfaceVerificationMethod(t T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Applies a join between the tables of T1 and T2
// previousTableName is the name of the table of T1
// It also applies the nested conditions
func (condition JoinCondition[T1, T2]) ApplyTo(query *gorm.DB, previousTableName string) (*gorm.DB, error) {
	// get the name of the table for T2
	toBeJoinedTableName, err := getTableName(query, *new(T2))
	if err != nil {
		return nil, err
	}

	// add a suffix to avoid tables with the same name when joining
	// the same table more than once
	nextTableName := toBeJoinedTableName + "_" + previousTableName

	// get the sql to do the join with T2
	joinQuery := condition.getSQLJoin(query, toBeJoinedTableName, nextTableName, previousTableName)

	whereConditions, joinConditions := divideConditionsByType(condition.Conditions)

	// apply WhereConditions to join in "on" clause
	connectionCondition := And(whereConditions...)

	onQuery, onValues, err := connectionCondition.GetSQL(query, nextTableName)
	if err != nil {
		return nil, err
	}

	if onQuery != "" {
		joinQuery += " AND " + onQuery
	}

	if !connectionCondition.affectsDeletedAt() {
		joinQuery += fmt.Sprintf(
			" AND %s.deleted_at IS NULL",
			nextTableName,
		)
	}

	// add the join to the query
	query = query.Joins(joinQuery, onValues...)

	// apply nested joins
	for _, joinCondition := range joinConditions {
		query, err = joinCondition.ApplyTo(query, nextTableName)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}

// Returns the SQL string to do a join between T1 and T2
// taking into account that the ID attribute necessary to do it
// can be either in T1's or T2's table.
func (condition JoinCondition[T1, T2]) getSQLJoin(query *gorm.DB, toBeJoinedTableName, nextTableName, previousTableName string) string {
	return fmt.Sprintf(
		`JOIN %[1]s %[2]s ON %[2]s.%[3]s = %[4]s.%[5]s
		`,
		toBeJoinedTableName,
		nextTableName,
		query.NamingStrategy.ColumnName(nextTableName, condition.T2Field),
		previousTableName,
		query.NamingStrategy.ColumnName(previousTableName, condition.T1Field),
	)
}

// Divides a list of conditions by its type: WhereConditions and JoinConditions
func divideConditionsByType[T any](
	conditions []Condition[T],
) (thisEntityConditions []WhereCondition[T], joinConditions []Condition[T]) {
	for _, condition := range conditions {
		typedCondition, ok := condition.(WhereCondition[T])
		if ok {
			thisEntityConditions = append(thisEntityConditions, typedCondition)
		} else {
			joinConditions = append(joinConditions, condition)
		}
	}

	return
}

// Condition that can be used to express conditions that are not supported (yet?) by BaDORM
// Example: table1.columnX = table2.columnY
type UnsafeCondition[T any] struct {
	SQLCondition string
	Values       []any
}

//nolint:unused // see inside
func (condition UnsafeCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition UnsafeCondition[T]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	return applyWhereCondition[T](condition, query, tableName)
}

func (condition UnsafeCondition[T]) GetSQL(_ *gorm.DB, tableName string) (string, []any, error) {
	return fmt.Sprintf(
		condition.SQLCondition,
		tableName,
	), condition.Values, nil
}

//nolint:unused // is used
func (condition UnsafeCondition[T]) affectsDeletedAt() bool {
	return false
}

// Condition that can be used to express conditions that are not supported (yet?) by BaDORM
// Example: table1.columnX = table2.columnY
func NewUnsafeCondition[T any](condition string, values []any) UnsafeCondition[T] {
	return UnsafeCondition[T]{
		SQLCondition: condition,
		Values:       values,
	}
}

// Condition used to returns an error when the query is executed
type InvalidCondition[T any] struct {
	Err error
}

//nolint:unused // see inside
func (condition InvalidCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition InvalidCondition[T]) ApplyTo(_ *gorm.DB, _ string) (*gorm.DB, error) {
	return nil, condition.Err
}

func (condition InvalidCondition[T]) GetSQL(_ *gorm.DB, _ string) (string, []any, error) {
	return "", nil, condition.Err
}

//nolint:unused // is used
func (condition InvalidCondition[T]) affectsDeletedAt() bool {
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

func And[T any](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition("AND", conditions...)
}

func Or[T any](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition("OR", conditions...)
}

func Not[T any](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewContainerCondition("NOT", conditions...)
}
