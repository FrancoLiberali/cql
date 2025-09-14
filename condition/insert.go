package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Insert[T model.Model] struct {
	tx                *gorm.DB
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

// TODO docs
func (insert *Insert[T]) Exec() (int64, error) {
	result := insert.tx.Create(insert.models)

	return result.RowsAffected, result.Error
}

// TODO docs
func (insert *Insert[T]) ExecInBatches(batchSize int) (int64, error) {
	result := insert.tx.CreateInBatches(insert.models, batchSize)

	return result.RowsAffected, result.Error
}
