package cmd

import (
	"errors"
	"go/types"

	"github.com/dave/jennifer/jen"
)

type JenParam struct {
	internalType *jen.Statement
	isBool       bool
	isString     bool
	isSlice      bool
	isNumeric    bool
}

func NewJenParam() *JenParam {
	return &JenParam{
		internalType: &jen.Statement{},
	}
}

func (param JenParam) GenericType() *jen.Statement {
	return param.internalType
}

func (param *JenParam) ToBasicKind(basicType *types.Basic) {
	switch basicType.Kind() {
	case types.Bool:
		param.ToBool()
	case types.Int:
		param.ToNumeric(param.internalType.Int())
	case types.Int8:
		param.ToNumeric(param.internalType.Int8())
	case types.Int16:
		param.ToNumeric(param.internalType.Int16())
	case types.Int32:
		param.ToNumeric(param.internalType.Int32())
	case types.Int64:
		param.ToNumeric(param.internalType.Int64())
	case types.Uint:
		param.ToNumeric(param.internalType.Uint())
	case types.Uint8:
		param.ToNumeric(param.internalType.Uint8())
	case types.Uint16:
		param.ToNumeric(param.internalType.Uint16())
	case types.Uint32:
		param.ToNumeric(param.internalType.Uint32())
	case types.Uint64:
		param.ToNumeric(param.internalType.Uint64())
	case types.Uintptr:
		param.internalType.Uintptr()
	case types.Float32:
		param.ToNumeric(param.internalType.Float32())
	case types.Float64:
		param.ToNumeric(param.internalType.Float64())
	case types.Complex64:
		param.internalType.Complex64()
	case types.Complex128:
		param.internalType.Complex128()
	case types.String:
		param.ToString()
	case types.Invalid, types.UnsafePointer,
		types.UntypedBool, types.UntypedInt,
		types.UntypedRune, types.UntypedFloat,
		types.UntypedComplex, types.UntypedString,
		types.UntypedNil:
		panic(errors.New("unreachable! untyped types can't be inside a struct"))
	}
}

func (param *JenParam) ToSlice() {
	param.isSlice = true
	param.internalType.Index()
}

func (param JenParam) ToCustomType(destPkg string, typeV Type) {
	param.internalType.Qual(
		getRelativePackagePath(destPkg, typeV),
		typeV.Name(),
	)
}

func (param *JenParam) SQLToBasicType(typeV Type) {
	switch typeV.String() {
	case nullString:
		param.ToString()
	case nullInt64:
		param.ToNumeric(param.internalType.Int64())
	case nullInt32:
		param.ToNumeric(param.internalType.Int32())
	case nullInt16:
		param.ToNumeric(param.internalType.Int16())
	case nullByte:
		param.ToNumeric(param.internalType.Int8())
	case nullFloat64:
		param.ToNumeric(param.internalType.Float64())
	case nullBool:
		param.ToBool()
	case nullTime, deletedAt:
		param.internalType.Qual(
			"time",
			"Time",
		)
	}
}

func (param *JenParam) ToBool() {
	if !param.isSlice {
		param.isBool = true
	}

	param.internalType.Bool()
}

func (param *JenParam) ToString() {
	if !param.isSlice {
		param.isString = true
	}

	param.internalType.String()
}

func (param *JenParam) ToNumeric(*jen.Statement) {
	if !param.isSlice {
		param.isNumeric = true
	}
}
