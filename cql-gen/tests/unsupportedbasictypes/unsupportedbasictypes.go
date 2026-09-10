package unsupportedbasictypes

import "github.com/FrancoLiberali/cql/model"

// UnsupportedBasicTypes contains exactly the Go basic types that cannot be
// stored in any SQL database:
//   - uintptr: a memory address; gorm rejects at AutoMigrate
//   - complex64 / complex128: no SQL type; gorm rejects at AutoMigrate
//
// cql-gen must fail loudly when it sees a model like this so the user finds
// out at codegen rather than at their first DB call. Users who want
// complex-shaped data can define a wrapper type implementing sql.Scanner +
// driver.Valuer (see the customtype fixture) — that goes through a different
// code path and works fine.
type UnsupportedBasicTypes struct {
	model.UUIDModel

	UIntptr    uintptr
	Complex64  complex64
	Complex128 complex128
}
