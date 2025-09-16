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

// WARNING: in postgres OnConflict can be used only with DoNothing,
// for UpdateAll, Update and Set, OnConflictOn must be used
func (insert *Insert[T]) OnConflict() *InsertOnConflict[T] {
	return &InsertOnConflict[T]{
		insert: insert,
	}
}

// TODO este ifield puede ser de cualquier model
// o necesita linter
// o hacer dos metodos uno para insertar solo un model y otro para insertar junto con las relations
// que ahi si necesitas Ifield. El linter lo necesitas de una forma o la otra
// Available for: postgres, sqlite
func (insert *Insert[T]) OnConflictOn(field IField, fields ...IField) *InsertOnConflict[T] {
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

func (insert *Insert[T]) getFieldNames(fields []IField) []string {
	fieldNames := make([]string, 0, len(fields))

	for _, field := range fields {
		fieldNames = append(
			fieldNames,
			field.columnName(insert.query, insert.query.initialTable),
		)
	}

	return fieldNames
}

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

func (insertOnConflict *InsertOnConflict[T]) UpdateAll() *InsertExec[T] {
	insertOnConflict.addPostgresErrorIfNotColumns("UpdateAll after OnConflict")

	insertOnConflict.insert.addOnConflictClause(clause.OnConflict{
		OnConstraint: insertOnConflict.onConstraint,
		Columns:      insertOnConflict.onConflictColumns,
		UpdateAll:    true,
	})

	return &InsertExec[T]{insert: insertOnConflict.insert}
}

func (insertOnConflict *InsertOnConflict[T]) Update(fields ...IField) *InsertExec[T] {
	insertOnConflict.addPostgresErrorIfNotColumns("Update after OnConflict")

	// TODO estos fields podrian ser de solo T
	fieldNames := insertOnConflict.insert.getFieldNames(fields)

	insertOnConflict.insert.addOnConflictClause(clause.OnConflict{
		OnConstraint: insertOnConflict.onConstraint,
		Columns:      insertOnConflict.onConflictColumns,
		DoUpdates:    clause.AssignmentColumns(fieldNames),
	})

	return &InsertExec[T]{insert: insertOnConflict.insert}
}

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

func (insertOnConflictSet *InsertOnConflictSet[T]) Exec() (int64, error) {
	insert := insertOnConflictSet.insertOnConflict.insert

	onConflictClause, err := insertOnConflictSet.getOnConflictClause()
	if err != nil {
		return 0, err
	}

	insert.addOnConflictClause(onConflictClause)

	return insert.Exec()
}

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

func (insertExec *InsertExec[T]) Exec() (int64, error) {
	return insertExec.insert.Exec()
}

// TODO docs
func (insertExec *InsertExec[T]) ExecInBatches(batchSize int) (int64, error) {
	return insertExec.insert.ExecInBatches(batchSize)
}

// TODO docs
// TODO comentario de que mysql retorna otra cosa
func (insert *Insert[T]) Exec() (int64, error) {
	if insert.err != nil {
		return 0, insert.err
	}

	result := insert.tx.Create(insert.models)

	return result.RowsAffected, result.Error
}

// TODO docs
func (insert *Insert[T]) ExecInBatches(batchSize int) (int64, error) {
	if insert.err != nil {
		return 0, insert.err
	}

	result := insert.tx.CreateInBatches(insert.models, batchSize)

	return result.RowsAffected, result.Error
}
