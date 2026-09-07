// CQL-specific benchmark scenarios that don't map onto gorm's canonical
// suite. Kept separate from benchmark_test.go (which mirrors gorm's
// tests/benchmark_test.go shapes) so that suite stays a clean 1:1
// comparison against gorm's numbers.
//
// Covered here:
//   - Multi-condition WHERE dispatch (5 conditions × 1 row / 100 rows) —
//     exercises the per-WHERE accumulate + flush path.
//   - First() — CQL's scanOne code path (distinct from Find's findWith).
//   - JoinedPreload (100 / 1K rows) — CQL's typed joined-preload scan
//     with generated RelationScanner + mount plan, the biggest win area
//     against gorm's reflective joined preload.
//   - HasMany (50 parents × 4 children) — post-scan loader path that
//     runs one child SELECT after the main query and mounts the grouped
//     children onto their parents.
//
// Shares openDB / newUser from benchmark_test.go (same package).

package tests_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/tests/conditions"
	"github.com/FrancoLiberali/cql/tests/models"
)

// -- helpers ---------------------------------------------------------

// seedUsersWithAccount inserts n users, each with an associated Account
// populated so joined-preload benches actually return joined rows.
func seedUsersWithAccount(b *testing.B, db *cql.DB, n int) {
	b.Helper()

	users := make([]*models.User, n)
	for i := range users {
		u := newUser(fmt.Sprintf("scan-%d", i))
		u.Account = models.Account{Number: fmt.Sprintf("acct-%d", i)}
		users[i] = u
	}

	if err := db.GormDB.CreateInBatches(&users, 100).Error; err != nil {
		b.Fatalf("seed: %v", err)
	}
}

// seedUsersWithPets inserts nParents users, each carrying nChildren
// Pets, so HasMany post-scan loader benches have realistic child fanout.
func seedUsersWithPets(b *testing.B, db *cql.DB, nParents, nChildren int) {
	b.Helper()

	users := make([]*models.User, nParents)
	for i := range users {
		u := newUser(fmt.Sprintf("pu-%d", i))

		u.Pets = make([]*models.Pet, nChildren)
		for j := range u.Pets {
			u.Pets[j] = &models.Pet{Name: fmt.Sprintf("pet-%d-%d", i, j)}
		}

		users[i] = u
	}

	if err := db.GormDB.CreateInBatches(&users, 50).Error; err != nil {
		b.Fatalf("seed: %v", err)
	}
}

// -- multi-condition WHERE dispatch ---------------------------------

func benchFindUsers5Where(b *testing.B, n int) {
	db := openDB(b)
	seedBenchUsers(b, db, n)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		out, err := cql.Query[models.User](ctx, db,
			conditions.User.Name.Is().NotEq(cql.String("nope1")),
			conditions.User.Name.Is().NotEq(cql.String("nope2")),
			conditions.User.Name.Is().NotEq(cql.String("nope3")),
			conditions.User.Name.Is().NotEq(cql.String("nope4")),
			conditions.User.Name.Is().NotEq(cql.String("nope5")),
		).Find()
		if err != nil {
			b.Fatal(err)
		}

		if len(out) != n {
			b.Fatalf("expected %d rows, got %d", n, len(out))
		}
	}
}

func BenchmarkFind_5Where_1Row(b *testing.B)    { benchFindUsers5Where(b, 1) }
func BenchmarkFind_5Where_100Rows(b *testing.B) { benchFindUsers5Where(b, 100) }

// -- First() (scanOne path) -----------------------------------------

// BenchmarkFirst — CQL's First() takes a different code path than Find:
// scanOne with LIMIT 1 vs findWith with row iteration. Worth measuring
// separately since some optimizations only touch one path.
func BenchmarkFirst(b *testing.B) {
	db := openDB(b)
	ctx := context.Background()

	u := newUser("bench-first")
	if _, err := cql.Insert(ctx, db, u).Exec(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, err := cql.Query[models.User](ctx, db,
			conditions.User.ID.Is().Eq(cql.UIntID(uint(u.ID))),
		).First()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// -- Joined preload (has-one) --------------------------------------

func benchJoinedPreload(b *testing.B, n int) {
	db := openDB(b)
	seedUsersWithAccount(b, db, n)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		out, err := cql.Query[models.User](ctx, db,
			conditions.User.Account().Preload(),
		).Find()
		if err != nil {
			b.Fatal(err)
		}

		if len(out) != n {
			b.Fatalf("expected %d rows, got %d", n, len(out))
		}
	}
}

func BenchmarkJoinedPreload_100(b *testing.B) { benchJoinedPreload(b, 100) }
func BenchmarkJoinedPreload_1K(b *testing.B)  { benchJoinedPreload(b, 1000) }

// -- HasMany (post-scan loader) ------------------------------------

// BenchmarkHasMany_50x4 — 50 parents × 4 pets each. Measures the
// HasMany fast loader: one child SELECT after the parent scan +
// per-parent mount. Old scan_bench_test.go used the same 50×4 shape.
func BenchmarkHasMany_50x4(b *testing.B) {
	db := openDB(b)
	seedUsersWithPets(b, db, 50, 4)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		out, err := cql.Query[models.User](ctx, db,
			conditions.User.Pets.Preload(),
		).Find()
		if err != nil {
			b.Fatal(err)
		}

		if len(out) != 50 {
			b.Fatalf("expected 50 users, got %d", len(out))
		}
	}
}
