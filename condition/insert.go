package condition

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

type Insert[T model.Model] struct {
	tx     *gorm.DB
	query  *CQLQuery
	err    error
	models []*T
}

func NewInsert[T model.Model](tx *gorm.DB, models []*T) *Insert[T] {
	return &Insert[T]{
		tx:     tx,
		query:  NewQuery[T](tx).cqlQuery,
		models: models,
	}
}

// OnConflict allows to set the action to be taken when any conflict happens when inserting the data.
//
// WARNING: in postgres OnConflict can be used only with DoNothing,
// for UpdateAll, Update and Set, OnConflictOn must be used
func (insert *Insert[T]) OnConflict() *InsertOnConflict[T] {
	return &InsertOnConflict[T]{
		insert: insert,
	}
}

// OnConflictOn allows to set the action to be taken when a conflict
// with the fields specified by parameter happens when inserting the data.
// For this, there must be a constraint on these fields, otherwise an error will be returned.
//
// Available for: postgres, sqlite
func (insert *Insert[T]) OnConflictOn(field FieldOfModel[T], fields ...FieldOfModel[T]) *InsertOnConflict[T] {
	if insert.query.Dialector() == sql.MySQL || insert.query.Dialector() == sql.SQLServer {
		insert.err = methodError(ErrUnsupportedByDatabase, "OnConflictOn")

		return &InsertOnConflict[T]{
			insert: insert,
		}
	}

	onConflictColumns := pie.Map(
		insert.getFieldNames(pie.Unshift(fields, field)),
		func(fieldName string) clause.Column {
			return clause.Column{Name: fieldName}
		},
	)

	return &InsertOnConflict[T]{
		insert:            insert,
		onConflictColumns: onConflictColumns,
	}
}

func (insert *Insert[T]) getFieldNames(fields []FieldOfModel[T]) []string {
	fieldNames := make([]string, 0, len(fields))

	for _, field := range fields {
		fieldNames = append(
			fieldNames,
			field.columnName(insert.query, insert.query.initialTable),
		)
	}

	return fieldNames
}

// OnConstraint allows to set the action to be taken when a conflict
// with the constraint specified by parameter happens when inserting the data.
// For this, the constraint must exist, otherwise an error will be returned.
//
// Available for: postgres
func (insert *Insert[T]) OnConstraint(constraintName string) *InsertOnConflict[T] {
	if insert.query.Dialector() != sql.Postgres {
		insert.err = methodError(ErrUnsupportedByDatabase, "OnConstraint")

		return &InsertOnConflict[T]{
			insert: insert,
		}
	}

	return &InsertOnConflict[T]{
		insert:       insert,
		onConstraint: constraintName,
	}
}

func (insert *Insert[T]) addOnConflictClause(clause clause.OnConflict) {
	insert.tx = insert.tx.Clauses(clause)
}

type InsertOnConflict[T model.Model] struct {
	insert *Insert[T]

	onConflictColumns []clause.Column
	onConstraint      string
}

// DoNothing will not take any action, simply preventing an error from being responded.
func (insertOnConflict *InsertOnConflict[T]) DoNothing() *InsertExec[T] {
	insertOnConflict.insert.addOnConflictClause(clause.OnConflict{
		OnConstraint: insertOnConflict.onConstraint,
		Columns:      insertOnConflict.onConflictColumns,
		DoNothing:    true,
	})

	return &InsertExec[T]{insert: insertOnConflict.insert}
}

func (insertOnConflict *InsertOnConflict[T]) addPostgresErrorIfNotColumns(msg string) {
	if insertOnConflict.insert.query.Dialector() == sql.Postgres &&
		len(insertOnConflict.onConflictColumns) == 0 &&
		insertOnConflict.onConstraint == "" {
		insertOnConflict.insert.err = methodError(ErrUnsupportedByDatabase, msg)
	}
}

// UpdateAll will update all model attributes with the values of the models that already exist.
func (insertOnConflict *InsertOnConflict[T]) UpdateAll() *InsertExec[T] {
	insertOnConflict.addPostgresErrorIfNotColumns("UpdateAll after OnConflict")

	insertOnConflict.insert.addOnConflictClause(clause.OnConflict{
		OnConstraint: insertOnConflict.onConstraint,
		Columns:      insertOnConflict.onConflictColumns,
		UpdateAll:    true,
	})

	return &InsertExec[T]{insert: insertOnConflict.insert}
}

// Update will update the attributes specified by parameter with the values of the models that already exist.
func (insertOnConflict *InsertOnConflict[T]) Update(fields ...FieldOfModel[T]) *InsertExec[T] {
	insertOnConflict.addPostgresErrorIfNotColumns("Update after OnConflict")

	fieldNames := insertOnConflict.insert.getFieldNames(fields)

	insertOnConflict.insert.addOnConflictClause(clause.OnConflict{
		OnConstraint: insertOnConflict.onConstraint,
		Columns:      insertOnConflict.onConflictColumns,
		DoUpdates:    clause.AssignmentColumns(fieldNames),
	})

	return &InsertExec[T]{insert: insertOnConflict.insert}
}

// Set allows to specify which updates to perform.
func (insertOnConflict *InsertOnConflict[T]) Set(sets ...*Set[T]) *InsertOnConflictSet[T] {
	insertOnConflict.addPostgresErrorIfNotColumns("Set after OnConflict")

	return &InsertOnConflictSet[T]{
		insertOnConflict: insertOnConflict,
		sets:             sets,
	}
}

