package cql

import (
	"context"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

// Benchmarks comparing CQL's fast-scan path (Query[T]().Find()) against
// **plain gorm** (db.Find(&dest)) on identical schema + data. This is the
// comparison end-users actually face: "should I use CQL or just call gorm
// directly?"
//
// Same sqlite :memory:, same seeded rows, same query semantics — the only
// difference is which API drives the scan. `_Cql` variants go through CQL
// + the code-generated Scanner; `_Gorm` variants call gorm.io/gorm methods
// directly.
//
// Run:
//   go test -run x -bench BenchmarkScan -benchmem -benchtime=2s ./...

func benchOpenDB(b *testing.B) (*DB, *gorm.DB) {
	b.Helper()

	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		b.Fatalf("open sqlite: %v", err)
	}

	if err := gdb.AutoMigrate(
		&models.Brand{},
		&models.Company{},
		&models.University{},
		&models.Seller{},
		&models.Phone{},
	); err != nil {
		b.Fatalf("automigrate: %v", err)
	}

	db, err := Open(sqlite.Open(":memory:"))
	if err != nil {
		b.Fatalf("cql open: %v", err)
	}
	// Share the migrated gorm handle with both APIs so benchmarks hit the
	// same schema + data without duplicating setup.
	db.GormDB = gdb

	return db, gdb
}

func seedBrands(b *testing.B, gdb *gorm.DB, n int) {
	b.Helper()

	for i := 0; i < n; i++ {
		if err := gdb.Create(&models.Brand{Name: fmt.Sprintf("brand_%d", i)}).Error; err != nil {
			b.Fatalf("seed brand: %v", err)
		}
	}
}

func seedPhones(b *testing.B, gdb *gorm.DB, n int) {
	b.Helper()

	for i := 0; i < n; i++ {
		brand := &models.Brand{Name: fmt.Sprintf("b_%d", i)}
		if err := gdb.Create(brand).Error; err != nil {
			b.Fatalf("seed brand: %v", err)
		}

		if err := gdb.Create(&models.Phone{Name: fmt.Sprintf("p_%d", i), BrandID: uint(brand.ID)}).Error; err != nil {
			b.Fatalf("seed phone: %v", err)
		}
	}
}

