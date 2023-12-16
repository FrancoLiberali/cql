package basicslicespointer

import "github.com/ditrit/badaas/orm"

type BasicSlicesPointer struct {
	orm.UUIDModel

	Bool       []*bool
	Int        []*int
	Int8       []*int8
	Int16      []*int16
	Int32      []*int32
	Int64      []*int64
	UInt       []*uint
	UInt8      []*uint8
	UInt16     []*uint16
	UInt32     []*uint32
	UInt64     []*uint64
	UIntptr    []*uintptr
	Float32    []*float32
	Float64    []*float64
	Complex64  []*complex64
	Complex128 []*complex128
	String     []*string
	Byte       []*byte
}
