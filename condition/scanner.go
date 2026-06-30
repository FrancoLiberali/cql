package condition

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/FrancoLiberali/cql/model"
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

// RelationScanner is a code-generated mounter for a single relation between
// Parent and Child. It reuses the existing per-model Scanner[Child] for
// column materialization and adds typed hooks for mounting (or nilling) the
// child onto the parent.
//
// Generated code wires one RelationScanner per relation onto its matching
// generated JoinCondition. When the user calls .Preload() on that condition
// at runtime, the JoinCondition registers the scanner against its
// runtime-computed alias (e.g. "Seller", "Seller__Company") so findWith
// dispatches joined columns through it.
type RelationScanner[Parent, Child any] struct {
	// RelationField is the Go field name on Parent (e.g. "Brand"). Matches
	// what CQL's Table.DeliverTable uses as the first-level alias.
	RelationField string
	// ChildScanner provides ScanValues + AssignValues for Child's columns.
	// Same instance the user gets via Query[Child].
	ChildScanner *Scanner[Child]
	// Mount sets the (non-nil) child onto parent. Pointer-relation
	// generators set `p.X = c`; value-relation generators set `p.X = *c`.
	Mount func(parent *Parent, child *Child)
	// SetNil resets the relation field on parent. For pointer relations
	// this sets to nil; for value relations it is a no-op because the
	// LEFT-JOIN-no-match case leaves the field at its zero value (probed
	// against gorm in TestProbe_ValueRelationNoMatch).
	SetNil func(parent *Parent)
}

// activeJoin is the type-erased runtime form of a RelationScanner bound to
// a specific alias path. CQLQuery accumulates these when joined preloads
// are applied.
type activeJoin struct {
	// alias is the column-name prefix CQL generates for this join's
	// columns (without the trailing "__"). E.g. "Seller", "Seller__Company".
	alias string
	// parentAlias points to the activeJoin that owns the mount target for
	// this one; empty for direct (top-level) joins which mount onto the
	// main row.
	parentAlias string
	// alloc creates a Child instance + its scan-value slots for the
	// columns matched to this join.
	alloc func(cols []string) (child any, vals []any, err error)
	// assign runs the child scanner against the scanned values and
	// reports whether every cell came back NULL (used to detect the
	// LEFT-JOIN-no-match case).
	assign func(child any, cols []string, vals []any) (allNull bool, err error)
	// mount writes the populated (or nil) child onto its parent
	// instance. parent is the typed-erased pointer (*MainT for direct
	// joins, the child instance of another activeJoin for nested).
	mount func(parent any, child any, allNull bool) error
}

// registerActiveJoin attaches an activeJoin to the query so findWith picks
// it up. Called from joinConditionImpl.applyTo when the user requested
// .Preload() on a generated JoinCondition.
func registerActiveJoin[Parent, Child any](
	q *CQLQuery, alias, parentAlias string, rs *RelationScanner[Parent, Child],
) {
	q.activeJoins = append(q.activeJoins, &activeJoin{
		alias:       alias,
		parentAlias: parentAlias,
		alloc: func(cols []string) (any, []any, error) {
			vals, err := rs.ChildScanner.ScanValues(cols)
			if err != nil {
				return nil, nil, err
			}

			return new(Child), vals, nil
		},
		assign: func(child any, cols []string, vals []any) (bool, error) {
			c, ok := child.(*Child)
			if !ok {
				return false, fmt.Errorf("cql join scan: expected *%T", *new(Child))
			}

			if err := rs.ChildScanner.AssignValues(c, cols, vals); err != nil {
				return false, err
			}

			return allValuesNull(vals), nil
		},
		mount: func(parent any, child any, allNull bool) error {
			p, ok := parent.(*Parent)
			if !ok {
				return fmt.Errorf("cql join mount: expected *%T parent, got %T", *new(Parent), parent)
			}

			if allNull {
				if rs.SetNil != nil {
					rs.SetNil(p)
				}

				return nil
			}

			c, ok := child.(*Child)
			if !ok {
				return fmt.Errorf("cql join mount: expected *%T child, got %T", *new(Child), child)
			}

			rs.Mount(p, c)

			return nil
		},
	})
}

