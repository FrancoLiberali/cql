// Package tests_test hosts CQL's benchmark suite, laid out to mirror
// gorm's own gorm/tests/benchmark_test.go so numbers from both can be
// compared side by side.
//
// Only CQL is measured here. To compare against plain gorm, run gorm's
// own suite from its repo (github.com/go-gorm/gorm/tests). Extra
// scenarios that don't map onto gorm's canonical shapes belong in a
// dedicated comparison repo, not here.
//
// The row schema — models.User — is a scalar-only mirror of gorm's
// utils/tests.User: same columns (id, timestamps, name, age, birthday,
// company_id, manager_id, active), just no associations. That keeps
// per-row scan cost directly comparable between CQL's runs here and
// gorm's own BenchmarkFind/Scan/... invocations.
//
// Shape parity with gorm:
//
//	gorm/tests/benchmark_test.go       cql/benchmarks/benchmark_test.go
//	--------------------------------   --------------------------------
//	BenchmarkCreate
//	BenchmarkFind (PK, 1 row)          BenchmarkFind
//	BenchmarkScan (raw SQL, 1 row)     BenchmarkScan (uses CQL DSL — see note)
//	BenchmarkScanSlice (10K rows)      BenchmarkScanSlice
//	BenchmarkScanSlicePointer
//	BenchmarkUpdate
//	BenchmarkDelete
package benchmarks_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	gormlogger "gorm.io/gorm/logger"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/benchmarks/conditions"
	"github.com/FrancoLiberali/cql/benchmarks/models"
)

func openDB(b *testing.B) *cql.DB {
	b.Helper()

	db, err := cql.Open(sqlite.Open(":memory:"), &cql.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		b.Fatalf("open db: %v", err)
	}

	if err := db.GormDB.AutoMigrate(&models.User{}, &models.Account{}, &models.Pet{}); err != nil {
		b.Fatalf("migrate: %v", err)
	}

	return db
}

// newUser matches gorm's helper GetUser(name, Config{}): scalar
// fields only, birthday rounded to the second so encoded values are
// stable across Insert/Find round trips.
func newUser(name string) *models.User {
	birthday := time.Now().Round(time.Second)

	return &models.User{
		Name:     name,
		Age:      18,
		Birthday: &birthday,
	}
}

func seedBenchUsers(b *testing.B, db *cql.DB, n int) {
	b.Helper()

	if n == 0 {
		return
	}

	users := make([]*models.User, n)
	for i := range users {
		users[i] = newUser(fmt.Sprintf("scan-%d", i))
	}

	// sqlite's default variable limit is 999. User has enough columns
	// that inserting all 10K in one shot blows past it; chunk into
	// batches gorm can bind safely.
	if err := db.GormDB.CreateInBatches(&users, 100).Error; err != nil {
		b.Fatalf("seed: %v", err)
	}
}

// BenchmarkCreate — mirrors gorm's `db.Create(&user)` shape.
func BenchmarkCreate(b *testing.B) {
	db := openDB(b)
	ctx := context.Background()

	user := newUser("bench")

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		user.ID = 0
		if _, err := cql.Insert(ctx, db, user).Exec(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFind — mirrors gorm's `db.Find(&User{}, "id = ?", id)` shape:
// primary-key lookup returning one row.
func BenchmarkFind(b *testing.B) {
	db := openDB(b)
	ctx := context.Background()

	user := newUser("bench-find")
	if _, err := cql.Insert(ctx, db, user).Exec(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, err := cql.Query[models.User](ctx, db,
			conditions.User.ID.Is().Eq(cql.UIntID(uint(user.ID))),
		).First()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkScan — gorm's original uses `db.Raw("select * from users where
// id = ?").Scan(&u)`. CQL has no raw-SQL surface, so this ends up being a
// PK-lookup Query.First() — same call as BenchmarkFind. Kept as a separate
// entry so the row layout matches gorm's suite.
func BenchmarkScan(b *testing.B) {
	db := openDB(b)
	ctx := context.Background()

	user := newUser("bench-scan")
	if _, err := cql.Insert(ctx, db, user).Exec(); err != nil {
		b.Fatal(err)
	}

	var out *models.User

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		got, err := cql.Query[models.User](ctx, db,
			conditions.User.ID.Is().Eq(cql.UIntID(uint(user.ID))),
		).First()
		if err != nil {
			b.Fatal(err)
		}

		out = got
	}

	_ = out
}

// BenchmarkScanSlice — mirrors gorm's 10K-row `.Scan(&[]User)` shape.
func BenchmarkScanSlice(b *testing.B) {
	db := openDB(b)
	seedBenchUsers(b, db, 10_000)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		out, err := cql.Query[models.User](ctx, db).Find()
		if err != nil {
			b.Fatal(err)
		}

		if len(out) != 10_000 {
			b.Fatalf("expected 10000 rows, got %d", len(out))
		}
	}
}

// BenchmarkScanSlicePointer — gorm's shape scans into `[]*User`. CQL's
// Query[T].Find() already returns []*T, so this is the same call as
// BenchmarkScanSlice above. Present for shape parity with gorm's suite.
func BenchmarkScanSlicePointer(b *testing.B) {
	db := openDB(b)
	seedBenchUsers(b, db, 10_000)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		out, err := cql.Query[models.User](ctx, db).Find()
		if err != nil {
			b.Fatal(err)
		}

		if len(out) != 10_000 {
			b.Fatalf("expected 10000 rows, got %d", len(out))
		}
	}
}

// BenchmarkUpdate — mirrors gorm's `db.Model(&user).Updates(map[…])` shape.
func BenchmarkUpdate(b *testing.B) {
	db := openDB(b)
	ctx := context.Background()

	user := newUser("bench-update")
	if _, err := cql.Insert(ctx, db, user).Exec(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := range b.N {
		_, err := cql.Update[models.User](ctx, db,
			conditions.User.ID.Is().Eq(cql.UIntID(uint(user.ID))),
		).Set(
			conditions.User.Age.Set().Eq(cql.UInt(uint(i))),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDelete — mirrors gorm's insert+delete-per-iteration shape.
func BenchmarkDelete(b *testing.B) {
	db := openDB(b)
	ctx := context.Background()

	user := newUser("bench-delete")

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		b.StopTimer()

		user.ID = 0
		if _, err := cql.Insert(ctx, db, user).Exec(); err != nil {
			b.Fatal(err)
		}

		b.StartTimer()

		if _, err := cql.Delete[models.User](ctx, db,
			conditions.User.ID.Is().Eq(cql.UIntID(uint(user.ID))),
		).Exec(); err != nil {
			b.Fatal(err)
		}
	}
}
