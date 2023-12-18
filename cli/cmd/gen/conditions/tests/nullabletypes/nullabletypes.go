package nullabletypes

import (
	"database/sql"

	"github.com/ditrit/badaas/orm/model"
)

type NullableTypes struct {
	model.UUIDModel

	String  sql.NullString
	Int64   sql.NullInt64
	Int32   sql.NullInt32
	Int16   sql.NullInt16
	Byte    sql.NullByte
	Float64 sql.NullFloat64
	Bool    sql.NullBool
	Time    sql.NullTime
}
