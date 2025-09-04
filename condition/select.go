package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

type Selection[T any] interface {
	Apply(value any, result *T) error
	ValueType() any
}

func Select[TResults any, TModel model.Model](
	query *Query[TModel],
	selections []Selection[TResults],
) ([]TResults, error) {
	if query.err != nil {
		return nil, query.err
	}

	// TODO aca poner las selecciones
	rows, err := query.cqlQuery.gormDB.Select("?, ?", 42, 43).Rows()
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
