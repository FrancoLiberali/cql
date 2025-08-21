package condition

type QueryGroup struct {
	gormQuery *CQLQuery
	err       error
	fields    []IField
}

// Having allows filter groups of rows based on conditions involving aggregate functions
func (query *QueryGroup) Having(conditions ...AggregationCondition) *QueryGroup {
	for _, condition := range conditions {
		sql, args, err := condition.toSQL(query.gormQuery)
		if err != nil {
			query.addError(methodError(err, "Having"))
			return query
		}

		query.gormQuery.Having(sql, args...)
	}

	return query
}

func (query *QueryGroup) Select(aggregation Aggregation, as string) *QueryGroup {
	selectSQL, values, err := aggregation.toSelectSQL(query.gormQuery, as)
	if err != nil {
		query.addError(methodError(err, "Select"))
		return query
	}

	query.gormQuery.AddSelectForAggregation(selectSQL, values)

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
