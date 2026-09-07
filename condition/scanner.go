package condition

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/FrancoLiberali/cql/model"
)

var (
	errJoinScanType  = errors.New("cql join scan: unexpected scan type")
	errJoinMountType = errors.New("cql join mount: unexpected mount type")
	errHasManyType   = errors.New("cql hasmany: unexpected parent type")
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
	// ReleaseValues returns each scanned cell to its sync.Pool. Optional;
	// when nil the runtime just drops the values for GC. Generated code
	// emits this alongside ScanValues so cells acquired from the per-type
	// pools (see scanner_pool.go) are returned after AssignValues has
	// copied the data into dest. Mirrors gorm's per-Field NewValuePool
	// usage in scan.go.
	ReleaseValues func(columns []string, values []any)
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
	// release returns scanned cells to their per-type pools.
	release func(cols []string, vals []any)
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
				return false, fmt.Errorf("%w: expected *%T", errJoinScanType, *new(Child))
			}

			if err := rs.ChildScanner.AssignValues(c, cols, vals); err != nil {
				return false, err
			}

			return allValuesNull(vals), nil
		},
		mount: func(parent any, child any, allNull bool) error {
			p, ok := parent.(*Parent)
			if !ok {
				return fmt.Errorf("%w: expected *%T parent, got %T", errJoinMountType, *new(Parent), parent)
			}

			if allNull {
				if rs.SetNil != nil {
					rs.SetNil(p)
				}

				return nil
			}

			c, ok := child.(*Child)
			if !ok {
				return fmt.Errorf("%w: expected *%T child, got %T", errJoinMountType, *new(Child), child)
			}

			rs.Mount(p, c)

			return nil
		},
		release: func(cols []string, vals []any) {
			if rs.ChildScanner != nil && rs.ChildScanner.ReleaseValues != nil {
				rs.ChildScanner.ReleaseValues(cols, vals)
			}
		},
	})
}

// allValuesNull checks whether every scanned cell carries the NULL sentinel
// for its concrete sql.Null* type. Returns true only when every cell either
// is a non-Valid sql.Null* or is nil. Used to detect LEFT-JOIN-no-match
// rows so we don't mount a zero-valued child.
func allValuesNull(values []any) bool {
	return !slices.ContainsFunc(values, valueHasData)
}