// allValuesNull checks whether every scanned cell carries the NULL sentinel
// for its concrete sql.Null* type. Returns true only when every cell either
// is a non-Valid sql.Null* or is nil. Used to detect LEFT-JOIN-no-match
// rows so we don't mount a zero-valued child.
func allValuesNull(values []any) bool {
	for _, v := range values {
		switch t := v.(type) {
		case *sql.NullBool:
			if t.Valid {
				return false
			}
		case *sql.NullString:
			if t.Valid {
				return false
			}
		case *sql.NullInt64:
			if t.Valid {
				return false
			}
		case *sql.NullInt32:
			if t.Valid {
				return false
			}
		case *sql.NullInt16:
			if t.Valid {
				return false
			}
		case *sql.NullByte:
			if t.Valid {
				return false
			}
		case *sql.NullFloat64:
			if t.Valid {
				return false
			}
		case *sql.NullTime:
			if t.Valid {
				return false
			}
		case *gorm.DeletedAt:
			if t.Valid {
				return false
			}
		case *NullSink:
			// always "null" sentinel — ignore
		case *NullableScanner:
			// custom scanners; we have no .Valid to inspect, so be
			// conservative: treat as non-null.
			return false
		case *model.UUID:
			// CQL convention: NilUUID is the sentinel for "no value".
			// Generated scanners scan a UUID column directly into
			// model.UUID; a LEFT-JOIN-no-match leaves it at NilUUID.
			if *t != model.NilUUID {
				return false
			}
		default:
			// unknown type (custom scanner direct, gorm types, etc.) —
			// be conservative: assume non-Null-wrapped values had data.
			return false
		}
	}

	return true
}

// HasManyLoader is a code-generated post-scan mounter for a HasMany
// relation. After the main query materializes parents, the runtime walks
// every registered loader exactly once:
//
//   - collect parent PK values via ParentID
//   - build a Query[Child] (which gets all the fast-path goodies, including
//     a nested fast scanner for Child) via BuildQuery and run Find
//   - group results by ChildFK
//   - call Mount(parent, groupedChildren) per parent
//
// Generated code constructs one HasManyLoader per HasMany relation on each
// parent model, with all the typed accessors baked in. The cost is one
// extra SQL round-trip per HasMany relation per query (same shape gorm's
// reflective preload uses).
type HasManyLoader[Parent, Child model.Model] struct {
	// CollectionField is the Go field name on Parent that holds the
	// children (e.g. "Sellers"). Used for diagnostic output.
	CollectionField string
	// ParentID extracts the PK value from a parent. Compared against
	// ChildFK to group children with their parent.
	ParentID func(*Parent) any
	// ChildFK extracts the foreign-key value from a child (e.g. its
	// CompanyID for a Seller→Company relation).
	ChildFK func(*Child) any
	// Mount writes the grouped children onto parent. Pointer slices
	// (`*[]Seller`) and value slices (`[]Seller`) both use Mount —
	// generated code adapts.
	Mount func(parent *Parent, children []*Child)
	// BuildQuery constructs the SELECT-children-WHERE-fk-IN(...) query.
	// Generated code uses the child's conditions struct so nested
	// preloads/joins/filters compose naturally:
	//   func(tx *gorm.DB, ids []any, nested []Condition[Seller]) (*Query[Seller], error) {
	//       conds := append(
	//         []Condition[Seller]{ conditions.Seller.CompanyID.IsUnsafe().In(ids) },
	//         nested...,
	//       )
	//       return NewQuery[Seller](tx, conds...), nil
	//   }
	BuildQuery func(tx *gorm.DB, parentIDs []any, nested []Condition[Child]) (*Query[Child], error)
}

// activeHasMany is the type-erased runtime form of a HasManyLoader bound
// to its registered nested preloads. CQLQuery accumulates these when
// collectionPreloadCondition.applyTo runs.
type activeHasMany struct {
	// collectionField is the Go field name on the parent model that
	// receives the children — purely for error messages.
	collectionField string
	// run executes the loader against the given parents. parents is
	// []*Parent type-erased to []any (so the runtime doesn't have to
	// be generic in Parent here).
	run func(tx *gorm.DB, parents []any) error
}

