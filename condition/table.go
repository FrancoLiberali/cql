package condition

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/FrancoLiberali/cql/model"
)

type Table struct {
	Name    string
	Alias   string
	Initial bool
	// Schema is the parsed gorm schema for this table's model, pinned so
	// column-name lookups skip the per-query NamingStrategy allocations.
	// Populated by NewTable / DeliverTable via getTableSchema (cache-shared).
	Schema *schema.Schema
}

// SQLName returns the name that must be used in a sql query to use this table:
// the alias if not empty or the table name
func (t Table) SQLName() string {
	if t.Alias != "" {
		return t.Alias
	}

	return t.Name
}

// Returns true if the Table is the initial table in a query
func (t Table) IsInitial() bool {
	return t.Initial
}

// Returns the related Table corresponding to the model
func (t Table) DeliverTable(query *CQLQuery, model model.Model, relationName string) (Table, error) {
	sch, err := getTableSchema(query.gormDB, model)
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
		Name:    sch.Table,
		Alias:   tableAlias,
		Initial: false,
		Schema:  sch,
	}, nil
}

func NewTable(db *gorm.DB, model model.Model) (Table, error) {
	sch, err := getTableSchema(db, model)
	if err != nil {
		return Table{}, err
	}

	return Table{
		Name:    sch.Table,
		Alias:   sch.Table,
		Initial: true,
		Schema:  sch,
	}, nil
}
