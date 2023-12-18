package query

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/model"
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
func (t Table) DeliverTable(query *Query, model model.Model, relationName string) (Table, error) {
	// get the name of the table for the model
	tableName, err := getTableName(query.GormDB, model)
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

func NewTable(db *gorm.DB, model model.Model) (Table, error) {
	initialTableName, err := getTableName(db, model)
	if err != nil {
		return Table{}, err
	}

	return Table{
		Name:    initialTableName,
		Alias:   initialTableName,
		Initial: true,
	}, nil
}