// registerHasManyLoader wraps a HasManyLoader into the type-erased
// activeHasMany form and registers it on the query. Called from
// collectionPreloadCondition.applyTo when generated code supplied a
// loader. nested are the JoinConditions passed to .Preload(nested...) so
// the child query can preload its own relations in turn.
func registerHasManyLoader[Parent, Child model.Model](
	q *CQLQuery,
	loader *HasManyLoader[Parent, Child],
	nested []Condition[Child],
) {
	q.activeHasMany = append(q.activeHasMany, &activeHasMany{
		collectionField: loader.CollectionField,
		run: func(tx *gorm.DB, parents []any) error {
			if len(parents) == 0 {
				return nil
			}

			ids := make([]any, 0, len(parents))

			for _, p := range parents {
				pp, ok := p.(*Parent)
				if !ok {
					return fmt.Errorf("cql hasmany %q: expected *%T parent, got %T",
						loader.CollectionField, *new(Parent), p)
				}

				ids = append(ids, loader.ParentID(pp))
			}

			// Fresh session: otherwise the child query inherits the
			// parent query's FROM clause + WHERE conditions, producing
			// nonsense SQL like `SELECT sellers.* FROM companies WHERE
			// companies.deleted_at IS NULL AND sellers.company_id IN (...)`.
			childTx := tx.Session(&gorm.Session{NewDB: true, Context: tx.Statement.Context})

			childQuery, err := loader.BuildQuery(childTx, ids, nested)
			if err != nil {
				return err
			}

			children, err := childQuery.Find()
			if err != nil {
				return err
			}

			groups := make(map[any][]*Child, len(parents))
			for _, c := range children {
				fk := loader.ChildFK(c)
				groups[fk] = append(groups[fk], c)
			}

			for _, p := range parents {
				pp := p.(*Parent)
				loader.Mount(pp, groups[loader.ParentID(pp)])
			}

			return nil
		},
	})
}

// runHasManyLoaders is called by findWith / scanOne after the main rows
// are materialized. Type-erases the typed *T list into []any for each
// registered loader.
func runHasManyLoaders[T any](q *CQLQuery, parents []*T) error {
	if len(q.activeHasMany) == 0 {
		return nil
	}

	asAny := make([]any, len(parents))
	for i, p := range parents {
		asAny[i] = p
	}

	for _, hm := range q.activeHasMany {
		if err := hm.run(q.gormDB, asAny); err != nil {
			return err
		}
	}

	return nil
}

// scannerProvider is the optional interface implemented by Conditions that
// can surface a scanner. Field-backed conditions return the scanner stored
// on their field; combinators (AND/OR/NOT) delegate to their children.
type scannerProvider interface {
	getScannerErased() any
}

// canUseFastScan reports whether the query is safe to materialize via the
// generated Scanner(s). Bails when:
//   - gorm Preload is in use AND no activeHasMany loader covers it (HasMany
//     without a generated loader falls back to gorm; with loader, findWith
//     strips Preloads before executing the main query and runs the loader
//     post-scan)
//   - the query added joined selects but no RelationScanner was registered
//     for any of them
func canUseFastScan(q *CQLQuery) bool {
	if q == nil || q.gormDB == nil || q.gormDB.Statement == nil {
		return false
	}

	if len(q.gormDB.Statement.Preloads) > 0 && len(q.activeHasMany) == 0 {
		return false
	}

	// Joined selects without an activeJoin to consume them → fall back.
	// (E.g. a join-for-filter that promotes selects for Postgres ordering,
	// or a future preload path we don't yet generate.)
	if q.hasJoinedSelects && len(q.activeJoins) == 0 {
		return false
	}

	return true
}

// stripGormPreloads removes any gorm preloads from the statement before
// the main SELECT runs — the activeHasMany loaders will populate the
// relation slices post-scan instead, avoiding a duplicate gorm-driven
// child query. The Preloads stay registered for non-Find paths (UPDATE
// RETURNING etc.) that don't go through findWith.
func stripGormPreloads(q *CQLQuery) {
	if q == nil || q.gormDB == nil || q.gormDB.Statement == nil {
		return
	}

	q.gormDB.Statement.Preloads = nil
}

