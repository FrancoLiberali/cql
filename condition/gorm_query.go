package condition

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

type GormQuery struct {
	GormDB          *gorm.DB
	ConcernedModels map[reflect.Type][]Table
	initialTable    Table
}

// Order specify order when retrieving models from database.
//
// if descending is true, the ordering is in descending direction.
func (query *GormQuery) Order(field IField, descending bool) error {
	table, err := query.GetModelTable(field)
	if err != nil {
		return err
	}

	switch query.Dialector() {
	case sql.Postgres:
		// postgres supports only order by selected fields
		query.AddSelect(table, field)
		query.GormDB = query.GormDB.Order(
			clause.OrderByColumn{
				Column: clause.Column{
					Name: query.getSelectAlias(table, field),
				},
				Desc: descending,
			},
		)

		return nil
	case sql.SQLServer, sql.SQLite, sql.MySQL:
		query.GormDB = query.GormDB.Order(
			clause.OrderByColumn{
				Column: clause.Column{
					Name: field.columnSQL(
						query,
						table,
					),
				},
				Desc: descending,
			},
		)

		return nil
	}

	return nil
}

// Offset specify the number of records to skip before starting to return the records
//
// Offset conditions can be cancelled by using `Offset(-1)`.
func (query *GormQuery) Offset(offset int) {
	query.GormDB = query.GormDB.Offset(offset)
}

// Limit specify the number of records to be retrieved
//
// Limit conditions can be cancelled by using `Limit(-1)`
func (query *GormQuery) Limit(limit int) {
	query.GormDB = query.GormDB.Limit(limit)
}

// Count returns the amount of models that fulfill the conditions
func (query *GormQuery) Count() (int64, error) {
	query.cleanSelects()

	var count int64

	return count, query.GormDB.Count(&count).Error
}

// First finds the first record ordered by primary key, matching given conditions
func (query *GormQuery) First(dest any) error {
	return query.GormDB.First(dest).Error
}

// Take finds the first record returned by the database in no specified order, matching given conditions
func (query *GormQuery) Take(dest any) error {
	return query.GormDB.Take(dest).Error
}

// Last finds the last record ordered by primary key, matching given conditions
func (query *GormQuery) Last(dest any) error {
	return query.GormDB.Last(dest).Error
}

// Find finds all models matching given conditions
func (query *GormQuery) Find(dest any) error {
	return query.GormDB.Find(dest).Error
}

func (query *GormQuery) AddSelect(table Table, fieldID IField) {
	query.GormDB.Statement.Selects = append(
		query.GormDB.Statement.Selects,
		fmt.Sprintf(
			"%s.%s AS %s",
			table.Alias,
			fieldID.columnName(query, table),
			query.getSelectAlias(table, fieldID),
		),
	)
}

func (query *GormQuery) getSelectAlias(table Table, fieldID IField) string {
	return fmt.Sprintf(
		"\"%[1]s__%[2]s\"", // name used by gorm to load the fields inside the models
		table.Alias,
		fieldID.columnName(query, table),
	)
}

func (query *GormQuery) Preload(preloadQuery string, args ...interface{}) {
	query.GormDB = query.GormDB.Preload(preloadQuery, args...)
}

func (query *GormQuery) Unscoped() {
	query.GormDB = query.GormDB.Unscoped()
}

func (query *GormQuery) Where(whereQuery interface{}, args ...interface{}) {
	query.GormDB = query.GormDB.Where(whereQuery, args...)
}

func (query *GormQuery) Joins(joinQuery string, isLeftJoin bool, args ...interface{}) {
	if isLeftJoin {
		query.GormDB = query.GormDB.Joins("LEFT JOIN "+joinQuery, args...)
	} else {
		query.GormDB = query.GormDB.InnerJoins("INNER JOIN "+joinQuery, args...)
	}
}

