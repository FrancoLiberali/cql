package condition

import (
	"strings"

	"github.com/FrancoLiberali/cql/model"
)

type Selection[T any] interface {
	Apply(value any, result *T) error
	ValueType() any
	ToSQL(query *CQLQuery) (string, []any, error)
}

func Select[TResults any, TModel model.Model](
	query *Query[TModel],
	selections []Selection[TResults],
) ([]TResults, error) {
	if query.err != nil {
		return nil, query.err
	}

	selectSQLs := make([]string, 0, len(selections))

	var allValues []any

	for _, selection := range selections {
		sql, values, err := selection.ToSQL(query.cqlQuery)
		if err != nil {
			return nil, err
		}

		selectSQLs = append(selectSQLs, sql)
		allValues = append(allValues, values...)
	}

	rows, err := query.cqlQuery.gormDB.Select(
		strings.Join(selectSQLs, ", "),
		allValues...,
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TResults

	cols := make([]any, 0, len(selections))

	for _, selection := range selections {
		cols = append(cols, selection.ValueType())
	}

	for rows.Next() {
		err = rows.Scan(cols...)
		if err != nil {
			return nil, err
		}

		var result TResults

		for i, selection := range selections {
			err = selection.Apply(cols[i], &result)
			if err != nil {
				return nil, err
			}
		}

		results = append(results, result)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return results, nil
}
