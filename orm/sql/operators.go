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

func (op Operator) Name() string {
	return operatorToName[op]
}

var operatorToName = map[Operator]string{
	Eq:            "Eq",
	NotEq:         "NotEq",
	Lt:            "Lt",
	LtOrEq:        "LtOrEq",
	Gt:            "Gt",
	GtOrEq:        "GtOrEq",
	Between:       "Between",
	NotBetween:    "NotBetween",
	IsDistinct:    "IsDistinct",
	IsNotDistinct: "IsNotDistinct",
	Like:          "Like",
	Escape:        "Escape",
	ArrayIn:       "ArrayIn",
	ArrayNotIn:    "ArrayNotIn",
	And:           "And",
	Or:            "Or",
	Not:           "Not",
}
