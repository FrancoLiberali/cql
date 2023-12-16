package conditions

import (
	"go/types"

	"github.com/dave/jennifer/jen"
)

type JenParam struct {
	statement *jen.Statement
}

func NewJenParam() *JenParam {
	return &JenParam{
		statement: jen.Id("v"),
	}
}

func (param JenParam) Statement() *jen.Statement {
	return param.statement
}

func (param JenParam) ToBasicKind(basicType *types.Basic) {
	switch basicType.Kind() {
	case types.Bool:
		param.statement.Bool()
	case types.Int:
		param.statement.Int()
	case types.Int8:
		param.statement.Int8()
	case types.Int16:
		param.statement.Int16()
	case types.Int32:
		param.statement.Int32()
	case types.Int64:
		param.statement.Int64()
	case types.Uint:
		param.statement.Uint()
	case types.Uint8:
		param.statement.Uint8()
	case types.Uint16:
		param.statement.Uint16()
	case types.Uint32:
		param.statement.Uint32()
	case types.Uint64:
		param.statement.Uint64()
	case types.Uintptr:
		param.statement.Uintptr()
	case types.Float32:
		param.statement.Float32()
	case types.Float64:
		param.statement.Float64()
	case types.Complex64:
		param.statement.Complex64()
	case types.Complex128:
		param.statement.Complex128()
	case types.String:
		param.statement.String()
	}
}

func (param JenParam) ToSlice() {
	param.statement.Index()
}

func (param JenParam) ToCustomType(destPkg string, typeV Type) {
	param.statement.Qual(
		getRelativePackagePath(destPkg, typeV),
		typeV.Name(),
	)
}
