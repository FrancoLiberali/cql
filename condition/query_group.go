package condition

type QueryGroup struct {
	gormQuery *GormQuery
	err       error
	fields    []IField
}

func (query *QueryGroup) Having(condition AggregationCondition) *QueryGroup {
	var table Table

	if condition.aggregation.field != nil { // CountAll
		var err error

		table, err = query.gormQuery.GetModelTable(condition.aggregation.field)
		if err != nil {
			query.addError(methodError(err, "Having"))
			return query
		}
	}

	sql, args, err := condition.toSQL(query.gormQuery, table)
	if err != nil {
		query.addError(methodError(err, "Having"))
		return query
	}

	query.gormQuery.Having(sql, args...)

	return query
}

func (query *QueryGroup) Select(aggregation Aggregation, as string) *QueryGroup {
	var table Table

	if aggregation.field != nil { // CountAll
		var err error

		table, err = query.gormQuery.GetModelTable(aggregation.field)
		if err != nil {
			query.addError(methodError(err, "Select"))
			return query
		}
	}

	selectSQL, err := aggregation.toSelectSQL(query.gormQuery, table, as)
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
