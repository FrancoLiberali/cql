package condition

import (
	"reflect"

	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

// scannerRegistry maps a model type to its generated *Scanner[T]. It is
// written only from the init() functions cql-gen emits in each *_scanner.go
// (which all run before main, single-threaded), and read-only afterwards, so a
// plain map is safe for the concurrent reads that happen during queries.
var scannerRegistry = map[reflect.Type]any{}

// RegisterScanner records the per-model Scanner so a Query built without any
// conditions (e.g. Query[T](ctx, db).Find()) can still take the fast-scan path
// by resolving the scanner by type. Generated *_scanner.go files call this
// from init().
func RegisterScanner[T model.Model](s *Scanner[T]) {
	var zero T

	scannerRegistry[reflect.TypeOf(zero)] = s
}

type Query[T model.Model] struct {
	cqlQuery *CQLQuery
	scanner  *Scanner[T]
	err      error
}

func (query *Query[T]) getError() error {
	return query.err
}

func (query *Query[T]) getCQLQuery() *CQLQuery {
	return query.cqlQuery
}

// Ascending specify an ascending order when retrieving models from database
func (query *Query[T]) Ascending(field IField) *Query[T] {
	return query.order(field, false)
}

// Descending specify a descending order when retrieving models from database
func (query *Query[T]) Descending(field IField) *Query[T] {
	return query.order(field, true)
}

// Order specify order when retrieving models from database
// if descending is true, the ordering is in descending direction
func (query *Query[T]) order(field IField, descending bool) *Query[T] {
	err := query.cqlQuery.Order(field, descending)
	if err != nil {
		methodName := "Ascending"
		if descending {
			methodName = "Descending"
		}

		query.addError(methodError(err, methodName))
	}

	return query
}

// Limit specify the number of models to be retrieved
//
// Limit conditions can be cancelled by using `Limit(-1)`
func (query *Query[T]) Limit(limit int) *Query[T] {
	query.cqlQuery.Limit(limit)

	return query
}

// Offset specify the number of models to skip before starting to return the results
//
// Offset conditions can be cancelled by using `Offset(-1)`
//
// Warning: in MySQL Offset can only be used if Limit is also used
func (query *Query[T]) Offset(offset int) *Query[T] {
	query.cqlQuery.Offset(offset)

	return query
}

// GroupBy arrange identical data into groups
func (query *Query[T]) GroupBy(fields ...IField) *QueryGroup {
	query.addError(query.cqlQuery.GroupBy(fields))

	return &QueryGroup{
		cqlQuery: query.cqlQuery,
		err:      query.err,
		fields:   fields,
	}
}

// Finishing methods

// Count returns the amount of models that fulfill the conditions
func (query *Query[T]) Count() (int64, error) {
	if query.err != nil {
		return 0, query.err
	}

	return query.cqlQuery.Count()
}

// First finds the first model ordered by primary key, matching given conditions
// or returns gorm.ErrRecordNotFound is if no model does it
func (query *Query[T]) First() (*T, error) {
	if query.err != nil {
		return nil, query.err
	}

	if query.scanner != nil && canUseFastScan(query.cqlQuery) {
		var model *T

		return model, firstWith[T](query.cqlQuery, &model, query.scanner)
	}

	var model *T

	return model, query.cqlQuery.First(&model)
}

// Take finds the first model returned by the database in no specified order, matching given conditions
// or returns gorm.ErrRecordNotFound is if no model does it
func (query *Query[T]) Take() (*T, error) {
	if query.err != nil {
		return nil, query.err
	}

	if query.scanner != nil && canUseFastScan(query.cqlQuery) {
		var model *T

		return model, takeWith[T](query.cqlQuery, &model, query.scanner)
	}

	var model *T

	return model, query.cqlQuery.Take(&model)
}

// Last finds the last model ordered by primary key, matching given conditions
// or returns gorm.ErrRecordNotFound is if no model does it
func (query *Query[T]) Last() (*T, error) {
	if query.err != nil {
		return nil, query.err
	}

	if query.scanner != nil && canUseFastScan(query.cqlQuery) {
		var model *T

		return model, lastWith[T](query.cqlQuery, &model, query.scanner)
	}

	var model *T

	return model, query.cqlQuery.Last(&model)
}

// FindOne finds the only one model that matches given conditions
// or returns error if 0 or more than 1 are found.
func (query *Query[T]) FindOne() (*T, error) {
	models, err := query.Find()
	if err != nil {
		return nil, err
	}

	switch {
	case len(models) == 1:
		return models[0], nil
	case len(models) == 0:
		return nil, ErrObjectNotFound
	default:
		return nil, ErrMoreThanOneObjectFound
	}
}

// Find finds all models matching given conditions
func (query *Query[T]) Find() ([]*T, error) {
	if query.err != nil {
		return nil, query.err
	}

	if query.scanner != nil && canUseFastScan(query.cqlQuery) {
		var models []*T

		return models, findWith[T](query.cqlQuery, &models, query.scanner)
	}

	var models []*T

	return models, query.cqlQuery.Find(&models)
}

func (query *Query[T]) addError(err error) {
	if err != nil && query.err == nil {
		query.err = err
	}
}

// Create a Query to which the conditions are applied inside transaction tx
func NewQuery[T model.Model](tx *gorm.DB, conditions ...Condition[T]) *Query[T] {
	gormQuery, err := ApplyConditions(tx, conditions)

	return &Query[T]{
		cqlQuery: gormQuery,
		scanner:  resolveScanner[T](conditions),
		err:      err,
	}
}

// resolveScanner walks the conditions once at query construction looking for
// the first one that carries a *Scanner[T]. All conditions built from the
// generated conditions struct share the same scanner pointer, so the first
// match is enough.
func resolveScanner[T model.Model](conditions []Condition[T]) *Scanner[T] {
	for _, c := range conditions {
		sp, ok := c.(scannerProvider)
		if !ok {
			continue
		}

		if s, ok := sp.getScannerErased().(*Scanner[T]); ok && s != nil {
			return s
		}
	}

	// No condition carried the scanner (e.g. a no-condition query) — fall back
	// to the per-type registry populated by generated init()s.
	var zero T
	if s, ok := scannerRegistry[reflect.TypeOf(zero)].(*Scanner[T]); ok {
		return s
	}

	return nil
}
