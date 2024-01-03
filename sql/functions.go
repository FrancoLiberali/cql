package sql

import "strings"

type Function interface {
	ApplyTo(internalSQL string, values int) string
}

type FunctionFunction struct {
	sqlFunction string
}

func (f FunctionFunction) ApplyTo(internalSQL string, values int) string {
	placeholders := strings.Repeat(", ?", values)

	return f.sqlFunction + "(" + internalSQL + placeholders + ")"
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
		Postgres: PreOperatorFunction{sqlOperator: "|/"},
		all:      FunctionFunction{sqlFunction: "SQRT"},
	}, Name: "SquareRoot"}
	Absolute = FunctionByDialector{functions: map[Dialector]Function{ //nolint:exhaustive // all present
		Postgres: PreOperatorFunction{sqlOperator: "@"},
		all:      FunctionFunction{sqlFunction: "abs"},
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
	Concat = FunctionByDialector{
		functions: map[Dialector]Function{ //nolint:exhaustive // all present
			Postgres: OperatorFunction{sqlOperator: "||"},
			all:      FunctionFunction{sqlFunction: "CONCAT"},
		},
		Name: "Concat",
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
