package condition_test

// This file probes gorm/database-sql behavior for the Go basic types that
// CQL's scanner generator currently skips (uintptr, complex64, complex128).
// The goal: decide whether cql-gen should silently skip or hard-fail when it
// encounters such a field. If gorm itself can't round-trip a value of a
// given type, failing fast at codegen is better than producing conditions
// that crash at query time.
//
// Each test uses an in-memory sqlite DB (no fixtures, no Docker) so the
// probe runs in unit-test mode.

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openProbeDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	return db
}

type uintptrModel struct {
	ID  uint
	Val uintptr
}

type complex64Model struct {
	ID  uint
	Val complex64
}

type complex128Model struct {
	ID  uint
	Val complex128
}

func TestGormProbe_Uintptr(t *testing.T) {
	db := openProbeDB(t)

	migrateErr := db.AutoMigrate(&uintptrModel{})
	t.Logf("uintptr AutoMigrate err: %v", migrateErr)

	if migrateErr != nil {
		// If migration fails, no need to test further.
		return
	}

	createErr := db.Create(&uintptrModel{Val: 42}).Error
	t.Logf("uintptr Create err: %v", createErr)

	var out uintptrModel
	findErr := db.First(&out).Error
	t.Logf("uintptr First err: %v, val: %v", findErr, out.Val)

	// Record verdict for the report.
	if createErr == nil && findErr == nil && out.Val == 42 {
		t.Log("VERDICT uintptr: round-trips cleanly through gorm")
	} else {
		t.Log("VERDICT uintptr: fails through gorm")
	}
}

func TestGormProbe_Complex64(t *testing.T) {
	db := openProbeDB(t)

	migrateErr := db.AutoMigrate(&complex64Model{})
	t.Logf("complex64 AutoMigrate err: %v", migrateErr)

	if migrateErr != nil {
		assert.Error(t, migrateErr, "expected migrate to fail for complex64")
		t.Log("VERDICT complex64: gorm rejects at AutoMigrate")

		return
	}

	createErr := db.Create(&complex64Model{Val: complex64(1 + 2i)}).Error
	t.Logf("complex64 Create err: %v", createErr)

	var out complex64Model
	findErr := db.First(&out).Error
	t.Logf("complex64 First err: %v, val: %v", findErr, out.Val)

	if createErr == nil && findErr == nil && out.Val == complex64(1+2i) {
		t.Log("VERDICT complex64: round-trips cleanly through gorm")
	} else {
		t.Log("VERDICT complex64: fails through gorm")
	}
}

func TestGormProbe_Complex128(t *testing.T) {
	db := openProbeDB(t)

	migrateErr := db.AutoMigrate(&complex128Model{})
	t.Logf("complex128 AutoMigrate err: %v", migrateErr)

	if migrateErr != nil {
		assert.Error(t, migrateErr, "expected migrate to fail for complex128")
		t.Log("VERDICT complex128: gorm rejects at AutoMigrate")

		return
	}

	createErr := db.Create(&complex128Model{Val: complex128(3 + 4i)}).Error
	t.Logf("complex128 Create err: %v", createErr)

	var out complex128Model
	findErr := db.First(&out).Error
	t.Logf("complex128 First err: %v, val: %v", findErr, out.Val)

	if createErr == nil && findErr == nil && out.Val == complex128(3+4i) {
		t.Log("VERDICT complex128: round-trips cleanly through gorm")
	} else {
		t.Log("VERDICT complex128: fails through gorm")
	}
}
