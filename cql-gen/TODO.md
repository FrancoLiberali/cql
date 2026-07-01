# cql-gen — known limitations

Discovered while mirroring gorm's `utils/tests.User` model for
`cql/tests/benchmark_test.go`. The following relationship shapes are not
currently supported by cql-gen and cause its generated conditions or
scanner to reference fields/types that don't exist. Fix when the
benchmark suite (or a user request) needs them.

## Polymorphic associations

`Toys []Toy \`gorm:"polymorphic:Owner"\`` — cql-gen assumes the FK on
`Toy` is named `<Parent>ID` (i.e. `Toy.UserID`), but polymorphic
associations use a pair of `OwnerID` / `OwnerType` columns instead.
Generated code references `Toy.UserID` which doesn't exist.

Same issue with `polymorphicType:Type;polymorphicId:CustomID` on
`Tools`.

## Self-referential relationships

`Manager *User` and `Team []User \`gorm:"foreignkey:ManagerID"\`` cause
two problems:

1. Init-cycle: the generated `userTeamHasManyLoader` references `User`
   and `User` (the conditions var) references the loader — Go rejects
   the cycle at compile time.
2. FK naming: cql-gen expects `User.UserID` for the self-ref, but the
   FK is `ManagerID`.

## Many-to-many

`Languages []Language \`gorm:"many2many:UserSpeak"\`` and
`Friends []*User \`gorm:"many2many:user_friends"\`` — cql-gen doesn't
recognize the `many2many` tag; it falls through to has-many handling
and expects direct FKs that don't exist.

## Associations to non-cql-Model types

`Company Company` where `Company` doesn't embed a `model.Model` (no
`gorm.Model`, no `UIntModelWithTimestamps`) — cql-gen skips generating
conditions for `Company`, then the User conditions/scanner reference
`User.Company` field which was never defined. Either promote such
types to full cql models, or teach cql-gen to skip fields whose target
type has no conditions.

## Suggested resolutions

- Add `gorm:"polymorphic:"` tag detection → generate `OwnerID`/`OwnerType`
  where clauses instead of `<Parent>ID`.
- Add `gorm:"many2many:"` tag detection → emit a join-table read (two
  scanners + a middle table).
- Break the self-ref init cycle by lazily initializing the loader (init
  in a `func init()` block, not in a var initializer).
- Skip fields whose target type has no generated conditions (currently
  silently referenced as missing).
