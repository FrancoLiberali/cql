package condition

import (
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Scanner is a code-generated row materializer for model type T. It avoids
// per-row reflection: ScanValues returns a []any of concrete pointer types
// matching the result columns (so database/sql.Rows.Scan does its native
// driver conversion), and AssignValues consumes that same slice with cheap
// type assertions and direct struct-field stores.
//
// Generated code (from cql-gen) creates one Scanner[T] per model and shares
// a pointer to it from every Field on that model's conditions struct.
// Query[T] picks it up from any condition built off those fields and
// dispatches the fast Find/First/Take/Last paths through it.
type Scanner[T any] struct {
	ScanValues   func(columns []string) ([]any, error)
	AssignValues func(dest *T, columns []string, values []any) error
}

// NullSink is a sql.Scanner that discards any value. Generated ScanValues
// uses it as a safe placeholder for columns the model doesn't know about
// (e.g. extra projections added by Order in Postgres), so the row can be
// drained cleanly without erroring.
type NullSink struct{}

func (*NullSink) Scan(_ any) error { return nil }

var _ sql.Scanner = (*NullSink)(nil)

// NullableScanner wraps any sql.Scanner so NULL values from the driver are
// silently skipped instead of forwarded to the inner Scan. Generated scanners
// use this for user-defined types that implement sql.Scanner but don't
// themselves accept nil (the common case — see models.MultiString as an
// example). gorm's reflective scan path applies the same filter before
// invoking field-level setters, so wrapping here keeps behavior identical.
type NullableScanner struct {
	Inner sql.Scanner
}

func (n *NullableScanner) Scan(src any) error {
	if src == nil {
		return nil
	}

	return n.Inner.Scan(src)
}

var _ sql.Scanner = (*NullableScanner)(nil)

// scannerProvider is the optional interface implemented by Conditions that
// can surface a scanner. Field-backed conditions return the scanner stored
// on their field; combinators (AND/OR/NOT) delegate to their children.
type scannerProvider interface {
	getScannerErased() any
}

// canUseFastScan reports whether the query is safe to materialize via
// Scanner. Bails when either:
//   - gorm Preload is in use (separate queries that populate relation slices)
//   - a joined preload added columns from a non-initial table (those columns
//     are meant to populate a relation field on T, which Scanner doesn't model)
func canUseFastScan(q *CQLQuery) bool {
	if q == nil || q.gormDB == nil || q.gormDB.Statement == nil {
		return false
	}

	return len(q.gormDB.Statement.Preloads) == 0 && !q.hasJoinedSelects
}

// findWith materializes every matching row into *dest using scanner. Caller
// must have validated that the fast path applies (scanner != nil and
// canUseFastScan(q) == true).
func findWith[T any](q *CQLQuery, dest *[]*T, scanner *Scanner[T]) error {
	rows, err := q.gormDB.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		values, err := scanner.ScanValues(columns)
		if err != nil {
			return err
		}

		if err := rows.Scan(values...); err != nil {
			return err
		}

		var m T
		if err := scanner.AssignValues(&m, columns, values); err != nil {
			return err
		}

		*dest = append(*dest, &m)
	}

	return rows.Err()
}

// firstWith materializes the first row (primary-key ordered) into dest.
// Returns gorm.ErrRecordNotFound when no row matches, matching the gorm
// fallback path's contract.
func firstWith[T any](q *CQLQuery, dest **T, scanner *Scanner[T]) error {
	return scanOne(q.gormDB.Order(clause.OrderByColumn{Column: clause.PrimaryColumn}), dest, scanner)
}

// takeWith materializes any one matching row (no ordering) into dest.
func takeWith[T any](q *CQLQuery, dest **T, scanner *Scanner[T]) error {
	return scanOne(q.gormDB, dest, scanner)
}

// lastWith materializes the last row (primary-key DESC) into dest.
func lastWith[T any](q *CQLQuery, dest **T, scanner *Scanner[T]) error {
	return scanOne(q.gormDB.Order(clause.OrderByColumn{Column: clause.PrimaryColumn, Desc: true}), dest, scanner)
}

func scanOne[T any](db *gorm.DB, dest **T, scanner *Scanner[T]) error {
	rows, err := db.Limit(1).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}

		return gorm.ErrRecordNotFound
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values, err := scanner.ScanValues(columns)
	if err != nil {
		return err
	}

	if err := rows.Scan(values...); err != nil {
		return err
	}

	var m T
	if err := scanner.AssignValues(&m, columns, values); err != nil {
		return err
	}

	*dest = &m

	return rows.Err()
}

