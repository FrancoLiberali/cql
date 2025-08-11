package cql

import "github.com/FrancoLiberali/cql/condition"

func Int(value int) condition.NumericValue[int] {
	return condition.NumericValue[int]{Value: value}
}

func Bool(value bool) condition.BoolValue {
	return condition.BoolValue{Value: value}
}