func (query *GormQuery) AddConcernedModel(model model.Model, table Table) {
	tableList, isPresent := query.ConcernedModels[reflect.TypeOf(model)]
	if !isPresent {
		query.ConcernedModels[reflect.TypeOf(model)] = []Table{table}
	} else {
		tableList = append(tableList, table)
		query.ConcernedModels[reflect.TypeOf(model)] = tableList
	}
}

func (query *GormQuery) GetTables(modelType reflect.Type) []Table {
	tableList, isPresent := query.ConcernedModels[modelType]
	if !isPresent {
		return nil
	}

	return tableList
}

const undefinedAppearance = -1

func (query *GormQuery) GetModelTable(field IField) (Table, error) {
	modelTables := query.GetTables(field.getModelType())
	if modelTables == nil {
		return Table{}, fieldModelNotConcernedError(field)
	}

	if len(modelTables) == 1 {
		return modelTables[0], nil
	}

	appearance := field.getAppearance()

	if appearance == undefinedAppearance {
		return Table{}, appearanceMustBeSelectedError(field)
	}

	if appearance > len(modelTables)-1 {
		return Table{}, appearanceOutOfRangeError(field)
	}

	return modelTables[appearance], nil
}

func (query GormQuery) ColumnName(table Table, fieldName string) string {
	return query.GormDB.NamingStrategy.ColumnName(table.Name, fieldName)
}

func (query GormQuery) Dialector() sql.Dialector {
	return sql.Dialector(query.GormDB.Dialector.Name())
}

