package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

type Selection[T any] interface {
	ResultsType() T
	Apply(value any, result *T) error
}

func Select[TResults any, TModel model.Model](
	query *Query[TModel],
	selections []Selection[TResults],
) ([]TResults, error) {
	if query.err != nil {
		return nil, query.err
	}

	// TODO aca poner las selecciones
	rows, err := query.cqlQuery.gormDB.Select("?", 42).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TResults

	// cols := make([]any, len(selections))

	for rows.Next() {
		// err = rows.Scan(cols...)
		var algo float64

		err = rows.Scan(&algo)
		if err != nil {
			return nil, err
		}

		var result TResults

		for _, selection := range selections {
			// err = selection.Apply(cols[i], &result)
			err = selection.Apply(algo, &result)
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
