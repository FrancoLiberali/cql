package gormdatabase

import "github.com/jackc/pgconn"

func IsDuplicateKeyError(err error) bool {
	// unique_violation code is equals to 23505
	return isPostgresError(err, "23505")
}

func isPostgresError(err error, errCode string) bool {
	postgresError, ok := err.(*pgconn.PgError)
	if ok {
		return postgresError.Code == errCode

	}
	return false
}
