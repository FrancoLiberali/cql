package condition

type QueryGroup struct {
	gormQuery *GormQuery
	err       error
	fields    []IField
}

func (query *QueryGroup) Select(aggregation Aggregation, as string) *QueryGroup {
	var table Table

	if aggregation.field != nil { // CountAll
		var err error

		table, err = query.gormQuery.GetModelTable(aggregation.field, UndefinedJoinNumber)
		if err != nil {
			query.addError(methodError(err, "Select"))
			return query
		}
	}

	selectSQL, err := aggregation.toSQL(query.gormQuery, table, as)
	if err != nil {
		query.addError(methodError(err, "Select"))
		return query
	}

	query.gormQuery.AddSelect(selectSQL)

	return query
}

func (query *QueryGroup) Into(result any) error {
	if query.err != nil {
		return query.err
	}

	return query.gormQuery.Find(result)
}

func (query *QueryGroup) addError(err error) {
	if err != nil && query.err == nil {
		query.err = err
	}
}
