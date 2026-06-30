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
