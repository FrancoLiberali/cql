package cql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

// These tests verify the typed-Scanner fast path for Find/First/Take/Last
// is selected (via the scanner pointer carried on conditions.Brand fields)
// and that it scans rows into models correctly without going through gorm's
// reflective scan.
//
// brand_scanner.go in test/conditions/ wires the scanner onto Brand fields.
// Any query built from those fields will pick up the scanner in
// resolveScanner and dispatch findWith/firstWith/takeWith/lastWith.

func openMockedDB(t *testing.T) (*DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := Open(postgres.New(postgres.Config{Conn: conn}))
	require.NoError(t, err)

	return db, mock, func() { conn.Close() }
}

func TestScannerFastPath_Find(t *testing.T) {
	db, mock, cleanup := openMockedDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT brands\.\* FROM "brands" WHERE brands\.name = \(\$1\)`).
		WithArgs("acme").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "acme").
				AddRow(2, "acme"),
		)

	brands, err := Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(String("acme")),
	).Find()

	require.NoError(t, err)
	require.Len(t, brands, 2)
	assert.EqualValues(t, 1, brands[0].ID)
	assert.Equal(t, "acme", brands[0].Name)
	assert.EqualValues(t, 2, brands[1].ID)
	assert.Equal(t, "acme", brands[1].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScannerFastPath_First(t *testing.T) {
	db, mock, cleanup := openMockedDB(t)
	defer cleanup()

	// First() adds ORDER BY primary key ASC and LIMIT 1.
	mock.ExpectQuery(`SELECT brands\.\* FROM "brands" WHERE brands\.name = \(\$1\) ORDER BY .* LIMIT \$2`).
		WithArgs("acme", 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(7, "acme"),
		)

	brand, err := Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(String("acme")),
	).First()

	require.NoError(t, err)
	require.NotNil(t, brand)
	assert.EqualValues(t, 7, brand.ID)
	assert.Equal(t, "acme", brand.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScannerFastPath_FirstReturnsErrRecordNotFound(t *testing.T) {
	db, mock, cleanup := openMockedDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT brands\.\* FROM "brands" WHERE brands\.name = \(\$1\) ORDER BY .* LIMIT \$2`).
		WithArgs("missing", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	_, err := Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(String("missing")),
	).First()

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScannerFastPath_Take(t *testing.T) {
	db, mock, cleanup := openMockedDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT brands\.\* FROM "brands" WHERE brands\.name = \(\$1\) LIMIT \$2`).
		WithArgs("acme", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(42, "acme"))

	brand, err := Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(String("acme")),
	).Take()

	require.NoError(t, err)
	assert.EqualValues(t, 42, brand.ID)
	assert.Equal(t, "acme", brand.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScannerFastPath_Last(t *testing.T) {
	db, mock, cleanup := openMockedDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT brands\.\* FROM "brands" WHERE brands\.name = \(\$1\) ORDER BY .*DESC.* LIMIT \$2`).
		WithArgs("acme", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(99, "acme"))

	brand, err := Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(String("acme")),
	).Last()

	require.NoError(t, err)
	assert.EqualValues(t, 99, brand.ID)
	assert.Equal(t, "acme", brand.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test that an extra column not known to the scanner (e.g. a projection
// added for Postgres ORDER BY) is drained into NullSink and ignored,
// without breaking the scan.
func TestScannerFastPath_ToleratesUnknownColumn(t *testing.T) {
	db, mock, cleanup := openMockedDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT brands\.\* FROM "brands" WHERE brands\.name = \(\$1\)`).
		WithArgs("acme").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "extra_alias"}).
				AddRow(1, "acme", "ignored"),
		)

	brands, err := Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(String("acme")),
	).Find()

	require.NoError(t, err)
	require.Len(t, brands, 1)
	assert.EqualValues(t, 1, brands[0].ID)
	assert.Equal(t, "acme", brands[0].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}
