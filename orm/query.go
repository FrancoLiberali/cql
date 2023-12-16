package orm

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type Table struct {
	Name    string
	Alias   string
	Initial bool
}

// Returns true if the Table is the initial table in a query
func (t Table) IsInitial() bool {
	return t.Initial
}

// Returns the related Table corresponding to the model
func (t Table) DeliverTable(query *Query, model Model, relationName string) (Table, error) {
	// get the name of the table for the model
	tableName, err := getTableName(query.gormDB, model)
	if err != nil {
		return Table{}, err
	}

	// add a suffix to avoid tables with the same name when joining
	// the same table more than once
	tableAlias := relationName
	if !t.IsInitial() {
		tableAlias = t.Alias + "__" + relationName
	}

	return Table{
		Name:    tableName,
		Alias:   tableAlias,
		Initial: false,
	}, nil
}

type Query struct {
	gormDB          *gorm.DB
	concernedModels map[reflect.Type][]Table
}

func (query *Query) AddSelect(table Table, fieldID iFieldIdentifier) {
	columnName := fieldID.ColumnName(query, table)

	query.gormDB.Statement.Selects = append(
		query.gormDB.Statement.Selects,
		fmt.Sprintf(
			"%[1]s.%[2]s AS \"%[1]s__%[2]s\"", // name used by gorm to load the fields inside the models
			table.Alias,
			columnName,
		),
	)
}

func (query *Query) Preload(preloadQuery string, args ...interface{}) {
	query.gormDB = query.gormDB.Preload(preloadQuery, args...)
}

func (query *Query) Unscoped() {
	query.gormDB = query.gormDB.Unscoped()
}

func (query *Query) Where(whereQuery interface{}, args ...interface{}) {
	query.gormDB = query.gormDB.Where(whereQuery, args...)
}

func (query *Query) Joins(joinQuery string, args ...interface{}) {
	query.gormDB = query.gormDB.Joins(joinQuery, args...)
}

func (query *Query) Find(dest interface{}, conds ...interface{}) error {
	query.gormDB = query.gormDB.Find(dest, conds...)

	return query.gormDB.Error
}

func (query *Query) AddConcernedModel(model Model, table Table) {
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
	return query.gormDB.NamingStrategy.ColumnName(table.Name, fieldName)
}

func NewQuery[T Model](db *gorm.DB, conditions []Condition[T]) (*Query, error) {
	model := *new(T)

	initialTableName, err := getTableName(db, model)
	if err != nil {
		return nil, err
	}

	initialTable := Table{
		Name:    initialTableName,
		Alias:   initialTableName,
		Initial: true,
	}

	query := &Query{
		gormDB:          db.Select(initialTableName + ".*"),
		concernedModels: map[reflect.Type][]Table{},
	}
	query.AddConcernedModel(model, initialTable)

	for _, condition := range conditions {
		err = condition.ApplyTo(query, initialTable)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}
