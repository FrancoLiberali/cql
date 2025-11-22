package condition

import (
	"strings"
)

type Selection[T any] interface {
	Apply(value any, result *T) error
	ValueType() any
	ToSQL(query *CQLQuery) (string, []any, error)
}

type IQuery interface {
	getError() error
	getCQLQuery() *CQLQuery
}

func Select[TResults any](
	query IQuery,
	selections []Selection[TResults],
) ([]TResults, error) {
	if query.getError() != nil {
		return nil, query.getError()
	}

	selectSQLs := make([]string, 0, len(selections))

	var allValues []any

	cqlQuery := query.getCQLQuery()

	for _, selection := range selections {
		sql, values, err := selection.ToSQL(cqlQuery)
		if err != nil {
			return nil, err
		}

		selectSQLs = append(selectSQLs, sql)
		allValues = append(allValues, values...)
	}

	var extraCols []any

	// add selects that where already in the query, for example for the order
	if cqlQuery.selectClause.SQL != "" {
		selectSQLs = append(selectSQLs, cqlQuery.selectClause.SQL)
		allValues = append(allValues, cqlQuery.selectClause.Vars...)

		for range len(strings.Split(cqlQuery.selectClause.SQL, ",")) {
			var anything any
			extraCols = append(extraCols, &anything)
		}
	}

	rows, err := cqlQuery.gormDB.Select(
		strings.Join(selectSQLs, ", "),
		allValues...,
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TResults

	cols := make([]any, 0, len(selections)+len(extraCols))

	for _, selection := range selections {
		cols = append(cols, selection.ValueType())
	}

	cols = append(cols, extraCols...)

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