type InsertOnConflictSet[T model.Model] struct {
	insertOnConflict *InsertOnConflict[T]

	sets []*Set[T]
}

// Exec execute the insert statement, returning the amount of rows inserted.
// It will also update the inserted model's primary key in their ids
//
// WARNING: the value returned may depend on the db engine, for example mysql
// returns the double of the other ones when there is conflict
func (insertOnConflictSet *InsertOnConflictSet[T]) Exec() (int64, error) {
	return insertOnConflictSet.internalExec(func(insert *Insert[T]) (int64, error) {
		return insert.Exec()
	})
}

// ExecInBatches execute the insert statement in batches of batchSize,
// returning the amount of rows inserted.
// It will also update the inserted model's primary key in their ids
//
// WARNING: the value returned may depend on the db engine, for example mysql
// returns the double of the other ones when there is conflict
func (insertOnConflictSet *InsertOnConflictSet[T]) ExecInBatches(batchSize int) (int64, error) {
	return insertOnConflictSet.internalExec(func(insert *Insert[T]) (int64, error) {
		return insert.ExecInBatches(batchSize)
	})
}

func (insertOnConflictSet *InsertOnConflictSet[T]) internalExec(execFunc func(*Insert[T]) (int64, error)) (int64, error) {
	insert := insertOnConflictSet.insertOnConflict.insert

	onConflictClause, err := insertOnConflictSet.getOnConflictClause()
	if err != nil {
		return 0, err
	}

	insert.addOnConflictClause(onConflictClause)

	return execFunc(insert)
}

// Where allows to set conditions on the models that generate conflicts when performing the updates.
//
// Available for: postgres, sqlite
func (insertOnConflictSet *InsertOnConflictSet[T]) Where(conditions ...Condition[T]) *InsertExec[T] {
	insert := insertOnConflictSet.insertOnConflict.insert

	if insert.query.Dialector() == sql.MySQL || insert.query.Dialector() == sql.SQLServer {
		insert.err = methodError(ErrUnsupportedByDatabase, "Where")

		return &InsertExec[T]{insert: insert}
	}

	insert.query = NewQuery(insert.query.gormDB, conditions...).cqlQuery

	onConflictClause, err := insertOnConflictSet.getOnConflictClause()
	if err != nil {
		insert.err = err
	}

	where, isWhere := insert.query.gormDB.Statement.Clauses["WHERE"].Expression.(clause.Where)
	if isWhere {
		onConflictClause.Where = where
	}

	insert.addOnConflictClause(onConflictClause)

	return &InsertExec[T]{insert: insert}
}

func (insertOnConflictSet *InsertOnConflictSet[T]) getOnConflictClause() (clause.OnConflict, error) {
	assignments := map[string]any{}

	insert := insertOnConflictSet.insertOnConflict.insert

	for _, set := range insertOnConflictSet.sets {
		setSQL, setValues, err := set.getValue().ToSQL(insert.query)
		if err != nil {
			return clause.OnConflict{}, err
		}

		fieldColumn := set.getField().columnName(insert.query, insert.query.initialTable)

		if setSQL == "" {
			assignments[fieldColumn] = setValues[0]
		} else {
			assignments[fieldColumn] = gorm.Expr(
				setSQL,
				setValues...,
			)
		}
	}

	return clause.OnConflict{
		OnConstraint: insertOnConflictSet.insertOnConflict.onConstraint,
		Columns:      insertOnConflictSet.insertOnConflict.onConflictColumns,
		DoUpdates:    clause.Assignments(assignments),
	}, nil
}

type InsertExec[T model.Model] struct {
	insert *Insert[T]
}

// Exec execute the insert statement, returning the amount of rows inserted.
// It will also update the inserted model's primary key in their ids
//
// WARNING: the value returned may depend on the db engine, for example mysql
// returns the double of the other ones when there is conflict
func (insertExec *InsertExec[T]) Exec() (int64, error) {
	return insertExec.insert.Exec()
}

// ExecInBatches execute the insert statement in batches of batchSize,
// returning the amount of rows inserted.
// It will also update the inserted model's primary key in their ids
//
// WARNING: the value returned may depend on the db engine, for example mysql
// returns the double of the other ones when there is conflict
func (insertExec *InsertExec[T]) ExecInBatches(batchSize int) (int64, error) {
	return insertExec.insert.ExecInBatches(batchSize)
}

// Exec execute the insert statement, returning the amount of rows inserted.
// It will also update the inserted model's primary key in their ids
//
// WARNING: the value returned may depend on the db engine, for example mysql
// returns the double of the other ones when there is conflict
func (insert *Insert[T]) Exec() (int64, error) {
	if insert.err != nil {
		return 0, insert.err
	}

	result := insert.tx.Omit(clause.Associations).Create(insert.models)

	return result.RowsAffected, result.Error
}

// ExecInBatches execute the insert statement in batches of batchSize,
// returning the amount of rows inserted.
// It will also update the inserted model's primary key in their ids
//
// WARNING: the value returned may depend on the db engine, for example mysql
// returns the double of the other ones when there is conflict
func (insert *Insert[T]) ExecInBatches(batchSize int) (int64, error) {
	if insert.err != nil {
		return 0, insert.err
	}

	result := insert.tx.Omit(clause.Associations).CreateInBatches(insert.models, batchSize)

	return result.RowsAffected, result.Error
}
