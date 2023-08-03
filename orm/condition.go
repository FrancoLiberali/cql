package orm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

const deletedAtField = "DeletedAt"

var (
	IDFieldID        = FieldIdentifier{Field: "ID"}
	CreatedAtFieldID = FieldIdentifier{Field: "CreatedAt"}
	UpdatedAtFieldID = FieldIdentifier{Field: "UpdatedAt"}
	DeletedAtFieldID = FieldIdentifier{Field: deletedAtField}
)

var ErrEmptyConditions = errors.New("condition must have at least one inner condition")

type Table struct {
	Name    string
	Alias   string
	Initial bool
}

// Returns true if the Table is the initial table in a query
func (table Table) IsInitial() bool {
	return table.Initial
}

// Returns the related Table corresponding to the model
func (table Table) DeliverTable(query *gorm.DB, model any, relationName string) (Table, error) {
	// get the name of the table for the model
	tableName, err := getTableName(query, model)
	if err != nil {
		return Table{}, err
	}

	// add a suffix to avoid tables with the same name when joining
	// the same table more than once
	tableAlias := relationName
	if !table.IsInitial() {
		tableAlias = table.Alias + "__" + relationName
	}

	return Table{
		Name:    tableName,
		Alias:   tableAlias,
		Initial: false,
	}, nil
}

type Condition[T any] interface {
	// Applies the condition to the "query"
	// using the "tableName" as name for the table holding
	// the data for object of type T
	ApplyTo(query *gorm.DB, table Table) (*gorm.DB, error)

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
	GetSQL(query *gorm.DB, table Table) (string, []any, error)

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

func (condition ContainerCondition[T]) ApplyTo(query *gorm.DB, table Table) (*gorm.DB, error) {
	return applyWhereCondition[T](condition, query, table)
}

func (condition ContainerCondition[T]) GetSQL(query *gorm.DB, table Table) (string, []any, error) {
	sqlString, values, err := condition.ConnectionCondition.GetSQL(query, table)
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

func (condition ConnectionCondition[T]) ApplyTo(query *gorm.DB, table Table) (*gorm.DB, error) {
	return applyWhereCondition[T](condition, query, table)
}

func (condition ConnectionCondition[T]) GetSQL(query *gorm.DB, table Table) (string, []any, error) {
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

type FieldIdentifier struct {
	Column       string
	Field        string
	ColumnPrefix string
}

func (columnID FieldIdentifier) ColumnName(db *gorm.DB, table Table) string {
	columnName := columnID.Column
	if columnName == "" {
		columnName = db.NamingStrategy.ColumnName(table.Name, columnID.Field)
	}

	// add column prefix and table name once we know the column name
	return columnID.ColumnPrefix + columnName
}

// Condition that verifies the value of a field,
// using the Operator
type FieldCondition[TObject any, TAtribute any] struct {
	FieldIdentifier FieldIdentifier
	Operator        Operator[TAtribute]
}

//nolint:unused // see inside
func (condition FieldCondition[TObject, TAtribute]) interfaceVerificationMethod(_ TObject) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
func (condition FieldCondition[TObject, TAtribute]) ApplyTo(query *gorm.DB, table Table) (*gorm.DB, error) {
	return applyWhereCondition[TObject](condition, query, table)
}

func applyWhereCondition[T any](condition WhereCondition[T], query *gorm.DB, table Table) (*gorm.DB, error) {
	sql, values, err := condition.GetSQL(query, table)
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
	return condition.FieldIdentifier.Field == deletedAtField
}

func (condition FieldCondition[TObject, TAtribute]) GetSQL(query *gorm.DB, table Table) (string, []any, error) {
	columnName := table.Alias + "." + condition.FieldIdentifier.ColumnName(query, table)
	return condition.Operator.ToSQL(columnName)
}

// Condition that joins with other table
type JoinCondition[T1 any, T2 any] struct {
	T1Field       string
	T2Field       string
	RelationField string
	Conditions    []Condition[T2]
}

func (condition JoinCondition[T1, T2]) interfaceVerificationMethod(t T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Applies a join between the tables of T1 and T2
// previousTableName is the name of the table of T1
// It also applies the nested conditions
func (condition JoinCondition[T1, T2]) ApplyTo(query *gorm.DB, t1Table Table) (*gorm.DB, error) {
	whereConditions, joinConditions := divideConditionsByType(condition.Conditions)

	// get the sql to do the join with T2
	t2Table, err := t1Table.DeliverTable(query, *new(T2), condition.RelationField)
	if err != nil {
		return nil, err
	}

	joinQuery := condition.getSQLJoin(
		query,
		t1Table,
		t2Table,
	)

	// apply WhereConditions to the join in the "on" clause
	connectionCondition := And(whereConditions...)

	onQuery, onValues, err := connectionCondition.GetSQL(query, t2Table)
	if err != nil {
		return nil, err
	}

	if onQuery != "" {
		joinQuery += " AND " + onQuery
	}

	if !connectionCondition.affectsDeletedAt() {
		joinQuery += fmt.Sprintf(
			" AND %s.deleted_at IS NULL",
			t2Table.Alias,
		)
	}

	// add the join to the query
	query = query.Joins(joinQuery, onValues...)

	// apply nested joins
	for _, joinCondition := range joinConditions {
		query, err = joinCondition.ApplyTo(query, t2Table)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}

// Returns the SQL string to do a join between T1 and T2
// taking into account that the ID attribute necessary to do it
// can be either in T1's or T2's table.
func (condition JoinCondition[T1, T2]) getSQLJoin(
	query *gorm.DB,
	t1Table Table,
	t2Table Table,
) string {
	return fmt.Sprintf(
		`JOIN %[1]s %[2]s ON %[2]s.%[3]s = %[4]s.%[5]s
		`,
		t2Table.Name,
		t2Table.Alias,
		query.NamingStrategy.ColumnName(t2Table.Name, condition.T2Field),
		t1Table.Alias,
		query.NamingStrategy.ColumnName(t1Table.Name, condition.T1Field),
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

func (condition UnsafeCondition[T]) ApplyTo(query *gorm.DB, table Table) (*gorm.DB, error) {
	return applyWhereCondition[T](condition, query, table)
}

func (condition UnsafeCondition[T]) GetSQL(_ *gorm.DB, table Table) (string, []any, error) {
	return fmt.Sprintf(
		condition.SQLCondition,
		table.Alias,
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

func (condition InvalidCondition[T]) ApplyTo(_ *gorm.DB, _ Table) (*gorm.DB, error) {
	return nil, condition.Err
}

func (condition InvalidCondition[T]) GetSQL(_ *gorm.DB, _ Table) (string, []any, error) {
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
