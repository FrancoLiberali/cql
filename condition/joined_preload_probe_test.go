package condition_test

// Probes the exact SQL shape, column aliases, and gorm semantics of CQL's
// joined preloads (JoinCondition.Preload). The output of these tests
// drives the runtime + generator design for the joined-preload fast path.
//
// What we need to learn before designing the runtime:
//   1. Column-alias format for direct and nested preloads
//   2. Order of columns in rows.Columns()
//   3. Behavior of LEFT JOIN when no child row matches (all-NULL handling)
//   4. Whether gorm distinguishes pointer vs value relations when nilling
//   5. How repeated joins of the same model are aliased (Appearance)

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type probeBrand struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func (probeBrand) TableName() string { return "brands" }

type probePhone struct {
	ID      uint `gorm:"primarykey"`
	Name    string
	BrandID uint
	Brand   probeBrand
}

func (probePhone) TableName() string { return "phones" }

func openProbe(t *testing.T) *gorm.DB {
	t.Helper()

	rec := &recordingLogger{Interface: logger.Default.LogMode(logger.Silent)}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: rec})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&probeBrand{}, &probePhone{}))

	return db
}

// recordingLogger captures the last SQL gorm executed so the probe can
// inspect it without parsing test output.
type recordingLogger struct {
	logger.Interface

	last string
}

func (r *recordingLogger) Trace(_ context.Context, _ time.Time, fc func() (string, int64), _ error) {
	if fc != nil {
		sql, _ := fc()
		r.last = sql
	}
}

// TestProbe_DirectPreloadColumns documents how a single-level joined preload
// is expressed by gorm — what the SELECT clause looks like, what columns
// rows.Columns() returns, and the order they appear in. Mirrors what CQL
// produces today via JoinCondition.Preload + AddSelectField (alias format
// `"<TableAlias>__<column>"`, e.g. `"Brand__id"`).
func TestProbe_DirectPreloadColumns(t *testing.T) {
	db := openProbe(t)

	brand := &probeBrand{Name: "acme"}
	require.NoError(t, db.Create(brand).Error)
	require.NoError(t, db.Create(&probePhone{Name: "p1", BrandID: brand.ID}).Error)

	// This mirrors what CQL's joined preload produces: an explicit JOIN +
	// AS-aliased selects of the joined columns. We construct it manually
	// here so we can lock in the exact format.
	rows, err := db.Raw(`
		SELECT phones.*,
		       Brand.id   AS "Brand__id",
		       Brand.name AS "Brand__name"
		  FROM phones
		  LEFT JOIN brands AS Brand ON Brand.id = phones.brand_id
	`).Rows()
	require.NoError(t, err)

	defer rows.Close()

	cols, err := rows.Columns()
	require.NoError(t, err)

	t.Logf("[ALIAS SPEC] cols=%v", cols)

	for i, c := range cols {
		t.Logf("  col[%d] = %q", i, c)
	}
}

// TestProbe_LeftJoinNoMatch documents what happens when LEFT JOIN finds no
// matching child row — every joined column comes back NULL. This is the
// signal the scanner uses to skip mounting the child onto the parent.
func TestProbe_LeftJoinNoMatch(t *testing.T) {
	db := openProbe(t)

	// Phone with no Brand row matching brand_id (use BrandID = 0 so no
	// brand satisfies the join condition).
	require.NoError(t, db.Create(&probePhone{Name: "p_no_brand", BrandID: 999}).Error)

	rows, err := db.Raw(`
		SELECT phones.*,
		       Brand.id   AS "Brand__id",
		       Brand.name AS "Brand__name"
		  FROM phones
		  LEFT JOIN brands AS Brand ON Brand.id = phones.brand_id
	`).Rows()
	require.NoError(t, err)

	defer rows.Close()

	for rows.Next() {
		var (
			pid      uint
			pname    string
			pbrand   uint
			bidVal   any
			bnameVal any
		)

		require.NoError(t, rows.Scan(&pid, &pname, &pbrand, &bidVal, &bnameVal))
		t.Logf("[LEFT JOIN NO MATCH] phone=%d/%s brand_id=%d  joined: Brand.id=%v Brand.name=%v",
			pid, pname, pbrand, bidVal, bnameVal)

		if bidVal != nil || bnameVal != nil {
			t.Errorf("expected both joined columns to be NULL when no match; got id=%v name=%v", bidVal, bnameVal)
		}
	}
}

