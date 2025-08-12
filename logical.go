package cql

import (
	"github.com/elliotchance/pie/v2"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
	"github.com/FrancoLiberali/cql/unsafe"
)

// Logical Operators
// ref: https://www.postgresql.org/docs/current/functions-logical.html

// And allows the connection of multiple conditions by the AND logical connector.
//
// Its use is optional as it is the default connector.
//
// Example:
//
// cql.And(conditions.City.Name.Is().Eq("Paris"), conditions.City.ZipCode.Is().Eq("75000"))
func And[T model.Model](firstCondition condition.WhereCondition[T], conditions ...condition.WhereCondition[T]) condition.WhereCondition[T] {
	return condition.And(pie.Unshift(conditions, firstCondition)...)
}

// Or allows the connection of multiple conditions by the OR logical connector.
//
// Example:
//
// cql.Or(conditions.City.Name.Is().Eq("Paris"), conditions.City.Name.Is().Eq("Buenos Aires"))
func Or[T model.Model](firstCondition condition.WhereCondition[T], conditions ...condition.WhereCondition[T]) condition.WhereCondition[T] {
	return condition.Or(pie.Unshift(conditions, firstCondition)...)
}

// Not allows the negation of the conditions within it. Multiple conditions are connected by an AND by default.
//
// Example:
//
// cql.Not(conditions.City.Name.Is().Eq("Paris"), conditions.City.Name.Is().Eq("Buenos Aires"))
//
// translates as
//
// NOT (name = "Paris" AND name = "Buenos Aires")
func Not[T model.Model](firstCondition condition.WhereCondition[T], conditions ...condition.WhereCondition[T]) condition.WhereCondition[T] {
	return condition.Not(pie.Unshift(conditions, firstCondition)...)
}

// True represents a condition that is always true.
//
// In general, it should not be used. It can only be useful in case you want to perform an operation on all models in a table.
func True[T model.Model]() condition.Condition[T] {
	return unsafe.NewCondition[T]("1 = 1")
}

// Logical Operators for having

// And allows the connection of multiple conditions by the AND logical connector.
//
// Its use is optional as it is the default connector.
//
// Example:
//
// cql.And(conditions.City.Name.Is().Eq("Paris"), conditions.City.ZipCode.Is().Eq("75000"))
func AndHaving(firstCondition condition.AggregationCondition, conditions ...condition.AggregationCondition) condition.AggregationCondition {
	return condition.ConnectionAggregationCondition(
		pie.Unshift(conditions, firstCondition),
		sql.And,
	)
}

// Or allows the connection of multiple conditions by the OR logical connector.
//
// Example:
//
// cql.Or(conditions.City.Name.Is().Eq("Paris"), conditions.City.Name.Is().Eq("Buenos Aires"))
func OrHaving(firstCondition condition.AggregationCondition, conditions ...condition.AggregationCondition) condition.AggregationCondition {
	return condition.ConnectionAggregationCondition(
		pie.Unshift(conditions, firstCondition),
		sql.Or,
	)
}

// Not allows the negation of the conditions within it. Multiple conditions are connected by an AND by default.
//
// Example:
//
// cql.Not(conditions.City.Name.Is().Eq("Paris"), conditions.City.Name.Is().Eq("Buenos Aires"))
//
// translates as
//
// NOT (name = "Paris" AND name = "Buenos Aires")
// TODO docs
func NotHaving(firstCondition condition.AggregationCondition, conditions ...condition.AggregationCondition) condition.AggregationCondition {
	return condition.ContainerAggregationCondition(
		pie.Unshift(conditions, firstCondition),
		sql.Not,
	)
}
