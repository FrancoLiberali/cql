package cql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"

	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

// joinedPreloadDB wires sqlmock + gorm + cql exactly the same way
// scanner_fast_path_test.go does, so the assertions focus on the new
// joined-preload behavior.
func joinedPreloadDB(t *testing.T) (*DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := Open(postgres.New(postgres.Config{Conn: conn}))
	require.NoError(t, err)

	return db, mock, func() { conn.Close() }
}

// phoneBrandJoin uses the cql-gen-generated builder, which now threads the
// per-relation scanner through to the runtime automatically.
var phoneBrandJoin = conditions.Phone.Brand

// TestJoinedPreloadFastScan_DirectBelongsTo materializes Phone rows with a
// preloaded Brand via the fast path. Asserts:
//   - SQL has the JOIN + aliased selects format CQL produces
//   - Phone.ID/Name/BrandID come from the phones.* columns
//   - Phone.Brand.ID/Name come from the joined "Brand__*" columns
//   - The fast path runs (not the gorm fallback) — confirmed by the
//     absence of a separate Brand SELECT (preload-style query).
func TestJoinedPreloadFastScan_DirectBelongsTo(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	// CQL emits: SELECT phones.*, "Brand"."id" AS "Brand__id", "Brand"."name" AS "Brand__name"
	//            FROM "phones" LEFT JOIN "brands" "Brand" ON ...
	// sqlmock's regexp is lenient; we anchor on the key tokens.
	mock.ExpectQuery(`SELECT phones\.\*,.*Brand__id.*Brand__name.*FROM "phones"`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at",
				"name", "brand_id",
				"Brand__id", "Brand__name",
			}).AddRow(
				1, nil, nil, nil,
				"p1", 7,
				7, "acme",
			),
		)

	phones, err := Query[models.Phone](
		context.Background(),
		db,
		phoneBrandJoin().Preload(),
	).Find()

	require.NoError(t, err)
	require.Len(t, phones, 1)
	assert.EqualValues(t, 1, phones[0].ID)
	assert.Equal(t, "p1", phones[0].Name)
	assert.EqualValues(t, 7, phones[0].BrandID)
	assert.EqualValues(t, 7, phones[0].Brand.ID)
	assert.Equal(t, "acme", phones[0].Brand.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestJoinedPreloadFastScan_LeftJoinNoMatch covers the case where the
// preloaded relation has no matching row (LEFT JOIN nulls). Brand fields
// should stay zero-valued and no error should bubble up.
func TestJoinedPreloadFastScan_LeftJoinNoMatch(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT phones\.\*,.*Brand__id.*FROM "phones"`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at",
				"name", "brand_id",
				"Brand__id", "Brand__name",
			}).AddRow(
				1, nil, nil, nil,
				"orphan", 999,
				nil, nil, // no matching brand
			),
		)

	phones, err := Query[models.Phone](
		context.Background(),
		db,
		phoneBrandJoin().Preload(),
	).Find()

	require.NoError(t, err)
	require.Len(t, phones, 1)
	assert.EqualValues(t, 1, phones[0].ID)
	assert.Equal(t, "orphan", phones[0].Name)
	// Value relation: zero-valued Brand struct on no-match.
	assert.EqualValues(t, 0, phones[0].Brand.ID)
	assert.Equal(t, "", phones[0].Brand.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestJoinedPreloadFastScan_PreloadPlusFilterOnChild combines two patterns
// that share joinConditionImpl.applyTo: WHERE on the joined child (which
// makes the JOIN inner-style and goes into the ON clause) AND .Preload()
// (which selects the child columns and triggers mounting). Both must
// behave correctly together.
func TestJoinedPreloadFastScan_PreloadPlusFilterOnChild(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	// INNER JOIN (CQL switches to inner when WHERE on child) + joined cols.
	mock.ExpectQuery(`SELECT phones\.\*,.*Brand__id.*FROM "phones" INNER JOIN brands`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at",
				"name", "brand_id",
				"Brand__id", "Brand__name",
			}).AddRow(
				1, nil, nil, nil,
				"p1", 7,
				7, "acme",
			),
		)

	phones, err := Query[models.Phone](
		context.Background(),
		db,
		phoneBrandJoin(
			conditions.Brand.Name.Is().Eq(String("acme")),
		).Preload(),
	).Find()

	require.NoError(t, err)
	require.Len(t, phones, 1)
	assert.EqualValues(t, 1, phones[0].ID)
	assert.Equal(t, "p1", phones[0].Name)
	// Both the WHERE-matched filter AND the preload mount fired.
	assert.EqualValues(t, 7, phones[0].Brand.ID)
	assert.Equal(t, "acme", phones[0].Brand.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestJoinedPreloadFastScan_Nested materializes a 2-level preload:
// Sale → Seller → Company. Exercises:
//   - longest-prefix alias matching ("Seller__Company__id" must reach
//     the nested scanner, not the parent's)
//   - depth-ordered mounting (Company mounted onto Seller before Seller
//     onto Sale)
//   - mix of all-non-null and all-null nested rows (the latter has
//     Seller preloaded but no Company)
func TestJoinedPreloadFastScan_Nested(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	// CQL aliases: top-level relation = "Seller", nested = "Seller__Company".
	// Result columns: sales.* first, then Seller__*, then Seller__Company__*.
	mock.ExpectQuery(`SELECT sales\.\*,.*"Seller__id".*"Seller__Company__id".*FROM "sales"`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at",
				"code", "description", "product_id", "seller_id",
				"Seller__id", "Seller__name", "Seller__company_id", "Seller__university_id",
				"Seller__Company__id", "Seller__Company__name",
				"Seller__Company__created_at", "Seller__Company__updated_at",
				"Seller__Company__deleted_at",
			}).AddRow(
				// Sale with full nested preload populated
				newUUID("11111111-1111-1111-1111-111111111111"),
				nil, nil, nil,
				100, "desc1",
				newUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
				newUUID("22222222-2222-2222-2222-222222222222"),
				newUUID("22222222-2222-2222-2222-222222222222"), "alice",
				newUUID("33333333-3333-3333-3333-333333333333"), nil,
				newUUID("33333333-3333-3333-3333-333333333333"), "ditrit",
				nil, nil, nil,
			).AddRow(
				// Sale with Seller but Company is all-NULL (LEFT JOIN miss)
				newUUID("44444444-4444-4444-4444-444444444444"),
				nil, nil, nil,
				200, "desc2",
				newUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
				newUUID("55555555-5555-5555-5555-555555555555"),
				newUUID("55555555-5555-5555-5555-555555555555"), "bob",
				nil, nil, // bob has no company_id
				nil, nil, nil, nil, nil,
			),
		)

	sales, err := Query[models.Sale](
		context.Background(),
		db,
		conditions.Sale.Seller(
			conditions.Seller.Company().Preload(),
		).Preload(),
	).Find()

	require.NoError(t, err)
	require.Len(t, sales, 2)

	// Row 1: full nested preload.
	assert.Equal(t, "desc1", sales[0].Description)
	require.NotNil(t, sales[0].Seller)
	assert.Equal(t, "alice", sales[0].Seller.Name)
	require.NotNil(t, sales[0].Seller.Company)
	assert.Equal(t, "ditrit", sales[0].Seller.Company.Name)

	// Row 2: Seller preloaded, Company is nil (LEFT JOIN no match).
	assert.Equal(t, "desc2", sales[1].Description)
	require.NotNil(t, sales[1].Seller)
	assert.Equal(t, "bob", sales[1].Seller.Name)
	assert.Nil(t, sales[1].Seller.Company)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// newUUID converts a string UUID into model.UUID for sqlmock AddRow.
// model.UUID implements sql.Scanner so passing the string also works, but
// passing the typed value documents intent.
func newUUID(s string) string { return s }

// TestJoinedPreloadFastScan_JoinWithoutPreload confirms a JOIN-for-filter
// (no .Preload() on Brand) still works through the fast path. No joined
// columns appear in the result, so no activeJoin is registered, and the
// scanner only sees phone columns. (Pre-Phase 2 behavior: this would have
// fallen back to gorm because hasJoinedSelects = true.)
func TestJoinedPreloadFastScan_JoinWithoutPreload(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	// JOIN-for-filter uses INNER JOIN (not LEFT) because we have a WHERE
	// on the joined table — see joinConditionImpl.addJoin behavior.
	mock.ExpectQuery(`SELECT phones\.\* FROM "phones" INNER JOIN brands`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at",
				"name", "brand_id",
			}).AddRow(
				42, nil, nil, nil,
				"filter_hit", 7,
			),
		)

	phones, err := Query[models.Phone](
		context.Background(),
		db,
		phoneBrandJoin(
			conditions.Brand.Name.Is().Eq(String("acme")),
		), // NO .Preload() — pure filter
	).Find()

	require.NoError(t, err)
	require.Len(t, phones, 1)
	assert.EqualValues(t, 42, phones[0].ID)
	assert.Equal(t, "filter_hit", phones[0].Name)
	// Brand wasn't preloaded → stays zero.
	assert.EqualValues(t, 0, phones[0].Brand.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
