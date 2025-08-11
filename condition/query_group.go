package condition

type QueryGroup struct {
	gormQuery *GormQuery
	err       error
	fields    []IField
}

// TODO docs
func (query *QueryGroup) Having(condition AggregationCondition) *QueryGroup {
	sql, args, err := condition.toSQL(query.gormQuery)
	if err != nil {
		query.addError(methodError(err, "Having"))
		return query
	}

	query.gormQuery.Having(sql, args...)

	return query
}

func (query *QueryGroup) Select(aggregation IAggregation, as string) *QueryGroup {
	selectSQL, err := aggregation.toSelectSQL(query.gormQuery, as)
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
