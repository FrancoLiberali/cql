package condition

type QueryGroup struct {
	cqlQuery *CQLQuery
	err      error
	fields   []IField
}

func (query *QueryGroup) getError() error {
	return query.err
}

func (query *QueryGroup) getCQLQuery() *CQLQuery {
	return query.cqlQuery
}

// Having allows filter groups of rows based on conditions involving aggregate functions
func (query *QueryGroup) Having(conditions ...AggregationCondition) *QueryGroup {
	for _, condition := range conditions {
		sql, args, err := condition.toSQL(query.cqlQuery)
		if err != nil {
			query.addError(methodError(err, "Having"))
			return query
		}

		query.cqlQuery.Having(sql, args...)
	}

	return query
}

func (query *QueryGroup) addError(err error) {
	if err != nil && query.err == nil {
		query.err = err
	}
}
