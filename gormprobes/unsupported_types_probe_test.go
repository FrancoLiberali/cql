package gormprobes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// cql-gen hard-fails (UnsupportedFieldError) when a model field is uintptr,
// complex64 or complex128. These probes lock in the justification: gorm itself
// rejects those types at AutoMigrate ("unsupported data type"), so failing at
// codegen is strictly better than emitting conditions that crash at query time.

func TestGormRejectsUintptr(t *testing.T) {
	t.Parallel()

	type uintptrModel struct {
		ID  uint
		Val uintptr
	}

	require.Error(t, openProbeDB(t).AutoMigrate(&uintptrModel{}),
		"gorm must reject uintptr at AutoMigrate — justifies cql-gen's hard-fail")
}

func TestGormRejectsComplex64(t *testing.T) {
	t.Parallel()

	type complex64Model struct {
		ID  uint
		Val complex64
	}

	require.Error(t, openProbeDB(t).AutoMigrate(&complex64Model{}),
		"gorm must reject complex64 at AutoMigrate — justifies cql-gen's hard-fail")
}

func TestGormRejectsComplex128(t *testing.T) {
	t.Parallel()

	type complex128Model struct {
		ID  uint
		Val complex128
	}

	require.Error(t, openProbeDB(t).AutoMigrate(&complex128Model{}),
		"gorm must reject complex128 at AutoMigrate — justifies cql-gen's hard-fail")
}