// findWith materializes every matching row into *dest using scanner. Caller
// must have validated that the fast path applies (scanner != nil and
// canUseFastScan(q) == true). After the main scan, runs registered
// HasMany loaders (one extra SELECT per relation) and mounts the grouped
// children onto each parent.
func findWith[T any](q *CQLQuery, dest *[]*T, scanner *Scanner[T]) error {
	// HasMany loaders own the child fetch; if gorm Preloads are also
	// registered (kept alive for non-Find paths) strip them here so gorm
	// doesn't duplicate the work + override the loader's mount.
	if len(q.activeHasMany) > 0 {
		stripGormPreloads(q)
	}

	rows, err := q.gormDB.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	plan, err := buildScanPlan(columns, q.activeJoins)
	if err != nil {
		return err
	}

	for rows.Next() {
		m, err := scanOneRow[T](rows, plan, columns, scanner)
		if err != nil {
			return err
		}

		*dest = append(*dest, m)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return runHasManyLoaders[T](q, *dest)
}

// firstWith materializes the first row (primary-key ordered) into dest.
// Returns gorm.ErrRecordNotFound when no row matches, matching the gorm
// fallback path's contract.
func firstWith[T any](q *CQLQuery, dest **T, scanner *Scanner[T]) error {
	return scanOne(q, q.gormDB.Order(clause.OrderByColumn{Column: clause.PrimaryColumn}), dest, scanner)
}

// takeWith materializes any one matching row (no ordering) into dest.
func takeWith[T any](q *CQLQuery, dest **T, scanner *Scanner[T]) error {
	return scanOne(q, q.gormDB, dest, scanner)
}

// lastWith materializes the last row (primary-key DESC) into dest.
func lastWith[T any](q *CQLQuery, dest **T, scanner *Scanner[T]) error {
	return scanOne(q, q.gormDB.Order(clause.OrderByColumn{Column: clause.PrimaryColumn, Desc: true}), dest, scanner)
}

// scanOne materializes a single row + runs registered HasMany loaders
// against the singleton parent.
func scanOne[T any](q *CQLQuery, db *gorm.DB, dest **T, scanner *Scanner[T]) error {
	if len(q.activeHasMany) > 0 {
		stripGormPreloads(q)
	}

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

	plan, err := buildScanPlan(columns, q.activeJoins)
	if err != nil {
		return err
	}

	m, err := scanOneRow[T](rows, plan, columns, scanner)
	if err != nil {
		return err
	}

	*dest = m

	if err := rows.Err(); err != nil {
		return err
	}

	return runHasManyLoaders[T](q, []*T{m})
}

// scanPlan precomputes per-column routing once (instead of doing prefix
// matching per row × per column). For each column in the result set,
// columnOwner records which active join (-1 = main) owns it, and joinCols /
// mainCols hold the per-owner column subset in order so each scanner sees
// the same column slice across all rows.
type scanPlan struct {
	mainColsIdx []int                 // result indices owned by the main scanner
	mainCols    []string              // names of main columns, parallel to mainColsIdx
	joins       []*activeJoin         // sorted parent-first (top-level → nested)
	joinColsIdx [][]int               // per-join result indices
	joinCols    [][]string            // per-join column names
}

// buildScanPlan inspects the row's columns once and produces a routing plan.
// Joined columns must be assigned to the deepest matching alias (so
// "Seller__Company__id" goes to the nested join, not the parent), which is
// achieved by sorting alias prefixes longest-first when matching.
func buildScanPlan(columns []string, joins []*activeJoin) (*scanPlan, error) {
	plan := &scanPlan{
		joins:       joins,
		joinColsIdx: make([][]int, len(joins)),
		joinCols:    make([][]string, len(joins)),
	}

	// Build (idx, prefix) pairs sorted by prefix length descending so
	// the deepest alias wins.
	type aliasIdx struct {
		idx    int
		prefix string
	}

	aliases := make([]aliasIdx, len(joins))
	for i, j := range joins {
		aliases[i] = aliasIdx{idx: i, prefix: j.alias + "__"}
	}

	// Insertion sort: small N (number of joins is usually <10), no need
	// for sort.Slice overhead.
	for i := 1; i < len(aliases); i++ {
		for k := i; k > 0 && len(aliases[k].prefix) > len(aliases[k-1].prefix); k-- {
			aliases[k], aliases[k-1] = aliases[k-1], aliases[k]
		}
	}

	for ci, col := range columns {
		owned := false

		for _, a := range aliases {
			if len(col) > len(a.prefix) && col[:len(a.prefix)] == a.prefix {
				suffix := col[len(a.prefix):]
				plan.joinColsIdx[a.idx] = append(plan.joinColsIdx[a.idx], ci)
				plan.joinCols[a.idx] = append(plan.joinCols[a.idx], suffix)
				owned = true

				break
			}
		}

		if !owned {
			plan.mainColsIdx = append(plan.mainColsIdx, ci)
			plan.mainCols = append(plan.mainCols, col)
		}
	}

	return plan, nil
}

// scanOneRow allocates per-owner value slots, runs rows.Scan, then walks
// the plan: assigns main columns, builds each child, and mounts children
// onto their parents. Returns the populated main instance.
func scanOneRow[T any](rows *sql.Rows, plan *scanPlan, columns []string, scanner *Scanner[T]) (*T, error) {
	values := make([]any, len(columns))

	mainVals, err := scanner.ScanValues(plan.mainCols)
	if err != nil {
		return nil, err
	}

	for i, ci := range plan.mainColsIdx {
		values[ci] = mainVals[i]
	}

	// Allocate per-join child + values.
	type joinScratch struct {
		child   any
		vals    []any
	}

	scratch := make([]joinScratch, len(plan.joins))
	for ji, j := range plan.joins {
		child, vals, err := j.alloc(plan.joinCols[ji])
		if err != nil {
			return nil, err
		}

		scratch[ji] = joinScratch{child: child, vals: vals}
		for i, ci := range plan.joinColsIdx[ji] {
			values[ci] = vals[i]
		}
	}

	if err := rows.Scan(values...); err != nil {
		return nil, err
	}

	// Materialize the main row.
	var m T
	if err := scanner.AssignValues(&m, plan.mainCols, mainVals); err != nil {
		return nil, err
	}

	// Materialize each child + detect all-null.
	allNullByAlias := make(map[string]bool, len(plan.joins))
	childByAlias := make(map[string]any, len(plan.joins))

	for ji, j := range plan.joins {
		allNull, err := j.assign(scratch[ji].child, plan.joinCols[ji], scratch[ji].vals)
		if err != nil {
			return nil, err
		}

		allNullByAlias[j.alias] = allNull
		childByAlias[j.alias] = scratch[ji].child
	}

	// Mount children onto their parents. Process deepest first so a
	// nested child is already populated when its parent gets mounted.
	// Plan join order is registration order; we sort by alias depth here.
	order := make([]int, len(plan.joins))
	for i := range order {
		order[i] = i
	}

	for i := 1; i < len(order); i++ {
		for k := i; k > 0 && joinDepth(plan.joins[order[k]].alias) > joinDepth(plan.joins[order[k-1]].alias); k-- {
			order[k], order[k-1] = order[k-1], order[k]
		}
	}

	for _, oi := range order {
		j := plan.joins[oi]

		var parent any
		if j.parentAlias == "" {
			parent = &m
		} else {
			parent = childByAlias[j.parentAlias]
		}

		// If the parent itself is all-null (LEFT JOIN miss), skip mounting
		// — the nested child has nothing to attach to.
		if j.parentAlias != "" && allNullByAlias[j.parentAlias] {
			continue
		}

		if err := j.mount(parent, childByAlias[j.alias], allNullByAlias[j.alias]); err != nil {
			return nil, err
		}
	}

	return &m, nil
}

// joinDepth counts alias segments separated by "__" (e.g. "Seller__Company"
// is depth 2). Used to mount nested children before their parents.
func joinDepth(alias string) int {
	depth := 1
	for i := 0; i+1 < len(alias); i++ {
		if alias[i] == '_' && alias[i+1] == '_' {
			depth++
			i++
		}
	}

	return depth
}