// TestProbe_NestedPreloadColumns documents the alias chain for 2-level
// preloads (Sale.Seller.Company-style). CQL's Table.DeliverTable appends
// "__<relation>" to the parent alias, so nested columns become
// `"Seller__Company__id"`.
type probeCompany struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func (probeCompany) TableName() string { return "companies" }

type probeSeller struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	CompanyID *uint
	Company   *probeCompany
}

func (probeSeller) TableName() string { return "sellers" }

type probeSale struct {
	ID       uint `gorm:"primarykey"`
	Amount   int
	SellerID *uint
	Seller   *probeSeller
}

func (probeSale) TableName() string { return "sales" }

func TestProbe_NestedPreloadColumns(t *testing.T) {
	db := openProbe(t)
	require.NoError(t, db.AutoMigrate(&probeCompany{}, &probeSeller{}, &probeSale{}))

	co := &probeCompany{Name: "ditrit"}
	require.NoError(t, db.Create(co).Error)

	s := &probeSeller{Name: "alice", CompanyID: &co.ID}
	require.NoError(t, db.Create(s).Error)

	require.NoError(t, db.Create(&probeSale{Amount: 100, SellerID: &s.ID}).Error)

	rows, err := db.Raw(`
		SELECT sales.*,
		       Seller.id            AS "Seller__id",
		       Seller.name          AS "Seller__name",
		       Seller.company_id    AS "Seller__company_id",
		       Seller__Company.id   AS "Seller__Company__id",
		       Seller__Company.name AS "Seller__Company__name"
		  FROM sales
		  LEFT JOIN sellers   AS Seller          ON Seller.id          = sales.seller_id
		  LEFT JOIN companies AS Seller__Company ON Seller__Company.id = Seller.company_id
	`).Rows()
	require.NoError(t, err)

	defer rows.Close()

	cols, err := rows.Columns()
	require.NoError(t, err)

	t.Logf("[NESTED ALIAS SPEC] cols=%v", cols)

	for i, c := range cols {
		t.Logf("  col[%d] = %q", i, c)
	}

	// Verify the prefix matcher we'd build at runtime works on these:
	for _, c := range cols {
		switch {
		case len(c) > len("Seller__Company__") && c[:len("Seller__Company__")] == "Seller__Company__":
			t.Logf("  matched NESTED prefix on %q", c)
		case len(c) > len("Seller__") && c[:len("Seller__")] == "Seller__":
			t.Logf("  matched DIRECT prefix on %q", c)
		default:
			t.Logf("  no prefix (main): %q", c)
		}
	}
}

// TestProbe_ValueRelationNoMatch confirms that when LEFT JOIN finds no
// match, gorm leaves a value-type relation as the zero struct. The fast
// path's Mount(allNull=true) is a no-op for value relations to match.
type probeBrandValue struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func (probeBrandValue) TableName() string { return "brands_value" }

type probePhoneValueRelation struct {
	ID      uint `gorm:"primarykey"`
	Name    string
	BrandID uint
	Brand   probeBrandValue // value, not pointer
}

func (probePhoneValueRelation) TableName() string { return "phones_value" }

func TestProbe_ValueRelationNoMatch(t *testing.T) {
	db := openProbe(t)
	require.NoError(t, db.AutoMigrate(&probeBrandValue{}, &probePhoneValueRelation{}))

	require.NoError(t, db.Create(&probePhoneValueRelation{Name: "no_brand", BrandID: 999}).Error)

	var phone probePhoneValueRelation
	require.NoError(t, db.Preload("Brand").First(&phone).Error)

	t.Logf("[VALUE RELATION, NO MATCH] phone.Brand=%+v", phone.Brand)

	if phone.Brand.ID != 0 || phone.Brand.Name != "" {
		t.Fatalf("expected zero-valued Brand, got %+v", phone.Brand)
	}
}
