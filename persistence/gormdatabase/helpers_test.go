package gormdatabase

import (
	"errors"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsDuplicateError(t *testing.T) {
	assert.False(t, IsDuplicateKeyError(errors.New("voila")))
	assert.False(t, IsDuplicateKeyError(&pgconn.PgError{
		Code: "235252551514",
	}))
	assert.True(t, IsDuplicateKeyError(&pgconn.PgError{
		Code: "23505",
	}))
}

func Test_isPostgresError(t *testing.T) {
	var postgresErrorAsError error = &pgconn.PgError{Code: "1234"}
	assert.True(t, isPostgresError(
		postgresErrorAsError,
		"1234",
	))
	postgresErrorAsError = &pgconn.PgError{Code: ""}
	assert.False(t, isPostgresError(
		postgresErrorAsError,
		"1234",
	))
	assert.False(t, isPostgresError(errors.New("a classic error"), "1234"))
}
