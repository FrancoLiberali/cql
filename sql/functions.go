package sql

import "strings"

type Function interface {
	ApplyTo(internalSQL string, values int) string
}

type FunctionFunction struct {
	sqlPrefix   string
	sqlFunction string
	sqlSuffix   string
}

func (f FunctionFunction) ApplyTo(internalSQL string, values int) string {
	finalSQL := f.sqlPrefix

	if f.sqlFunction != "" {
		placeholders := strings.Repeat(", ?", values)

		finalSQL += f.sqlFunction + "(" + internalSQL + placeholders + ")"
	}

	if f.sqlSuffix != "" {
		finalSQL += " " + f.sqlSuffix
	}

	return finalSQL
}

type OperatorFunction struct {
	sqlOperator string
}

func (f OperatorFunction) ApplyTo(internalSQL string, _ int) string {
	return "(" + internalSQL + " " + f.sqlOperator + " ?)"
}

type PreOperatorFunction struct {
	sqlOperator string
}

func (f PreOperatorFunction) ApplyTo(internalSQL string, _ int) string {
	return f.sqlOperator + internalSQL
}

type FunctionByDialector struct {
	functions map[Dialector]Function
	Name      string
}

const all = "all"

var (
	// Numeric

	Plus    = FunctionByDialector{functions: map[Dialector]Function{all: OperatorFunction{sqlOperator: "+"}}, Name: "Plus"}    //nolint:exhaustive // all present
	Minus   = FunctionByDialector{functions: map[Dialector]Function{all: OperatorFunction{sqlOperator: "-"}}, Name: "Minus"}   //nolint:exhaustive // all present
	Times   = FunctionByDialector{functions: map[Dialector]Function{all: OperatorFunction{sqlOperator: "*"}}, Name: "Times"}   //nolint:exhaustive // all present
	Divided = FunctionByDialector{functions: map[Dialector]Function{all: OperatorFunction{sqlOperator: "/"}}, Name: "Divided"} //nolint:exhaustive // all present
	Modulo  = FunctionByDialector{functions: map[Dialector]Function{all: OperatorFunction{sqlOperator: "%"}}, Name: "Modulo"}  //nolint:exhaustive // all present
	Power   = FunctionByDialector{functions: map[Dialector]Function{                                                           //nolint:exhaustive // all present
		Postgres: OperatorFunction{sqlOperator: "^"},
		all:      FunctionFunction{sqlFunction: "POWER"},
	}, Name: "Power"}
	SquareRoot = FunctionByDialector{functions: map[Dialector]Function{ //nolint:exhaustive // all present
		all: FunctionFunction{sqlFunction: "SQRT"},
	}, Name: "SquareRoot"}
	Absolute = FunctionByDialector{functions: map[Dialector]Function{ //nolint:exhaustive // all present
		all: FunctionFunction{sqlFunction: "ABS"},
	}, Name: "Absolute"}
	BitAnd = FunctionByDialector{functions: map[Dialector]Function{all: OperatorFunction{sqlOperator: "&"}}, Name: "And"} //nolint:exhaustive // all present
	BitOr  = FunctionByDialector{functions: map[Dialector]Function{all: OperatorFunction{sqlOperator: "|"}}, Name: "Or"}  //nolint:exhaustive // all present
	BitXor = FunctionByDialector{functions: map[Dialector]Function{                                                       //nolint:exhaustive // supported
		Postgres:  OperatorFunction{sqlOperator: "#"},
		MySQL:     OperatorFunction{sqlOperator: "^"},
		SQLServer: OperatorFunction{sqlOperator: "^"},
	}, Name: "Xor"}
	BitNot = FunctionByDialector{
		functions: map[Dialector]Function{all: PreOperatorFunction{sqlOperator: "~"}}, //nolint:exhaustive // all present
		Name:      "Not",
	}
	BitShiftLeft = FunctionByDialector{
		functions: map[Dialector]Function{all: OperatorFunction{sqlOperator: "<<"}}, //nolint:exhaustive // all present
		Name:      "ShiftLeft",
	}
	BitShiftRight = FunctionByDialector{
		functions: map[Dialector]Function{all: OperatorFunction{sqlOperator: ">>"}}, //nolint:exhaustive // all present
		Name:      "ShiftRight",
	}

	// String

	Concat = FunctionByDialector{
		functions: map[Dialector]Function{ //nolint:exhaustive // all present
			Postgres: OperatorFunction{sqlOperator: "||"},
			all:      FunctionFunction{sqlFunction: "CONCAT"},
		},
		Name: "Concat",
	}

	// Aggregators
	// All

	Count = FunctionByDialector{
		functions: map[Dialector]Function{all: FunctionFunction{sqlFunction: "COUNT"}}, //nolint:exhaustive // all present
		Name:      "Count",
	}
	CountAll = FunctionByDialector{
		functions: map[Dialector]Function{all: FunctionFunction{sqlPrefix: "COUNT(*)"}}, //nolint:exhaustive // all present
		Name:      "CountAll",
	}
	Min = FunctionByDialector{
		functions: map[Dialector]Function{all: FunctionFunction{sqlFunction: "MIN"}}, //nolint:exhaustive // all present
		Name:      "Min",
	}
	Max = FunctionByDialector{
		functions: map[Dialector]Function{all: FunctionFunction{sqlFunction: "MAX"}}, //nolint:exhaustive // all present
		Name:      "Max",
	}

	// Numeric

	Sum = FunctionByDialector{
		functions: map[Dialector]Function{all: FunctionFunction{sqlFunction: "SUM"}}, //nolint:exhaustive // all present
		Name:      "Sum",
	}
	Average = FunctionByDialector{
		functions: map[Dialector]Function{all: FunctionFunction{sqlFunction: "AVG"}}, //nolint:exhaustive // all present
		Name:      "Average",
	}
	BitAndAggregation = FunctionByDialector{
		functions: map[Dialector]Function{ //nolint:exhaustive // supported
			Postgres: FunctionFunction{sqlFunction: "BIT_AND"},
			MySQL:    FunctionFunction{sqlFunction: "BIT_AND"},
		},
		Name: "And",
	}
	BitOrAggregation = FunctionByDialector{
		functions: map[Dialector]Function{ //nolint:exhaustive // supported
			Postgres: FunctionFunction{sqlFunction: "BIT_OR"},
			MySQL:    FunctionFunction{sqlFunction: "BIT_OR"},
		},
		Name: "Or",
	}

	// Bool
	All = FunctionByDialector{
		functions: map[Dialector]Function{ //nolint:exhaustive // all present
			Postgres: FunctionFunction{sqlFunction: "EVERY"},
			all:      FunctionFunction{sqlFunction: "AVG", sqlSuffix: "== 1"},
		},
		Name: "All",
	}
	Any = FunctionByDialector{
		functions: map[Dialector]Function{ //nolint:exhaustive // all present
			Postgres: FunctionFunction{sqlFunction: "BOOL_OR"},
			all:      FunctionFunction{sqlFunction: "AVG", sqlSuffix: "> 0"},
		},
		Name: "Any",
	}
	None = FunctionByDialector{
		functions: map[Dialector]Function{ //nolint:exhaustive // all present
			Postgres: FunctionFunction{sqlFunction: "NOT BOOL_OR"},
			all:      FunctionFunction{sqlFunction: "AVG", sqlSuffix: "== 0"},
		},
		Name: "None",
	}
)

func (f FunctionByDialector) Get(dialector Dialector) (Function, bool) {
	dialectorFunc, dialectorPresent := f.functions[dialector]

	if dialectorPresent {
		return dialectorFunc, true
	}

	allFunc, allPresent := f.functions[all]

	return allFunc, allPresent
}
