// This file hosts the bench-only model graph that mirrors gorm's
// utils/tests.User schema verbatim (User + its associations). Kept in
// its own package so the type names (User, Account, Pet, Toy, Tools,
// Company, Language) don't collide with the existing cql/test/models
// integration-test types, and so cql-gen can produce bench-only
// conditions alongside without touching the main test suite.

package models

import (
	"database/sql"
	"time"

	"github.com/FrancoLiberali/cql/model"
)

// User mirrors gorm's utils/tests.User for the benchmark suite. The
// scalar column layout (id, timestamps, name, age, birthday, company_id,
// active) matches gorm's users table exactly, so per-row scan cost is
// directly comparable.
//
// The association graph is TRIMMED to the shapes cql-gen currently
// supports: Account (has-one), Pets (has-many), NamedPet (has-one),
// CompanyID (belongs-to scalar). Gorm's original User also has
// polymorphic (Toys, Tools), self-referential (Manager, Team), and
// many-to-many (Languages, Friends) relationships plus a Company field
// whose type isn't a cql model — cql-gen can't produce correct
// conditions/scanners for any of those yet. See cql/cql-gen/TODO.md.
// The benchmarks themselves don't query the association graph, so the
// missing fields don't affect BenchmarkFind/Scan/ScanSlice cost; they
// slightly under-charge BenchmarkCreate/Update's schema-parse work
// versus gorm's full-User row (a few hundred nanoseconds).
type User struct {
	model.UIntModelWithTimestamps

	Name      string
	Age       uint
	Birthday  *time.Time
	Account   Account
	Pets      []*Pet
	NamedPet  *Pet
	CompanyID *int
	Active    bool
}

type Account struct {
	model.UIntModelWithTimestamps

	UserID sql.NullInt64
	Number string
}

type Pet struct {
	model.UIntModelWithTimestamps

	UserID *model.UIntID
	Name   string
}
