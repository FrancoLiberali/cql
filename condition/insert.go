package condition

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/FrancoLiberali/cql/model"
)

type Insert[T model.Model] struct {
	tx     *gorm.DB
	err    error
	models []*T
}

func NewInsert[T model.Model](tx *gorm.DB, models []*T) *Insert[T] {
	return &Insert[T]{
		tx:     tx,
		models: models,
	}
}

// TODO este ifield puede ser de cualquier model
// o necesita linter
// o hacer dos metodos uno para insertar solo un model y otro para insertar junto con las relations
// que ahi si necesitas Ifield. El linter lo necesitas de una forma o la otra
func (insert *Insert[T]) OnConflict(fields ...IField) *InsertOnConflict[T] {
	// TODO postgresql necesita al menos 1 mientras que sqllite no
	// menos para do nothing
	onConflictColumns := pie.Map(
		insert.getFieldNames(fields),
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
	query := NewQuery[T](insert.tx)

	fieldNames := make([]string, 0, len(fields))

	for _, field := range fields {
		fieldNames = append(
			fieldNames,
			field.columnName(query.cqlQuery, query.cqlQuery.initialTable),
		)
	}

	return fieldNames
}

func (insert *Insert[T]) OnConstraint(constraintName string) *InsertOnConflict[T] {
	return &InsertOnConflict[T]{
		insert:       insert,
		onConstraint: constraintName,
	}
}

type InsertOnConflict[T model.Model] struct {
	insert *Insert[T]

	onConflictColumns []clause.Column
	onConstraint      string
}

func (insertOnConflict *InsertOnConflict[T]) DoNothing() *Insert[T] {
	// TODO esto podria ir en cql query
	insertOnConflict.insert.tx = insertOnConflict.insert.tx.Clauses(clause.OnConflict{
		OnConstraint: insertOnConflict.onConstraint,
		Columns:      insertOnConflict.onConflictColumns,
		DoNothing:    true,
	})

	return insertOnConflict.insert
}

func (insertOnConflict *InsertOnConflict[T]) UpdateAll() *Insert[T] {
	insertOnConflict.insert.tx = insertOnConflict.insert.tx.Clauses(clause.OnConflict{
		OnConstraint: insertOnConflict.onConstraint,
		Columns:      insertOnConflict.onConflictColumns,
		UpdateAll:    true,
	})

	return insertOnConflict.insert
}

func (insertOnConflict *InsertOnConflict[T]) Update(fields ...IField) *Insert[T] {
	// TODO estos fields podrian ser de solo T
	fieldNames := insertOnConflict.insert.getFieldNames(fields)

	insertOnConflict.insert.tx = insertOnConflict.insert.tx.Clauses(clause.OnConflict{
		OnConstraint: insertOnConflict.onConstraint,
		Columns:      insertOnConflict.onConflictColumns,
		DoUpdates:    clause.AssignmentColumns(fieldNames),
	})

	return insertOnConflict.insert
}

func (insertOnConflict *InsertOnConflict[T]) Set(sets ...*Set[T]) *InsertOnConflictSet[T] {
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

	insert.tx = insert.tx.Clauses(onConflictClause)

	return insert.Exec()
}

func (insertOnConflictSet *InsertOnConflictSet[T]) Where(conditions ...Condition[T]) *Insert[T] {
	insert := insertOnConflictSet.insertOnConflict.insert

	query := NewQuery(insert.tx, conditions...)

	onConflictClause, err := insertOnConflictSet.getOnConflictClause()
	if err != nil {
		insert.err = err
	}

	where, isWhere := query.cqlQuery.gormDB.Statement.Clauses["WHERE"].Expression.(clause.Where)
	if isWhere {
		onConflictClause.Where = where
	}

	insert.tx = insert.tx.Clauses(onConflictClause)

	return insert
}

func (insertOnConflictSet *InsertOnConflictSet[T]) getOnConflictClause() (clause.OnConflict, error) {
	assignments := map[string]any{}

	insert := insertOnConflictSet.insertOnConflict.insert

	query := NewQuery[T](insert.tx)

	for _, set := range insertOnConflictSet.sets {
		setSQL, setValues, err := set.getValue().ToSQL(query.cqlQuery)
		if err != nil {
			return clause.OnConflict{}, err
		}

		fieldColumn := set.getField().columnName(query.cqlQuery, query.cqlQuery.initialTable)

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

// TODO docs
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
