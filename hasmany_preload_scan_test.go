package cql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

// TestHasManyPreloadFastScan_Basic exercises the simplest HasMany preload
// case: Company.Sellers.Preload(). Mocks both queries:
//
//   1) SELECT companies.* FROM "companies" WHERE ...
//   2) SELECT sellers.* FROM "sellers" WHERE sellers.company_id IN (?, ?)
//
// and asserts that the second query's results are grouped by company_id
// and mounted onto the right Company via the hand-written
// CompanySellersHasManyLoader.
func TestHasManyPreloadFastScan_Basic(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	companyA := newUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	companyB := newUUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	// Main query (Company list).
	mock.ExpectQuery(`SELECT companies\.\* FROM "companies"`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
			}).
				AddRow(companyA, nil, nil, nil, "acme").
				AddRow(companyB, nil, nil, nil, "globex"),
		)

	// HasMany loader query — fetches sellers for both companies in ONE go.
	// Argument order in WHERE IN(...) matches the order companies came back.
	mock.ExpectQuery(`SELECT sellers\.\* FROM "sellers" WHERE sellers\.company_id IN \(\$1,\$2\)`).
		WithArgs(companyA, companyB).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "name", "company_id", "university_id",
			}).
				AddRow(newUUID("11111111-1111-1111-1111-111111111111"), "alice", companyA, nil).
				AddRow(newUUID("22222222-2222-2222-2222-222222222222"), "bob", companyA, nil).
				AddRow(newUUID("33333333-3333-3333-3333-333333333333"), "carol", companyB, nil),
		)

	companies, err := Query[models.Company](
		context.Background(),
		db,
		conditions.Company.Sellers.Preload(),
	).Find()

	require.NoError(t, err)
	require.Len(t, companies, 2)

	assert.Equal(t, "acme", companies[0].Name)
	require.NotNil(t, companies[0].Sellers)
	sellers0 := *companies[0].Sellers
	require.Len(t, sellers0, 2)
	assert.ElementsMatch(t, []string{"alice", "bob"}, []string{sellers0[0].Name, sellers0[1].Name})

	assert.Equal(t, "globex", companies[1].Name)
	require.NotNil(t, companies[1].Sellers)
	sellers1 := *companies[1].Sellers
	require.Len(t, sellers1, 1)
	assert.Equal(t, "carol", sellers1[0].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestHasManyPreloadFastScan_WithNestedJoinedPreload exercises the supported
// "nesting" shape for HasMany: Company.Sellers.Preload(Seller.University().Preload()).
// The inner Seller.University() is a JoinCondition that flows through
// HasManyLoader.BuildQuery as the `nested` slice, becoming a condition on
// the internal Query[Seller]. That Query[Seller] then registers its own
// activeJoin for University and runs a single JOIN-style child SELECT.
//
// Validates that the whole chain stays on the fast path: one main query +
// one child query (with JOIN), no per-row reflection anywhere.
func TestHasManyPreloadFastScan_WithNestedJoinedPreload(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	companyA := newUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	companyB := newUUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	uniX := newUUID("99999999-9999-9999-9999-999999999999")

	mock.ExpectQuery(`SELECT companies\.\* FROM "companies"`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
			}).
				AddRow(companyA, nil, nil, nil, "acme").
				AddRow(companyB, nil, nil, nil, "globex"),
		)

	// Child query MUST include the LEFT JOIN + University__* select +
	// IN on company_id. This is the proof that the nested joined preload
	// is composed into the child query and not fired as a separate gorm
	// preload.
	mock.ExpectQuery(`SELECT sellers\.\*,.*"University__id".*FROM "sellers".*LEFT JOIN universities.*sellers\.company_id IN \(\$1,\$2\)`).
		WithArgs(companyA, companyB).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "name", "company_id", "university_id",
				"University__id", "University__name",
			}).
				AddRow(newUUID("11111111-1111-1111-1111-111111111111"), "alice", companyA, uniX, uniX, "MIT").
				AddRow(newUUID("22222222-2222-2222-2222-222222222222"), "bob", companyA, nil, nil, nil).
				AddRow(newUUID("33333333-3333-3333-3333-333333333333"), "carol", companyB, uniX, uniX, "MIT"),
		)

	companies, err := Query[models.Company](
		context.Background(),
		db,
		conditions.Company.Sellers.Preload(
			conditions.Seller.University().Preload(),
		),
	).Find()

	require.NoError(t, err)
	require.Len(t, companies, 2)

	// Company A → 2 sellers, alice has University, bob doesn't.
	require.NotNil(t, companies[0].Sellers)
	sellersA := *companies[0].Sellers
	require.Len(t, sellersA, 2)

	alice := findSellerByName(sellersA, "alice")
	require.NotNil(t, alice)
	require.NotNil(t, alice.University)
	assert.Equal(t, "MIT", alice.University.Name)

	bob := findSellerByName(sellersA, "bob")
	require.NotNil(t, bob)
	assert.Nil(t, bob.University)

	// Company B → 1 seller (carol with University).
	require.NotNil(t, companies[1].Sellers)
	sellersB := *companies[1].Sellers
	require.Len(t, sellersB, 1)
	require.NotNil(t, sellersB[0].University)
	assert.Equal(t, "MIT", sellersB[0].University.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func findSellerByName(sellers []models.Seller, name string) *models.Seller {
	for i := range sellers {
		if sellers[i].Name == name {
			return &sellers[i]
		}
	}

	return nil
}

// TestHasManyPreloadFastScan_First validates that single-row finishers
// (First/Take/Last) also run registered HasMany loaders against the lone
// parent. Regression for the scanOne path — symmetrically wires
// runHasManyLoaders just like findWith does.
func TestHasManyPreloadFastScan_First(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	companyA := newUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	mock.ExpectQuery(`SELECT companies\.\* FROM "companies".*LIMIT \$\d`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "created_at", "updated_at", "deleted_at", "name",
			}).AddRow(companyA, nil, nil, nil, "acme"),
		)

	mock.ExpectQuery(`SELECT sellers\.\* FROM "sellers" WHERE sellers\.company_id IN \(\$1\)`).
		WithArgs(companyA).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "name", "company_id", "university_id",
			}).
				AddRow(newUUID("11111111-1111-1111-1111-111111111111"), "alice", companyA, nil).
				AddRow(newUUID("22222222-2222-2222-2222-222222222222"), "bob", companyA, nil),
		)

	company, err := Query[models.Company](
		context.Background(),
		db,
		conditions.Company.Sellers.Preload(),
	).First()

	require.NoError(t, err)
	require.NotNil(t, company)
	assert.Equal(t, "acme", company.Name)

	require.NotNil(t, company.Sellers)
	require.Len(t, *company.Sellers, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestHasManyPreloadFastScan_EmptyParentSet verifies the loader short-circuits
// when there are no parents to load children for — no child query is sent.
func TestHasManyPreloadFastScan_EmptyParentSet(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT companies\.\* FROM "companies"`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "name",
		}))

	// NO second query mocked — loader must skip when parents is empty.

	companies, err := Query[models.Company](
		context.Background(),
		db,
		conditions.Company.Sellers.Preload(),
	).Find()

	require.NoError(t, err)
	assert.Empty(t, companies)
	assert.NoError(t, mock.ExpectationsWereMet())
}
