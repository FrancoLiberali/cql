package sql

type Operator uint

const (
	Eq Operator = iota
	NotEq
	Lt
	LtOrEq
	Gt
	GtOrEq
	Between
	NotBetween
	IsDistinct
	IsNotDistinct
	Like
	Escape
	ArrayIn
	ArrayNotIn
	And
	Or
	Not
)

func (op Operator) String() string {
	return operatorToSQL[op]
}

var operatorToSQL = map[Operator]string{
	Eq:            "=",
	NotEq:         "<>",
	Lt:            "<",
	LtOrEq:        "<=",
	Gt:            ">",
	GtOrEq:        ">=",
	Between:       "BETWEEN",
	NotBetween:    "NOT BETWEEN",
	IsDistinct:    "IS DISTINCT FROM",
	IsNotDistinct: "IS NOT DISTINCT FROM",
	Like:          "LIKE",
	Escape:        "ESCAPE",
	ArrayIn:       "IN",
	ArrayNotIn:    "NOT IN",
	And:           "AND",
	Or:            "OR",
	Not:           "NOT",
}