func seedCompaniesWithSellers(b *testing.B, gdb *gorm.DB, companies, sellersPerCompany int) {
	b.Helper()

	for i := 0; i < companies; i++ {
		c := &models.Company{Name: fmt.Sprintf("c_%d", i)}
		if err := gdb.Create(c).Error; err != nil {
			b.Fatalf("seed company: %v", err)
		}

		for j := 0; j < sellersPerCompany; j++ {
			cid := c.ID
			if err := gdb.Create(&models.Seller{Name: fmt.Sprintf("s_%d_%d", i, j), CompanyID: &cid}).Error; err != nil {
				b.Fatalf("seed seller: %v", err)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Find: return all rows.
// ---------------------------------------------------------------------------

func benchFindBrandsCql(b *testing.B, n int) {
	db, gdb := benchOpenDB(b)
	seedBrands(b, gdb, n)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		out, err := Query[models.Brand](ctx, db,
			conditions.Brand.Name.IsUnsafe().NotEq(condition.String("___nope")),
		).Find()
		if err != nil {
			b.Fatal(err)
		}
		if len(out) != n {
			b.Fatalf("expected %d rows, got %d", n, len(out))
		}
	}
}

func benchFindBrandsGorm(b *testing.B, n int, gdb *gorm.DB) {
	if gdb == nil {
		_, gdb = benchOpenDB(b)
		seedBrands(b, gdb, n)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var out []*models.Brand
		if err := gdb.Where("name != ?", "___nope").Find(&out).Error; err != nil {
			b.Fatal(err)
		}
		if len(out) != n {
			b.Fatalf("expected %d rows, got %d", n, len(out))
		}
	}
}

func BenchmarkScan_Find_1Row_Cql(b *testing.B)     { benchFindBrandsCql(b, 1) }
func BenchmarkScan_Find_1Row_Gorm(b *testing.B)    { benchFindBrandsGorm(b, 1, nil) }
func BenchmarkScan_Find_100Rows_Cql(b *testing.B)  { benchFindBrandsCql(b, 100) }
func BenchmarkScan_Find_100Rows_Gorm(b *testing.B) { benchFindBrandsGorm(b, 100, nil) }
func BenchmarkScan_Find_1KRows_Cql(b *testing.B)   { benchFindBrandsCql(b, 1000) }
func BenchmarkScan_Find_1KRows_Gorm(b *testing.B)  { benchFindBrandsGorm(b, 1000, nil) }

// ---------------------------------------------------------------------------
// First: single row.
// ---------------------------------------------------------------------------

func BenchmarkScan_First_Cql(b *testing.B) {
	db, gdb := benchOpenDB(b)
	seedBrands(b, gdb, 50)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if _, err := Query[models.Brand](ctx, db,
			conditions.Brand.Name.IsUnsafe().NotEq(condition.String("___nope")),
		).First(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkScan_First_Gorm(b *testing.B) {
	_, gdb := benchOpenDB(b)
	seedBrands(b, gdb, 50)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var out models.Brand
		if err := gdb.Where("name != ?", "___nope").First(&out).Error; err != nil {
			b.Fatal(err)
		}
	}
}

// ---------------------------------------------------------------------------
// Joined preload (BelongsTo). Phone with Brand.
// CQL: conditions.Phone.Brand().Preload(). Gorm: db.Preload("Brand").
// ---------------------------------------------------------------------------

func benchJoinedPreloadCql(b *testing.B, n int) {
	db, gdb := benchOpenDB(b)
	seedPhones(b, gdb, n)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		out, err := Query[models.Phone](ctx, db,
			conditions.Phone.Brand().Preload(),
		).Find()
		if err != nil {
			b.Fatal(err)
		}
		if len(out) != n {
			b.Fatalf("expected %d rows, got %d", n, len(out))
		}
	}
}

func benchJoinedPreloadGorm(b *testing.B, n int) {
	_, gdb := benchOpenDB(b)
	seedPhones(b, gdb, n)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var out []*models.Phone
		if err := gdb.Preload("Brand").Find(&out).Error; err != nil {
			b.Fatal(err)
		}
		if len(out) != n {
			b.Fatalf("expected %d rows, got %d", n, len(out))
		}
	}
}

func BenchmarkScan_JoinedPreload_100_Cql(b *testing.B)  { benchJoinedPreloadCql(b, 100) }
func BenchmarkScan_JoinedPreload_100_Gorm(b *testing.B) { benchJoinedPreloadGorm(b, 100) }
func BenchmarkScan_JoinedPreload_1K_Cql(b *testing.B)   { benchJoinedPreloadCql(b, 1000) }
func BenchmarkScan_JoinedPreload_1K_Gorm(b *testing.B)  { benchJoinedPreloadGorm(b, 1000) }

// Phone-only baseline (no relation loading) so the joined-preload
// number can be split into "main scan cost" + "join composition cost".
func BenchmarkScan_Phone_NoPreload_1K_Cql(b *testing.B) {
	db, gdb := benchOpenDB(b)
	seedPhones(b, gdb, 1000)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		out, err := Query[models.Phone](ctx, db,
			conditions.Phone.Name.IsUnsafe().NotEq(condition.String("___nope")),
		).Find()
		if err != nil {
			b.Fatal(err)
		}
		if len(out) != 1000 {
			b.Fatalf("expected 1000 rows, got %d", len(out))
		}
	}
}

func BenchmarkScan_Phone_NoPreload_1K_Gorm(b *testing.B) {
	_, gdb := benchOpenDB(b)
	seedPhones(b, gdb, 1000)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var out []*models.Phone
		if err := gdb.Where("name != ?", "___nope").Find(&out).Error; err != nil {
			b.Fatal(err)
		}
		if len(out) != 1000 {
			b.Fatalf("expected 1000 rows, got %d", len(out))
		}
	}
}

// ---------------------------------------------------------------------------
// HasMany preload. Company.Sellers.
// ---------------------------------------------------------------------------

func BenchmarkScan_HasMany_50x4_Cql(b *testing.B) {
	db, gdb := benchOpenDB(b)
	seedCompaniesWithSellers(b, gdb, 50, 4)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		out, err := Query[models.Company](ctx, db,
			conditions.Company.Sellers.Preload(),
		).Find()
		if err != nil {
			b.Fatal(err)
		}
		if len(out) != 50 {
			b.Fatalf("expected 50 companies, got %d", len(out))
		}
	}
}

func BenchmarkScan_HasMany_50x4_Gorm(b *testing.B) {
	_, gdb := benchOpenDB(b)
	seedCompaniesWithSellers(b, gdb, 50, 4)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var out []*models.Company
		if err := gdb.Preload("Sellers").Find(&out).Error; err != nil {
			b.Fatal(err)
		}
		if len(out) != 50 {
			b.Fatalf("expected 50 companies, got %d", len(out))
		}
	}
}
