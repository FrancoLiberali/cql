package query

import (
	"fmt"
	"reflect"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ditrit/badaas/orm/model"
)

type Query struct {
	GormDB          *gorm.DB
	concernedModels map[reflect.Type][]Table
}

func (query *Query) AddSelect(table Table, fieldID IFieldIdentifier) {
	columnName := fieldID.ColumnName(query, table)

	query.GormDB.Statement.Selects = append(
		query.GormDB.Statement.Selects,
		fmt.Sprintf(
			"%[1]s.%[2]s AS \"%[1]s__%[2]s\"", // name used by gorm to load the fields inside the models
			table.Alias,
			columnName,
		),
	)
}

func (query *Query) Preload(preloadQuery string, args ...interface{}) {
	query.GormDB = query.GormDB.Preload(preloadQuery, args...)
}

func (query *Query) Unscoped() {
	query.GormDB = query.GormDB.Unscoped()
}

func (query *Query) Where(whereQuery interface{}, args ...interface{}) {
	query.GormDB = query.GormDB.Where(whereQuery, args...)
}

func (query *Query) Joins(joinQuery string, args ...interface{}) {
	query.GormDB = query.GormDB.Joins(joinQuery, args...)
}

func (query *Query) Find(dest interface{}, conds ...interface{}) error {
	query.GormDB = query.GormDB.Find(dest, conds...)

	return query.GormDB.Error
}

func (query *Query) AddConcernedModel(model model.Model, table Table) {
	tableList, isPresent := query.concernedModels[reflect.TypeOf(model)]
	if !isPresent {
		query.concernedModels[reflect.TypeOf(model)] = []Table{table}
	} else {
		tableList = append(tableList, table)
		query.concernedModels[reflect.TypeOf(model)] = tableList
	}
}

func (query *Query) GetTables(modelType reflect.Type) []Table {
	tableList, isPresent := query.concernedModels[modelType]
	if !isPresent {
		return nil
	}

	return tableList
}

func (query Query) ColumnName(table Table, fieldName string) string {
	return query.GormDB.NamingStrategy.ColumnName(table.Name, fieldName)
}

func NewQuery(db *gorm.DB, initialModel model.Model, initialTable Table) *Query {
	query := &Query{
		GormDB:          db.Select(initialTable.Name + ".*"),
		concernedModels: map[reflect.Type][]Table{},
	}

	query.AddConcernedModel(initialModel, initialTable)

	return query
}

// Get the name of the table in "db" in which the data for "entity" is saved
// returns error is table name can not be found by gorm,
// probably because the type of "entity" is not registered using AddModel
func getTableName(db *gorm.DB, entity any) (string, error) {
	schemaName, err := schema.Parse(entity, &sync.Map{}, db.NamingStrategy)
	if err != nil {
		return "", err
	}

	return schemaName.Table, nil
}
