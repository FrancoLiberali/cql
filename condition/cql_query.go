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

type CQLQuery struct {
	gormDB          *gorm.DB
	concernedModels map[reflect.Type][]Table
	initialTable    Table
	selectClause    clause.Expr
}

// Order specify order when retrieving models from database.
//
// if descending is true, the ordering is in descending direction.
func (query *CQLQuery) Order(field IField, descending bool) error {
	table, err := query.GetModelTable(field)
	if err != nil {
		return err
	}

	switch query.Dialector() {
	case sql.Postgres:
		// postgres supports only order by selected fields
		query.AddSelectField(table, field, true)
		query.gormDB = query.gormDB.Order(
			clause.OrderByColumn{
				Column: clause.Column{
					Name: query.getSelectAlias(table, field),
				},
				Desc: descending,
			},
		)

		return nil
	case sql.SQLServer, sql.SQLite, sql.MySQL:
		query.gormDB = query.gormDB.Order(
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
func (query *CQLQuery) Offset(offset int) {
	query.gormDB = query.gormDB.Offset(offset)
}

// Limit specify the number of records to be retrieved
//
// Limit conditions can be cancelled by using `Limit(-1)`
func (query *CQLQuery) Limit(limit int) {
	query.gormDB = query.gormDB.Limit(limit)
}

// GroupBy arrange identical data into groups
func (query *CQLQuery) GroupBy(fields []IField) error {
	query.cleanSelects()

	for _, field := range fields {
		fieldSQL, _, err := field.ToSQL(query)
		if err != nil {
			return err
		}

		query.AddSelectForAggregation(fieldSQL, nil)

		query.gormDB.Group(fieldSQL)
	}

	return nil
}

// Having allows filter groups of rows based on conditions involving aggregate functions
func (query *CQLQuery) Having(sql string, args ...any) {
	query.gormDB.Having(sql, args...)
}

// Count returns the amount of models that fulfill the conditions
func (query *CQLQuery) Count() (int64, error) {
	query.cleanSelects()

	var count int64

	return count, query.gormDB.Count(&count).Error
}

// First finds the first record ordered by primary key, matching given conditions
func (query *CQLQuery) First(dest any) error {
	return query.gormDB.First(dest).Error
}

// Take finds the first record returned by the database in no specified order, matching given conditions
func (query *CQLQuery) Take(dest any) error {
	return query.gormDB.Take(dest).Error
}

// Last finds the last record ordered by primary key, matching given conditions
func (query *CQLQuery) Last(dest any) error {
	return query.gormDB.Last(dest).Error
}

// Find finds all models matching given conditions
func (query *CQLQuery) Find(dest any) error {
	return query.gormDB.Find(dest).Error
}

// Select specify fields that you want when doing group bys
func (query *CQLQuery) AddSelectForAggregation(sql string, values []any) {
	newSQL := query.selectClause.SQL
	if newSQL != "" {
		newSQL = query.selectClause.SQL + "," + sql
	} else {
		newSQL = sql
	}

	query.selectClause = clause.Expr{
		SQL:  newSQL,
		Vars: append(query.selectClause.Vars, values...),
	}

	query.gormDB.Statement.AddClause(clause.Select{
		Expression: query.selectClause,
	})
}

// Select specify fields that you want when querying, creating, updating
func (query *CQLQuery) AddSelect(sql string) {
	query.gormDB.Statement.Selects = append(
		query.gormDB.Statement.Selects,
		sql,
	)
}

func (query *CQLQuery) AddSelectField(table Table, fieldID IField, addAs bool) {
	columnName := fieldID.columnSQL(query, table)

	if addAs {
		columnName += " AS " + query.getSelectAlias(table, fieldID)
	}

	query.AddSelect(columnName)
}

func (query *CQLQuery) getSelectAlias(table Table, fieldID IField) string {
	return fmt.Sprintf(
		"\"%[1]s__%[2]s\"", // name used by gorm to load the fields inside the models
		table.Alias,
		fieldID.columnName(query, table),
	)
}

func (query *CQLQuery) Preload(preloadQuery string, args ...interface{}) {
	query.gormDB = query.gormDB.Preload(preloadQuery, args...)
}

func (query *CQLQuery) Unscoped() {
	query.gormDB = query.gormDB.Unscoped()
}

func (query *CQLQuery) Where(whereQuery interface{}, args ...interface{}) {
	query.gormDB = query.gormDB.Where(whereQuery, args...)
}

func (query *CQLQuery) Joins(joinQuery string, isLeftJoin bool, args ...interface{}) {
	if isLeftJoin {
		query.gormDB = query.gormDB.Joins("LEFT JOIN "+joinQuery, args...)
	} else {
		query.gormDB = query.gormDB.InnerJoins("INNER JOIN "+joinQuery, args...)
	}
}

func (query *CQLQuery) AddConcernedModel(model model.Model, table Table) {
	tableList, isPresent := query.concernedModels[reflect.TypeOf(model)]
	if !isPresent {
		query.concernedModels[reflect.TypeOf(model)] = []Table{table}
	} else {
		tableList = append(tableList, table)
		query.concernedModels[reflect.TypeOf(model)] = tableList
	}
}

func (query *CQLQuery) GetTables(modelType reflect.Type) []Table {
	tableList, isPresent := query.concernedModels[modelType]
	if !isPresent {
		return nil
	}

	return tableList
}

func (query *CQLQuery) GetModelTable(field IField) (Table, error) {
	modelTables := query.GetTables(field.getModelType())
	if modelTables == nil {
		return Table{}, fieldModelNotConcernedError(field)
	}

	if len(modelTables) == 1 {
		return modelTables[0], nil
	}

	appearance, isPresent := field.getAppearance()

	if !isPresent {
		return Table{}, appearanceMustBeSelectedError(field)
	}

	if appearance > uint(len(modelTables))-1 {
		return Table{}, appearanceOutOfRangeError(field)
	}

	return modelTables[appearance], nil
}

func (query CQLQuery) ColumnName(table Table, fieldName string) string {
	return query.gormDB.NamingStrategy.ColumnName(table.Name, fieldName)
}

func (query CQLQuery) Dialector() sql.Dialector {
	return sql.Dialector(query.gormDB.Dialector.Name())
}

func NewGormQuery(db *gorm.DB, initialModel model.Model, initialTable Table) *CQLQuery {
	query := &CQLQuery{
		gormDB:          db.Model(&initialModel).Select(initialTable.Name + ".*"),
		concernedModels: map[reflect.Type][]Table{},
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
func (query *CQLQuery) Returning(dest any) error {
	query.gormDB.Model(dest)

	switch query.Dialector() {
	case sql.Postgres: // support RETURNING from any table
		columns := []clause.Column{}

		for _, selectClause := range query.gormDB.Statement.Selects {
			selectSplit := strings.Split(selectClause, ".")
			columns = append(columns, clause.Column{
				Table: selectSplit[0],
				Name:  selectSplit[1],
				Raw:   true,
			})
		}

		query.gormDB.Clauses(clause.Returning{Columns: columns})
	case sql.SQLite, sql.SQLServer: // supports RETURNING only from main table
		if len(query.gormDB.Statement.Selects) > 1 {
			return preloadsInReturningNotAllowed(query.Dialector())
		}

		query.gormDB.Clauses(clause.Returning{})
	case sql.MySQL: // RETURNING not supported
		return ErrUnsupportedByDatabase
	}

	return nil
}

func (query *CQLQuery) cleanSelects() {
	query.gormDB.Statement.Selects = []string{}
}

// Find finds all models matching given conditions
func (query *CQLQuery) Update(sets []ISet) (int64, error) {
	updateMap := map[string]any{}

	query.cleanSelects()

	switch query.Dialector() {
	case sql.Postgres, sql.SQLServer, sql.SQLite: // support UPDATE SET FROM
		for _, set := range sets {
			field := set.getField()

			updateValue, err := getUpdateValue(query, set)
			if err != nil {
				return 0, err
			}

			table, err := query.GetModelTable(field)
			if err != nil {
				return 0, err
			}

			updateMap[field.columnName(query, table)] = updateValue
		}

		query.joinsToFrom()
	case sql.MySQL: // support UPDATE JOIN SET
		// if at least one join is done,
		// allow UPDATE without WHERE as the condition can be the join
		if len(query.gormDB.Statement.Joins) > 0 {
			query.gormDB.AllowGlobalUpdate = true
		}

		setClause := clause.Set{}
		updatedTables := []Table{}

		for _, set := range sets {
			field := set.getField()

			updateValue, err := getUpdateValue(query, set)
			if err != nil {
				return 0, err
			}

			table, err := query.GetModelTable(field)
			if err != nil {
				return 0, err
			}

			setClause = append(setClause, clause.Assignment{
				Column: clause.Column{
					Name:  field.columnName(query, table),
					Table: table.SQLName(),
				},
				Value: updateValue,
			})

			updatedTables = append(updatedTables, table)
		}

		now := time.Now()
		for _, table := range pie.Unique(updatedTables) {
			setClause = append(setClause, clause.Assignment{
				Column: clause.Column{
					Name:  "updated_at",
					Table: table.SQLName(),
				},
				Value: now,
			})
		}

		query.gormDB.Clauses(setClause)
	}

	update := query.gormDB.Updates(updateMap)

	return update.RowsAffected, update.Error
}

func (query *CQLQuery) joinsToFrom() {
	joinTables := []clause.Table{}

	for _, join := range query.gormDB.Statement.Joins {
		tableName, tableAlias, onStatement := splitJoin(join.Name)

		joinTables = append(joinTables, clause.Table{
			Name:  tableName,
			Alias: tableAlias,
			Raw:   true, // prevent gorm from putting the alias in quotes
		})

		query.gormDB = query.gormDB.Where(onStatement, join.Conds...)
	}

	if len(joinTables) > 0 {
		query.gormDB.Clauses(
			clause.From{
				Tables: joinTables,
			},
		)
	}

	query.gormDB.Statement.Joins = nil
}

func getUpdateValue(query *CQLQuery, set ISet) (any, error) {
	if value := set.getValue(); value != nil {
		valueSQL, valueValues, err := set.getValue().ToSQL(query)
		if err != nil {
			return nil, err
		}

		if valueSQL != "" {
			return gorm.Expr(valueSQL, valueValues...), nil
		}

		return valueValues[0], nil
	}

	return nil, nil //nolint:nilnil // is necessary
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

func (query *CQLQuery) Delete() (int64, error) {
	switch query.Dialector() {
	case sql.Postgres, sql.SQLServer, sql.SQLite: // support UPDATE SET FROM
		query.joinsToFrom()
	case sql.MySQL:
		// if at least one join is done,
		// allow UPDATE without WHERE as the condition can be the join
		if len(query.gormDB.Statement.Joins) > 0 {
			query.gormDB.AllowGlobalUpdate = true
		}

		query.gormDB.Clauses(clause.Set{clause.Assignment{
			Column: clause.Column{
				Name:  "deleted_at",
				Table: query.initialTable.SQLName(),
			},
			Value: time.Now(),
		}})
	}

	query.gormDB.Statement.Selects = []string{}

	deleteTx := query.gormDB.Delete(query.gormDB.Statement.Model)

	return deleteTx.RowsAffected, deleteTx.Error
}