// valueHasData reports whether a single scanned cell holds a real (non-NULL)
// value. allValuesNull uses it to detect LEFT-JOIN-no-match rows (every cell
// NULL) so they aren't mounted as zero-valued children.
func valueHasData(v any) bool {
	switch t := v.(type) {
	case *sql.NullBool:
		return t.Valid
	case *sql.NullString:
		return t.Valid
	case *sql.NullInt64:
		return t.Valid
	case *sql.NullInt32:
		return t.Valid
	case *sql.NullInt16:
		return t.Valid
	case *sql.NullByte:
		return t.Valid
	case *sql.NullFloat64:
		return t.Valid
	case *sql.NullTime:
		return t.Valid
	case *gorm.DeletedAt:
		return t.Valid
	case *NullSink:
		// always the "null" sentinel
		return false
	case *NullableScanner:
		// custom scanners; we have no .Valid to inspect, so be
		// conservative: treat as having data.
		return true
	case *model.UUID:
		// CQL convention: NilUUID is the sentinel for "no value". Generated
		// scanners scan a UUID column directly into model.UUID; a
		// LEFT-JOIN-no-match leaves it at NilUUID.
		return *t != model.NilUUID
	default:
		// unknown type (custom scanner direct, gorm types, etc.) — be
		// conservative: assume non-Null-wrapped values had data.
		return true
	}
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
					return fmt.Errorf("%w: %q expected *%T parent, got %T",
						errHasManyType, loader.CollectionField, *new(Parent), p)
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
				pp, _ := p.(*Parent)
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

// prepopulateSelectAndFromForFind stuffs the SELECT + FROM clauses into
// gorm's Clauses map ahead of BuildQuerySQL so the callback's
// `make([]clause.Column, ...)` + `AddClauseIfNotExists(clause.From{}) /
// AddClauseIfNotExists(clauseSelect)` trio become no-ops (map lookup +
// present-so-skip).
//
// FROM is always safe on the Find/Take/First code path: only Rows-style
// query callbacks reach here, so UPDATE/DELETE's alternate SQL builders
// (which would trip over a leftover FROM clause) are never involved. When
// gorm needs to inject JOINs into FROM, GenJoinClauses' caller overrides
// via AddClause anyway.
//
// SELECT is safe only when CQL isn't tracking joined selects. For joined
// preloads, CQL appends `alias.column AS "alias__column"` fragments onto
// Statement.Selects via AddSelectField. If we pre-populate SELECT before
// BuildQuerySQL runs, its `AddClauseIfNotExists(clauseSelect)` no-ops and
// those joined columns never make it into the rendered SQL. In that case
// we let gorm handle SELECT and keep just the FROM win.
func prepopulateSelectAndFromForFind(q *CQLQuery) {
	stmt := q.gormDB.Statement

	if _, has := stmt.Clauses["FROM"]; !has {
		fc := clause.Clause{Name: "FROM"}
		clause.From{}.MergeClause(&fc)
		stmt.Clauses["FROM"] = fc
	}

	// Extra selects beyond the base "<table>.*" (len > 1) mean gorm must
	// build the SELECT itself: e.g. postgres ORDER BY adds a
	// "<table>.<col> AS \"<table>__<col>\"" alias to Statement.Selects on
	// the initial table (so hasJoinedSelects stays false), and clobbering
	// it here would drop that column and break the ORDER BY.
	if q.hasJoinedSelects || len(q.activeJoins) > 0 || len(stmt.Selects) > 1 {
		return
	}

	if _, has := stmt.Clauses["SELECT"]; !has {
		sc := clause.Clause{Name: "SELECT"}
		clause.Select{Columns: []clause.Column{{Name: q.initialTable.Name + ".*", Raw: true}}}.MergeClause(&sc)
		stmt.Clauses["SELECT"] = sc
		// Nil out Selects so BuildQuerySQL doesn't re-enter its
		// per-column loop and re-allocate a clause.Column slice that
		// AddClauseIfNotExists then drops anyway.
		stmt.Selects = nil
	}
}

// findWith materializes every matching row into *dest using scanner. Caller
// must have validated that the fast path applies (scanner != nil and
// canUseFastScan(q) == true). After the main scan, runs registered
// HasMany loaders (one extra SELECT per relation) and mounts the grouped
// children onto each parent.
func findWith[T any](q *CQLQuery, dest *[]*T, scanner *Scanner[T]) error {
	q.flushPending()

	// HasMany loaders own the child fetch; if gorm Preloads are also
	// registered (kept alive for non-Find paths) strip them here so gorm
	// doesn't duplicate the work + override the loader's mount.
	if len(q.activeHasMany) > 0 {
		stripGormPreloads(q)
	}

	// Skip BuildQuerySQL's SELECT/FROM merge dance when there are no
	// JOINs (neither user-added nor joined-preload). Pre-populating the
	// two clauses here is safe: with no JOINs, GenJoinClauses never runs
	// and never needs to append columns to the SELECT that would be lost
	// to AddClauseIfNotExists' skip-when-present rule. Update/Delete
	// paths never enter this function, so their SQL builders aren't
	// affected. Saves ~5 allocs per query on the flat path.
	prepopulateSelectAndFromForFind(q)

	rows, err := q.gormDB.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	plan := buildScanPlan(columns, q.activeJoins)

	// Allocate scratch once per query; reused across rows. The Acquire/
	// Release cell pool refills the inner slots per row, so the outer
	// slices never carry stale pool-held pointers across iterations.
	scratch := newRowScratch(len(columns), len(plan.joins))

	for rows.Next() {
		m, err := scanOneRow[T](rows, plan, scanner, scratch)
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

// rowScratch holds per-query (not per-row) work buffers reused across the
// findWith hot loop. Saves nRows × 4 slice allocations vs the previous
// "allocate-per-row" approach.
type rowScratch struct {
	values       []any
	joinChildren []any
	joinVals     [][]any
	joinAllNull  []bool
}

func newRowScratch(totalCols, nJoins int) *rowScratch {
	s := &rowScratch{
		values: make([]any, totalCols),
	}
	if nJoins > 0 {
		s.joinChildren = make([]any, nJoins)
		s.joinVals = make([][]any, nJoins)
		s.joinAllNull = make([]bool, nJoins)
	}

	return s
}

// Precomputed clause.Clause values reused across queries to skip the
// alloc-heavy DB.Order / DB.Limit paths — each of which allocates via
// `getInstance + type-switch/variadic + AddClause` on every call. The
// clause.Clause struct is copied into Statement.Clauses by value, and
// each Expression's MergeClause reads from — never mutates — the
// stored value, so sharing these package-level clauses across queries
// is safe. Name is empty on the LIMIT variant because Limit.Build
// already writes the "LIMIT " prefix itself (see clause.Clause.Build).
var (
	firstOrderByClause = clause.Clause{
		Name: "ORDER BY",
		Expression: clause.OrderBy{
			Columns: []clause.OrderByColumn{{Column: clause.PrimaryColumn}},
		},
	}
	lastOrderByClause = clause.Clause{
		Name: "ORDER BY",
		Expression: clause.OrderBy{
			Columns: []clause.OrderByColumn{{Column: clause.PrimaryColumn, Desc: true}},
		},
	}
	limitOneValue  = 1
	limitOneClause = clause.Clause{
		Expression: clause.Limit{Limit: &limitOneValue},
	}
)

// firstWith materializes the first row (primary-key ordered) into dest.
// Returns gorm.ErrRecordNotFound when no row matches, matching the gorm
// fallback path's contract.
func firstWith[T any](q *CQLQuery, dest **T, scanner *Scanner[T]) error {
	q.gormDB.Statement.Clauses["ORDER BY"] = firstOrderByClause
	return scanOne(q, q.gormDB, dest, scanner)
}

// takeWith materializes any one matching row (no ordering) into dest.
func takeWith[T any](q *CQLQuery, dest **T, scanner *Scanner[T]) error {
	return scanOne(q, q.gormDB, dest, scanner)
}

// lastWith materializes the last row (primary-key DESC) into dest.
func lastWith[T any](q *CQLQuery, dest **T, scanner *Scanner[T]) error {
	q.gormDB.Statement.Clauses["ORDER BY"] = lastOrderByClause
	return scanOne(q, q.gormDB, dest, scanner)
}

// scanOne materializes a single row + runs registered HasMany loaders
// against the singleton parent.
func scanOne[T any](q *CQLQuery, db *gorm.DB, dest **T, scanner *Scanner[T]) error {
	q.flushPending()

	if len(q.activeHasMany) > 0 {
		stripGormPreloads(q)
	}

	prepopulateSelectAndFromForFind(q)

	// Directly install a shared LIMIT 1 clause instead of db.Limit(1)
	// (getInstance + type-switch + AddClause allocs).
	db.Statement.Clauses["LIMIT"] = limitOneClause

	rows, err := db.Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		if rowsErr := rows.Err(); rowsErr != nil {
			return rowsErr
		}

		return gorm.ErrRecordNotFound
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	plan := buildScanPlan(columns, q.activeJoins)

	scratch := newRowScratch(len(columns), len(plan.joins))

	m, err := scanOneRow[T](rows, plan, scanner, scratch)
	if err != nil {
		return err
	}

	*dest = m

	if err := rows.Err(); err != nil {
		return err
	}

	return runHasManyLoaders[T](q, []*T{m})
}

// scanPlan precomputes per-column routing AND per-join mount sequencing once
// (instead of doing prefix matching + alias-map lookups per row × per
// column). Built once per query in buildScanPlan; consumed by scanOneRow
// for every row.
type scanPlan struct {
	mainColsIdx []int         // result indices owned by the main scanner
	mainCols    []string      // names of main columns, parallel to mainColsIdx
	joins       []*activeJoin // registration order (top-level → nested)
	joinColsIdx [][]int       // per-join result indices
	joinCols    [][]string    // per-join column names
	// parentIdx is the join index of each join's parent, or -1 for top-
	// level joins (mount onto the main row). Avoids the per-row
	// childByAlias map lookup.
	parentIdx []int
	// mountOrder is the join indices sorted by alias depth descending so
	// nested children mount onto their parents before the parents mount
	// onto theirs. Precomputed once instead of re-sorting per row.
	mountOrder []int
}

// buildScanPlan inspects the row's columns once and produces a routing plan.
// Joined columns must be assigned to the deepest matching alias (so
// "Seller__Company__id" goes to the nested join, not the parent), which is
// achieved by sorting alias prefixes longest-first when matching.
func buildScanPlan(columns []string, joins []*activeJoin) *scanPlan {
	plan := &scanPlan{
		joins:       joins,
		joinColsIdx: make([][]int, len(joins)),
		joinCols:    make([][]string, len(joins)),
		// Pre-size main-column slices: they can never exceed len(columns)
		// (join-owned columns are removed from the main pool). Skipping
		// the growth reallocations here saves ~3 allocs per query on the
		// common no-join Read path — the biggest single contributor to
		// CQL's DSL-setup alloc overhead on the profile.
		mainColsIdx: make([]int, 0, len(columns)),
		mainCols:    make([]string, 0, len(columns)),
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

	// Precompute parentIdx + mountOrder once. Both used per row in
	// scanOneRow to avoid map allocations + sort work in the hot path.
	aliasToIdx := make(map[string]int, len(joins))
	for i, j := range joins {
		aliasToIdx[j.alias] = i
	}

	plan.parentIdx = make([]int, len(joins))
	for i, j := range joins {
		if j.parentAlias == "" {
			plan.parentIdx[i] = -1
		} else if pi, ok := aliasToIdx[j.parentAlias]; ok {
			plan.parentIdx[i] = pi
		} else {
			plan.parentIdx[i] = -1 // parent not registered; mount onto main
		}
	}

	plan.mountOrder = make([]int, len(joins))
	for i := range plan.mountOrder {
		plan.mountOrder[i] = i
	}

	for i := 1; i < len(plan.mountOrder); i++ {
		for k := i; k > 0 && joinDepth(joins[plan.mountOrder[k]].alias) > joinDepth(joins[plan.mountOrder[k-1]].alias); k-- {
			plan.mountOrder[k], plan.mountOrder[k-1] = plan.mountOrder[k-1], plan.mountOrder[k]
		}
	}

	return plan
}

// scanOneRow allocates per-row scan destinations into the caller-provided
// scratch (reused across rows by findWith), runs rows.Scan, then walks
// the precomputed plan to assign + mount. Per-row hot path stays
// allocation-light: the only allocs are pool-Acquired cell wrappers (via
// ScanValues), the child structs, and the result *T.
func scanOneRow[T any](rows *sql.Rows, plan *scanPlan, scanner *Scanner[T], scratch *rowScratch) (*T, error) {
	values := scratch.values

	mainVals, err := scanner.ScanValues(plan.mainCols)
	if err != nil {
		return nil, err
	}

	for i, ci := range plan.mainColsIdx {
		values[ci] = mainVals[i]
	}

	// Per-join children + their scanned values. Reused from scratch — no
	// per-row alloc of the outer slices.
	joinChildren := scratch.joinChildren
	joinVals := scratch.joinVals
	joinAllNull := scratch.joinAllNull

	for ji, j := range plan.joins {
		child, vals, err := j.alloc(plan.joinCols[ji])
		if err != nil {
			return nil, err
		}

		joinChildren[ji] = child
		joinVals[ji] = vals
		joinAllNull[ji] = false

		for i, ci := range plan.joinColsIdx[ji] {
			values[ci] = vals[i]
		}
	}

	if err := rows.Scan(values...); err != nil {
		return nil, err
	}

	var m T
	if err := scanner.AssignValues(&m, plan.mainCols, mainVals); err != nil {
		return nil, err
	}

	if scanner.ReleaseValues != nil {
		scanner.ReleaseValues(plan.mainCols, mainVals)
	}

	for ji, j := range plan.joins {
		allNull, err := j.assign(joinChildren[ji], plan.joinCols[ji], joinVals[ji])
		if err != nil {
			return nil, err
		}

		joinAllNull[ji] = allNull
		// Return joined cells to their pools as soon as the child is
		// materialized — mirrors gorm's per-row Put at scan.go:111.
		if j.release != nil {
			j.release(plan.joinCols[ji], joinVals[ji])
		}
	}

	if err := mountJoins(&m, plan, joinChildren, joinAllNull); err != nil {
		return nil, err
	}

	return &m, nil
}

// mountJoins mounts each scanned join child onto its parent — the main row
// (parentIdx == -1) or a shallower join child — in the plan's precomputed
// depth-descending order, skipping any child whose parent join was a
// LEFT-JOIN-no-match (no instance to mount onto).
func mountJoins[T any](m *T, plan *scanPlan, joinChildren []any, joinAllNull []bool) error {
	for _, ji := range plan.mountOrder {
		j := plan.joins[ji]
		parentJoinIdx := plan.parentIdx[ji]

		if parentJoinIdx >= 0 && joinAllNull[parentJoinIdx] {
			continue
		}

		var parent any
		if parentJoinIdx < 0 {
			parent = m
		} else {
			parent = joinChildren[parentJoinIdx]
		}

		if err := j.mount(parent, joinChildren[ji], joinAllNull[ji]); err != nil {
			return err
		}
	}

	return nil
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
