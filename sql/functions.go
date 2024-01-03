package sql

type Function interface {
	Name() string
	ApplyTo(internalSQL string) string
}

type FunctionFunction struct {
	name        string
	sqlFunction string
}

func (f FunctionFunction) ApplyTo(internalSQL string) string {
	return f.sqlFunction + "(" + internalSQL + ", ?)"
}

func (f FunctionFunction) Name() string {
	return f.name + " (" + f.sqlFunction + ")"
}

type OperatorFunction struct {
	name        string
	sqlOperator string
}

func (f OperatorFunction) ApplyTo(internalSQL string) string {
	return "(" + internalSQL + " " + f.sqlOperator + " ?)"
}

func (f OperatorFunction) Name() string {
	return f.name + " (" + f.sqlOperator + ")"
}

type FunctionByDialector map[Dialector]Function

const all = "all"

var (
	Plus    = FunctionByDialector{all: OperatorFunction{sqlOperator: "+", name: "Plus"}}
	Minus   = FunctionByDialector{all: OperatorFunction{sqlOperator: "-", name: "Minus"}}
	Times   = FunctionByDialector{all: OperatorFunction{sqlOperator: "*", name: "Times"}}
	Divided = FunctionByDialector{all: OperatorFunction{sqlOperator: "/", name: "Divided"}}
	Modulo  = FunctionByDialector{all: OperatorFunction{sqlOperator: "%", name: "Modulo"}}
	Power   = FunctionByDialector{all: OperatorFunction{sqlOperator: "^", name: "Power"}}
	Concat  = FunctionByDialector{
		Postgres: OperatorFunction{sqlOperator: "||", name: "Concat"},
		all:      FunctionFunction{sqlFunction: "CONCAT", name: "Concat"},
	}
)

func (f FunctionByDialector) Get(dialector Dialector) (Function, bool) {
	dialectorFunc, dialectorPresent := f[dialector]

	if dialectorPresent {
		return dialectorFunc, true
	}

	allFunc, allPresent := f[all]

	return allFunc, allPresent
}
