package nullabletypes

import (
	"database/sql"

	"github.com/ditrit/badaas/orm"
)

type NullableTypes struct {
	orm.UUIDModel

	String  sql.NullString
	Int64   sql.NullInt64
	Int32   sql.NullInt32
	Int16   sql.NullInt16
	Byte    sql.NullByte
	Float64 sql.NullFloat64
	Bool    sql.NullBool
	Time    sql.NullTime
}