func NewGormQuery(db *gorm.DB, initialModel model.Model, initialTable Table) *GormQuery {
	query := &GormQuery{
		GormDB:          db.Model(&initialModel).Select(initialTable.Name + ".*"),
		ConcernedModels: map[reflect.Type][]Table{},
		initialTable:    initialTable,
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

// available for: postgres, sqlite, sqlserver
//
// warning: in sqlite, sqlserver preloads are not allowed
func (query *GormQuery) Returning(dest any) error {
	query.GormDB.Model(dest)

	switch query.Dialector() {
	case sql.Postgres: // support RETURNING from any table
		columns := []clause.Column{}

		for _, selectClause := range query.GormDB.Statement.Selects {
			selectSplit := strings.Split(selectClause, ".")
			columns = append(columns, clause.Column{
				Table: selectSplit[0],
				Name:  selectSplit[1],
				Raw:   true,
			})
		}

		query.GormDB.Clauses(clause.Returning{Columns: columns})
	case sql.SQLite, sql.SQLServer: // supports RETURNING only from main table
		if len(query.GormDB.Statement.Selects) > 1 {
			return preloadsInReturningNotAllowed(query.Dialector())
		}

		query.GormDB.Clauses(clause.Returning{})
	case sql.MySQL: // RETURNING not supported
		return ErrUnsupportedByDatabase
	}

	return nil
}

func (query *GormQuery) cleanSelects() {
	query.GormDB.Statement.Selects = []string{}
}

// Find finds all models matching given conditions
func (query *GormQuery) Update(sets []ISet) (int64, error) {
	tablesAndValues, err := getUpdateTablesAndValues(query, sets)
	if err != nil {
		return 0, err
	}

	updateMap := map[string]any{}

	query.cleanSelects()

	switch query.Dialector() {
	case sql.Postgres, sql.SQLServer, sql.SQLite: // support UPDATE SET FROM
		for field, tableAndValue := range tablesAndValues {
			updateMap[field.columnName(query, tableAndValue.table)] = tableAndValue.value
		}

		query.joinsToFrom()
	case sql.MySQL: // support UPDATE JOIN SET
		// if at least one join is done,
		// allow UPDATE without WHERE as the condition can be the join
		if len(query.GormDB.Statement.Joins) > 0 {
			query.GormDB.AllowGlobalUpdate = true
		}

		sets := clause.Set{}
		updatedTables := []Table{}

		for field, tableAndValue := range tablesAndValues {
			sets = append(sets, clause.Assignment{
				Column: clause.Column{
					Name:  field.columnName(query, tableAndValue.table),
					Table: tableAndValue.table.SQLName(),
				},
				Value: tableAndValue.value,
			})

			updatedTables = append(updatedTables, tableAndValue.table)
		}

		now := time.Now()
		for _, table := range pie.Unique(updatedTables) {
			sets = append(sets, clause.Assignment{
				Column: clause.Column{
					Name:  "updated_at",
					Table: table.SQLName(),
				},
				Value: now,
			})
		}

		query.GormDB.Clauses(sets)
	}

	update := query.GormDB.Updates(updateMap)

	return update.RowsAffected, update.Error
}

func (query *GormQuery) joinsToFrom() {
	joinTables := []clause.Table{}

	for _, join := range query.GormDB.Statement.Joins {
		tableName, tableAlias, onStatement := splitJoin(join.Name)

		joinTables = append(joinTables, clause.Table{
			Name:  tableName,
			Alias: tableAlias,
			Raw:   true, // prevent gorm from putting the alias in quotes
		})

		query.GormDB = query.GormDB.Where(onStatement, join.Conds...)
	}

	if len(joinTables) > 0 {
		query.GormDB.Clauses(
			clause.From{
				Tables: joinTables,
			},
		)
	}

	query.GormDB.Statement.Joins = nil
}

type TableAndValue struct {
	table Table
	value any
}

func getUpdateTablesAndValues(query *GormQuery, sets []ISet) (map[IField]TableAndValue, error) {
	tables := map[IField]TableAndValue{}

	for _, set := range sets {
		field := set.getField()

		table, err := query.GetModelTable(field)
		if err != nil {
			return nil, err
		}

		updateValue, err := getUpdateValue(query, set)
		if err != nil {
			return nil, err
		}

		_, fieldAlreadyPresent := tables[field]
		if fieldAlreadyPresent {
			return nil, fieldIsRepeatedError(field)
		}

		tables[field] = TableAndValue{
			table: table,
			value: updateValue,
		}
	}

	return tables, nil
}

func getUpdateValue(query *GormQuery, set ISet) (any, error) {
	value := set.getValue()

	if iValue, isIValue := value.(IValue); isIValue {
		table, err := query.GetModelTable(iValue.getField())
		if err != nil {
			return nil, err
		}

		valueSQL, valueValues, err := iValue.toSQL(query, table)
		if err != nil {
			return nil, err
		}

		return gorm.Expr(valueSQL, valueValues...), nil
	}

	return value, nil
}

// Splits a JOIN statement into the table name, table alias and ON statement
func splitJoin(joinStatement string) (string, string, string) {
	// remove INNER JOIN and LEFT JOIN
	joinStatement = strings.ReplaceAll(joinStatement, "INNER JOIN ", "")
	joinStatement = strings.ReplaceAll(joinStatement, "LEFT JOIN ", "")

	// divide table and on statement
	joinStatementSplit := strings.Split(joinStatement, " ON ")
	table := joinStatementSplit[0]
	onStatement := joinStatementSplit[1]

	// divide table name and alias
	tableSplit := strings.Split(table, " ")

	return tableSplit[0], tableSplit[1], onStatement
}

func (query *GormQuery) Delete() (int64, error) {
	switch query.Dialector() {
	case sql.Postgres, sql.SQLServer, sql.SQLite: // support UPDATE SET FROM
		query.joinsToFrom()
	case sql.MySQL:
		// if at least one join is done,
		// allow UPDATE without WHERE as the condition can be the join
		if len(query.GormDB.Statement.Joins) > 0 {
			query.GormDB.AllowGlobalUpdate = true
		}

		query.GormDB.Clauses(clause.Set{clause.Assignment{
			Column: clause.Column{
				Name:  "deleted_at",
				Table: query.initialTable.SQLName(),
			},
			Value: time.Now(),
		}})
	}

	query.GormDB.Statement.Selects = []string{}

	deleteTx := query.GormDB.Delete(query.GormDB.Statement.Model)

	return deleteTx.RowsAffected, deleteTx.Error
}
