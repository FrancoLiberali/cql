package gormprobes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// openProbeDB opens a fresh in-memory SQLite DB with logging silenced (the
// type-support probes deliberately trigger AutoMigrate errors, whose gorm log
// output would otherwise be noise).
func openProbeDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	return db
}
