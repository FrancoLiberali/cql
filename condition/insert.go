package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Insert[T model.Model] struct {
	tx                *gorm.DB
	err               error
	models            []*T
	onConflictColumns []string
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
	onConflictColumns := make([]clause.Column, 0, len(fields))

	for _, field := range fields {
		onConflictColumns = append(
			onConflictColumns,
			// TODO deberia ser el nombre completo por si los nombres se repiten
			clause.Column{Name: field.fieldName()},
		)
	}

	return &InsertOnConflict[T]{
		insert:            insert,
		onConflictColumns: onConflictColumns,
	}
}

type InsertOnConflict[T model.Model] struct {
	insert *Insert[T]

	onConflictColumns []clause.Column
}

func (insertOnConflict *InsertOnConflict[T]) DoNothing() *Insert[T] {
	// TODO esto podria ir en cql query
	insertOnConflict.insert.tx = insertOnConflict.insert.tx.Clauses(clause.OnConflict{
		// TODO ver todas las opciones que tiene esto
		Columns:   insertOnConflict.onConflictColumns,
		DoNothing: true,
	})

	return insertOnConflict.insert
}

func (insertOnConflict *InsertOnConflict[T]) UpdateAll() *Insert[T] {
	insertOnConflict.insert.tx = insertOnConflict.insert.tx.Clauses(clause.OnConflict{
		Columns:   insertOnConflict.onConflictColumns,
		UpdateAll: true,
	})

	return insertOnConflict.insert
}

func (insertOnConflict *InsertOnConflict[T]) Update(fields ...IField) *Insert[T] {
	fieldNames := make([]string, 0, len(fields))

	for _, field := range fields {
		fieldNames = append(
			fieldNames,
			// TODO deberia ser el nombre completo por si los nombres se repiten
			field.fieldName(),
		)
	}

	insertOnConflict.insert.tx = insertOnConflict.insert.tx.Clauses(clause.OnConflict{
		Columns:   insertOnConflict.onConflictColumns,
		DoUpdates: clause.AssignmentColumns(fieldNames),
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

		if setSQL == "" && len(setValues) == 1 {
			assignments[set.getField().fieldName()] = setValues[0]
		} else {
			// TODO esto no anda creo
			assignments[set.getField().fieldName()] = gorm.Expr(
				setSQL,
				setValues...,
			)
		}
	}

	return clause.OnConflict{
		Columns:   insertOnConflictSet.insertOnConflict.onConflictColumns,
		DoUpdates: clause.Assignments(assignments),
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
