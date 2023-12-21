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
	// mysql
	MySQLXor
	MySQLRegexp
	MySQLNullSafeEqual
	// postgresql
	PostgreSQLILike
	PostgreSQLSimilarTo
	PostgreSQLPosixMatch
	PostgreSQLPosixIMatch
	// sqlite
	SQLiteGlob
)

func (op Operator) String() string {
	return operatorToSQL[op]
}

var operatorToSQL = map[Operator]string{
	Eq:                    "=",
	NotEq:                 "<>",
	Lt:                    "<",
	LtOrEq:                "<=",
	Gt:                    ">",
	GtOrEq:                ">=",
	Between:               "BETWEEN",
	NotBetween:            "NOT BETWEEN",
	IsDistinct:            "IS DISTINCT FROM",
	IsNotDistinct:         "IS NOT DISTINCT FROM",
	Like:                  "LIKE",
	Escape:                "ESCAPE",
	ArrayIn:               "IN",
	ArrayNotIn:            "NOT IN",
	And:                   "AND",
	Or:                    "OR",
	Not:                   "NOT",
	MySQLXor:              "XOR",
	MySQLRegexp:           "REGEXP",
	MySQLNullSafeEqual:    "<=>",
	PostgreSQLILike:       "ILIKE",
	PostgreSQLSimilarTo:   "SIMILAR TO",
	PostgreSQLPosixMatch:  "~",
	PostgreSQLPosixIMatch: "~*",
	SQLiteGlob:            "GLOB",
}

func (op Operator) Name() string {
	return operatorToName[op]
}

var operatorToName = map[Operator]string{
	Eq:                    "Eq",
	NotEq:                 "NotEq",
	Lt:                    "Lt",
	LtOrEq:                "LtOrEq",
	Gt:                    "Gt",
	GtOrEq:                "GtOrEq",
	Between:               "Between",
	NotBetween:            "NotBetween",
	IsDistinct:            "IsDistinct",
	IsNotDistinct:         "IsNotDistinct",
	Like:                  "Like",
	Escape:                "Escape",
	ArrayIn:               "ArrayIn",
	ArrayNotIn:            "ArrayNotIn",
	And:                   "And",
	Or:                    "Or",
	Not:                   "Not",
	MySQLXor:              "mysql.Xor",
	MySQLRegexp:           "mysql.Regexp",
	MySQLNullSafeEqual:    "IsNotDistinct",
	PostgreSQLILike:       "psql.ILike",
	PostgreSQLSimilarTo:   "psql.SimilarTo",
	PostgreSQLPosixMatch:  "psql.POSIXMatch",
	PostgreSQLPosixIMatch: "psql.POSIXIMatch",
	SQLiteGlob:            "sqlite.Glob",
}

func (op Operator) Supports(dialector Dialector) bool {
	supportedDialector, present := operatorDialector[op]
	if !present {
		// supports all dialectors
		return true
	}

	return supportedDialector == dialector
}

var operatorDialector = map[Operator]Dialector{ //nolint:exhaustive // missing key is supported
	MySQLXor:              MySQL,
	MySQLRegexp:           MySQL,
	MySQLNullSafeEqual:    MySQL,
	PostgreSQLILike:       Postgres,
	PostgreSQLSimilarTo:   Postgres,
	PostgreSQLPosixMatch:  Postgres,
	PostgreSQLPosixIMatch: Postgres,
	SQLiteGlob:            SQLite,
}
