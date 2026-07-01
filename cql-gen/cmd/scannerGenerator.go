package cmd

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"

	"github.com/FrancoLiberali/cql/cql-gen/cmd/log"
)

// UnsupportedFieldError is returned when cql-gen encounters a field whose Go
// type cannot be round-tripped through any SQL database (verified: gorm itself
// rejects it at AutoMigrate). Failing at codegen surfaces the problem now
// rather than at the user's first migrate or query call.
type UnsupportedFieldError struct {
	Model string
	Field string
	Type  string
}

func (e *UnsupportedFieldError) Error() string {
	return fmt.Sprintf(
		"cql-gen: model %s, field %s: Go type %q has no SQL representation "+
			"(gorm rejects this type at AutoMigrate). Remove the field or change its type.",
		e.Model, e.Field, e.Type,
	)
}

// DuplicateColumnError is returned when two model fields would resolve to the
// same SQL column name (e.g. an embedded struct with no prefix shadowing a
// top-level column). Such a model is unusable: gorm would reject the table
// at AutoMigrate, and the scanner can't disambiguate which field a value
// belongs to. Fail at codegen.
type DuplicateColumnError struct {
	Model   string
	Column  string
	FieldA  string
	FieldB  string
}

func (e *DuplicateColumnError) Error() string {
	return fmt.Sprintf(
		"cql-gen: model %s: column %q is produced by both field %s and field %s "+
			"(gorm would reject the table at AutoMigrate). Add a gorm:\"column:...\" tag "+
			"or an embeddedPrefix to disambiguate.",
		e.Model, e.Column, e.FieldA, e.FieldB,
	)
}

// ScannerGenerator emits a Scanner[T] for a model plus an init() that rewires
// the conditions struct's Field values to carry the scanner pointer. Skips
// emission silently when the model has fields of types we don't yet know how
// to materialize without reflection — those models keep working through the
// gorm fallback path.
type ScannerGenerator struct {
	object     types.Object
	objectType Type
}

func NewScannerGenerator(object types.Object) *ScannerGenerator {
	return &ScannerGenerator{
		object:     object,
		objectType: Type{Type: object.Type()},
	}
}

// scannerField describes one column's worth of scan/assign/release code
// for a single model field. valueExpr is what ScanValues emits to obtain
// the per-row destination cell (a pooled `condition.AcquireNullX()` call
// where possible; falls back to `new(sql.NullX)` for unpooled types).
// releaseFn is the body of one case in the generated ReleaseValues method
// — it returns the cell to its sync.Pool after AssignValues has copied
// the data out. Mirrors gorm's per-cell `field.NewValuePool.Put` at
// scan.go:111. May be nil for cells that aren't pooled (e.g. custom
// scanner types wrapped in NullableScanner).
type scannerField struct {
	columnName string             // final SQL column name (snake_case + any prefix override)
	valueExpr  *jen.Statement     // ScanValues cell allocator (Acquire call when pooled)
	assignFn   scannerAssignFunc  // AssignValues body for this column's case
	releaseFn  scannerAssignFunc  // ReleaseValues body for this column's case; nil = no-op
}

// scannerAssignFunc receives the case index variable name and the values slice
// access expression. It returns the body statements for one switch case.
type scannerAssignFunc func(idxIdent string) []jen.Code

// Validate runs all classification and consistency checks on the model
// without writing any files. Called from main.go BEFORE conditions
// generation so that a fatal scanner error (unsupported type, duplicate
// column) leaves no stray output on disk. Returns:
//   - (true, nil) if a scanner CAN be emitted for this model
//   - (false, nil) if scanner emission should be skipped silently (model not
//     a cql model, or has fields whose type is valid SQL but the scanner
//     generator doesn't yet support — gorm fallback still works)
//   - (false, *UnsupportedFieldError) if a field's Go type can't be
//     persisted at all (uintptr, complex*)
//   - (false, *DuplicateColumnError) if two fields resolve to the same
//     column (gorm would reject the table; scanner couldn't disambiguate)
func (sg ScannerGenerator) Validate() (bool, error) {
	flatFields, err := sg.flattenFields()
	if err != nil {
		return false, nil //nolint:nilerr // not a cql model — skip silently
	}

	scannerFields, err := sg.classifyAll(flatFields)
	if err != nil {
		return false, err
	}

	if scannerFields == nil {
		return false, nil
	}

	if err := sg.checkDuplicateColumns(flatFields, scannerFields); err != nil {
		return false, err
	}

	return true, nil
}

// Into emits the scanner file alongside the conditions file. Returns:
//   - (true, nil) if a scanner was emitted
//   - (false, nil) if the model has fields whose type is valid SQL but the
//     scanner generator doesn't yet support (e.g. pointer-to-named) — these
//     keep working through gorm's fallback scan path
//   - (false, *UnsupportedFieldError) / *DuplicateColumnError if validation
//     fails — callers should normally have called Validate first to catch
//     these before any other file was written
//   - (false, otherErr) on I/O failure
func (sg ScannerGenerator) Into(destPkg, dir string) (bool, error) {
	flatFields, err := sg.flattenFields()
	if err != nil {
		return false, nil //nolint:nilerr // not a cql model — skip silently
	}

	scannerFields, err := sg.classifyAll(flatFields)
	if err != nil {
		return false, err
	}

	if scannerFields == nil {
		log.Logger.Debugf("Scanner skipped for %s: model has fields whose Go type isn't yet supported by fast scan",
			sg.object.Name())

		return false, nil
	}

	if err := sg.checkDuplicateColumns(flatFields, scannerFields); err != nil {
		return false, err
	}

	file := NewFile(destPkg, fileNameForScanner(dir, sg.object.Name()))

	objectQual := jen.Qual(
		getRelativePackagePath(destPkg, sg.objectType),
		sg.objectType.Name(),
	)

	scannerVar := strcase.ToCamel(sg.object.Name()) + "Scanner"

	sg.emitScannerVar(file, objectQual, scannerVar, scannerFields)
	sg.emitRelationScanners(file, destPkg)
	sg.emitInitRewiring(file, scannerVar)

	if err := file.Save(); err != nil {
		return false, err
	}

	return true, nil
}

// flattenFields returns the top-level fields of the model with embedded base
// models / structs walked one level deep, matching what conditionsGenerator
// produces. Fields whose path requires deeper recursion (or that are
// relations/collections) are returned with their leaf metadata so
// classifyAll can decide whether they're scan-supported.
func (sg ScannerGenerator) flattenFields() ([]Field, error) {
	fields, err := getFields(sg.objectType)
	if err != nil {
		return nil, err
	}

	flat := []Field{}

	for _, f := range fields {
		flat = append(flat, sg.flattenOne(f, "", nil)...)
	}

	return flat, nil
}

// flattenOne expands an embedded struct one level deep, carrying the gorm
// column prefix and the Go access path from the model root down to each
// leaf. Mirrors generateForEmbeddedField for column handling, and
// additionally tracks the access path so the emitted scanner can write to
// dest.GormEmbedded.Int instead of (incorrectly) dest.Int.
//
// Access-path rules:
//   - Go-anonymous embed (`SomeType` with no field name): fields are
//     promoted — child path inherits parent's path unchanged.
//   - Named gorm-tagged embed (`Foo SomeType `gorm:"embedded"``): child
//     path prepends the parent's Go field name.
func (sg ScannerGenerator) flattenOne(f Field, parentColumnPrefix string, parentAccessPath []string) []Field {
	if !f.Embedded {
		// Leaf field. Final column prefix and access path are fixed here.
		f.ColumnPrefix = parentColumnPrefix + f.ColumnPrefix
		f.AccessPath = append(append([]string{}, parentAccessPath...), f.Name)

		return []Field{f}
	}

	embeddedStruct, ok := f.Type.Underlying().(*types.Struct)
	if !ok {
		return []Field{f}
	}

	sub, err := getStructFields(embeddedStruct)
	if err != nil {
		return nil
	}

	childPrefix := parentColumnPrefix
	if !isBaseModel(f.TypeString()) {
		childPrefix += f.Tags.getEmbeddedPrefix()
	}

	childAccessPath := parentAccessPath
	if !f.GoAnonymous {
		// Named gorm-embedded — inner fields aren't promoted.
		childAccessPath = append(append([]string{}, parentAccessPath...), f.Name)
	}

	out := []Field{}
	for _, s := range sub {
		out = append(out, sg.flattenOne(s, childPrefix, childAccessPath)...)
	}

	return out
}

// checkDuplicateColumns returns a DuplicateColumnError when two scanner
// fields map to the same SQL column name. The companion flatFields slice is
// used to recover the Go field paths for the error message.
func (sg ScannerGenerator) checkDuplicateColumns(flatFields []Field, scanFields []scannerField) error {
	seen := map[string]int{}

	for i, sf := range scanFields {
		if prev, dup := seen[sf.columnName]; dup {
			return &DuplicateColumnError{
				Model:  sg.object.Name(),
				Column: sf.columnName,
				FieldA: pathString(flatFields[prev]),
				FieldB: pathString(flatFields[i]),
			}
		}

		seen[sf.columnName] = i
	}

	return nil
}

func pathString(f Field) string {
	if len(f.AccessPath) == 0 {
		return f.Name
	}

	out := f.AccessPath[0]
	for _, p := range f.AccessPath[1:] {
		out += "." + p
	}

	return out
}

// classifyAll attempts to produce a scannerField for every field. Returns:
//   - ([]scannerField, nil) on success (one entry per scannable field;
//     non-scannable fields are skipped silently so the scanner can still
//     materialize a partial model — relations and unsupported-but-valid
//     types are left at their zero value, identical to what gorm does
//     when the user doesn't preload them)
//   - (nil, *UnsupportedFieldError) if any field has a Go type that no SQL
//     database can persist (caller should fail the build)
//   - (nil, nil) if there are no scannable fields at all (rare; emission
//     is skipped silently)
func (sg ScannerGenerator) classifyAll(fields []Field) ([]scannerField, error) {
	out := make([]scannerField, 0, len(fields))

	for _, f := range fields {
		sf, ok, err := sg.classifyField(f)
		if err != nil {
			return nil, err
		}

		if !ok {
			// Field is a relation, an unsupported-but-valid-SQL type, or
			// otherwise outside fast-scan scope. Skip just this field;
			// keep going so other fields still get scan cases.
			continue
		}

		out = append(out, sf)
	}

	if len(out) == 0 {
		return nil, nil
	}

	return out, nil
}

// classifyField determines the scan/assign code for one field. Returns:
//   - (sf, true, nil) on success
//   - (zero, false, *UnsupportedFieldError) if the field's Go type can never
//     work in any SQL DB (uintptr, complex*)
//   - (zero, false, nil) if the field's Go type is valid SQL but not yet
//     handled by the fast-path generator
func (sg ScannerGenerator) classifyField(field Field) (scannerField, bool, error) {
	colName := field.getColumnName()
	if colName == "" {
		colName = strcase.ToSnake(field.Name)
	}
	// gorm embeddedPrefix (e.g. "gorm_embedded_") is prepended to every
	// embedded leaf's actual column name; the scanner must switch on that
	// final name so the row maps to dest correctly.
	colName = field.ColumnPrefix + colName

	dest := destAccessor(field.AccessPath, field.Name)

	switch t := field.Type.Type.(type) {
	case *types.Basic:
		sf, ok, err := classifyBasic(colName, dest, t)
		if err != nil {
			// Annotate with the owning model so the error message points
			// at User.Foo, not just the bare field.
			if ufe, isUFE := err.(*UnsupportedFieldError); isUFE && ufe.Model == "" {
				ufe.Model = sg.object.Name()
				ufe.Field = field.Name
			}

			return scannerField{}, false, err
		}

		return sf, ok, nil
	case *types.Named:
		sf, ok := sg.classifyNamed(colName, dest, field)

		return sf, ok, nil
	case *types.Pointer:
		// Pointer-to-basic: scan into the matching sql.Null* wrapper,
		// assign pointer when valid, nil otherwise.
		if basic, ok := t.Elem().(*types.Basic); ok {
			sf, ok := classifyPointerToBasic(colName, dest, basic)

			return sf, ok, nil
		}

		// Pointer-to-named (most commonly *model.UUID or *model.UIntID
		// for nullable foreign keys): scan into the named type itself,
		// then assign nil-or-pointer based on a zero-value sentinel.
		if named, ok := t.Elem().(*types.Named); ok {
			named := Type{Type: named}
			switch named.String() {
			case modelPath + "." + uuid:
				return pointerUUIDField(colName, dest), true, nil
			case modelPath + "." + uIntID:
				return pointerUIntIDField(colName, dest), true, nil
			}
			// Unknown pointer-to-named — skip (gorm fallback handles).
		}

		return scannerField{}, false, nil
	case *types.Slice:
		if elem, ok := t.Elem().(*types.Basic); ok && elem.Kind() == types.Uint8 {
			return classifyByteSlice(colName, dest), true, nil
		}

		return scannerField{}, false, nil
	default:
		return scannerField{}, false, nil
	}
}

// destAccessor builds the `dest.X.Y.Z` jen statement chain for assigning to
// a leaf field. Falls back to `dest.<leafName>` when the AccessPath is empty
// (e.g. for fields the older path didn't populate).
func destAccessor(accessPath []string, leafName string) *jen.Statement {
	if len(accessPath) == 0 {
		accessPath = []string{leafName}
	}

	expr := jen.Id("dest")
	for _, part := range accessPath {
		expr = expr.Dot(part)
	}

	return expr
}

// classifyBasic returns:
//   - (sf, true, nil) for kinds we can scan
//   - (zero, false, *UnsupportedFieldError) for kinds that no SQL DB can
//     persist (uintptr, complex64, complex128) — the build should fail
//   - (zero, false, nil) for anything else (currently unreachable for Basic)
func classifyBasic(col string, dest *jen.Statement, t *types.Basic) (scannerField, bool, error) {
	switch t.Kind() {
	case types.Bool:
		return nullValueField(col, dest, "NullBool", "Bool", nil), true, nil
	case types.String:
		return nullValueField(col, dest, "NullString", "String", nil), true, nil
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		castName := basicKindName(t.Kind())

		return nullValueField(col, dest, "NullInt64", "Int64", castFn(castName)), true, nil
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		castName := basicKindName(t.Kind())

		return nullValueField(col, dest, "NullInt64", "Int64", castFn(castName)), true, nil
	case types.Float32:
		return nullValueField(col, dest, "NullFloat64", "Float64", castFn("float32")), true, nil
	case types.Float64:
		return nullValueField(col, dest, "NullFloat64", "Float64", nil), true, nil
	case types.Uintptr, types.Complex64, types.Complex128:
		// Verified empirically: gorm itself rejects these at AutoMigrate
		// with "unsupported data type". Fail at codegen so the user finds
		// out now instead of at the first DB call.
		return scannerField{}, false, &UnsupportedFieldError{Type: t.Name()}
	default:
		return scannerField{}, false, nil
	}
}

func basicKindName(k types.BasicKind) string {
	switch k {
	case types.Int:
		return "int"
	case types.Int8:
		return "int8"
	case types.Int16:
		return "int16"
	case types.Int32:
		return "int32"
	case types.Int64:
		return "int64"
	case types.Uint:
		return "uint"
	case types.Uint8:
		return "uint8"
	case types.Uint16:
		return "uint16"
	case types.Uint32:
		return "uint32"
	case types.Uint64:
		return "uint64"
	default:
		return ""
	}
}

// castFn wraps the scanned value with a type conversion (e.g. int(v.Int64)).
// Pass nil for no conversion.
func castFn(typeName string) func(inner *jen.Statement) *jen.Statement {
	return func(inner *jen.Statement) *jen.Statement {
		return jen.Id(typeName).Call(inner)
	}
}

// nullValueField builds a scannerField for a column scanned into one of the
// sql.Null* wrappers. cast wraps the .X access (e.g. v.Int64) before the
// assignment; pass nil to assign the raw value. dest is the prebuilt
// `dest.X.Y.Z` chain so embedded gorm-tagged fields land in the right slot.
//
// Cells are pooled via condition.AcquireNullX / ReleaseNullX so per-row
// allocation cost amortizes across the result set.
func nullValueField(
	col string, dest *jen.Statement,
	nullType, valueAccessor string,
	cast func(*jen.Statement) *jen.Statement,
) scannerField {
	return scannerField{
		columnName: col,
		valueExpr:  jen.Qual(conditionPath, "Acquire"+nullType).Call(),
		assignFn: func(idx string) []jen.Code {
			vAssertion := jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
				jen.Id("values").Index(jen.Id(idx)).Assert(
					jen.Op("*").Qual("database/sql", nullType),
				),
			)
			notOk := jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad type %T"),
					jen.Id("values").Index(jen.Id(idx)),
				)),
			)
			rhs := jen.Id("v").Dot(valueAccessor)
			if cast != nil {
				rhs = cast(rhs)
			}

			assign := dest.Clone().Op("=").Add(rhs)

			validGuard := jen.If(jen.Id("v").Dot("Valid")).Block(assign)

			return []jen.Code{vAssertion, notOk, validGuard}
		},
		releaseFn: releaseNullFn(col, nullType),
	}
}

// releaseNullFn emits the body of one case in ReleaseValues for a pooled
// sql.Null* cell. Type-asserts and calls condition.ReleaseNullX.
func releaseNullFn(col, nullType string) scannerAssignFunc {
	return func(idx string) []jen.Code {
		return []jen.Code{
			jen.If(
				jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
					jen.Id("values").Index(jen.Id(idx)).Assert(
						jen.Op("*").Qual("database/sql", nullType),
					),
				),
				jen.Id("ok"),
			).Block(
				jen.Qual(conditionPath, "Release"+nullType).Call(jen.Id("v")),
			),
		}
	}
}

func classifyPointerToBasic(col string, dest *jen.Statement, t *types.Basic) (scannerField, bool) {
	sf, ok, _ := classifyBasic(col, dest, t)
	if !ok {
		return scannerField{}, false
	}

	// Override the assign to set a pointer to the value when valid, nil otherwise.
	innerKind := basicKindName(t.Kind())
	if innerKind == "" {
		switch t.Kind() {
		case types.Bool:
			innerKind = "bool"
		case types.String:
			innerKind = "string"
		case types.Float32:
			innerKind = "float32"
		case types.Float64:
			innerKind = "float64"
		default:
			return scannerField{}, false
		}
	}

	nullType, valueAccessor := nullTypeFor(t.Kind())
	if nullType == "" {
		return scannerField{}, false
	}

	sf.assignFn = func(idx string) []jen.Code {
		vAssertion := jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
			jen.Id("values").Index(jen.Id(idx)).Assert(
				jen.Op("*").Qual("database/sql", nullType),
			),
		)
		notOk := jen.If(jen.Op("!").Id("ok")).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(
				jen.Lit("cql scanner "+col+": bad type %T"),
				jen.Id("values").Index(jen.Id(idx)),
			)),
		)
		validBranch := jen.If(jen.Id("v").Dot("Valid")).Block(
			jen.Id("tmp").Op(":=").Id(innerKind).Call(jen.Id("v").Dot(valueAccessor)),
			dest.Clone().Op("=").Op("&").Id("tmp"),
		).Else().Block(
			dest.Clone().Op("=").Nil(),
		)

		return []jen.Code{vAssertion, notOk, validBranch}
	}
	sf.releaseFn = releaseNullFn(col, nullType)

	return sf, true
}

func nullTypeFor(k types.BasicKind) (string, string) {
	switch k {
	case types.Bool:
		return "NullBool", "Bool"
	case types.String:
		return "NullString", "String"
	case types.Float32, types.Float64:
		return "NullFloat64", "Float64"
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		return "NullInt64", "Int64"
	default:
		return "", ""
	}
}

func classifyByteSlice(col string, dest *jen.Statement) scannerField {
	return scannerField{
		columnName: col,
		valueExpr:  jen.New(jen.Index().Byte()),
		assignFn: func(idx string) []jen.Code {
			vAssertion := jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
				jen.Id("values").Index(jen.Id(idx)).Assert(jen.Op("*").Index().Byte()),
			)
			notOk := jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad type %T"),
					jen.Id("values").Index(jen.Id(idx)),
				)),
			)
			assign := dest.Clone().Op("=").Op("*").Id("v")

			return []jen.Code{vAssertion, notOk, assign}
		},
	}
}

// classifyNamed handles named types: cql model IDs (UIntID, UUID), time.Time,
// sql.Null* wrappers, and gorm custom types implementing Scanner+Valuer.
// Returns ok=false for relations to other CQL models — those would require
// joined-preload handling, which is outside the fast-path scope.
func (sg ScannerGenerator) classifyNamed(col string, dest *jen.Statement, field Field) (scannerField, bool) {
	if _, err := field.Type.CQLModelStruct(); err == nil {
		// Relation to another model — falls outside fast-path scope.
		return scannerField{}, false
	}

	typeStr := field.Type.String()

	switch typeStr {
	case modelPath + "." + uIntID:
		return uintIDField(col, dest), true
	case modelPath + "." + uuid:
		return uuidField(col, dest), true
	case "time.Time":
		return timeField(col, dest), true
	}

	if field.Type.IsSQLNullableType() {
		return nullableTypeField(col, dest, field.Type), true
	}

	if field.Type.IsGormCustomType() {
		return customScannerField(col, dest, field.Type), true
	}

	return scannerField{}, false
}

func uintIDField(col string, dest *jen.Statement) scannerField {
	return scannerField{
		columnName: col,
		valueExpr:  jen.Qual(conditionPath, "AcquireNullInt64").Call(),
		assignFn: func(idx string) []jen.Code {
			vAssertion := jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
				jen.Id("values").Index(jen.Id(idx)).Assert(jen.Op("*").Qual("database/sql", "NullInt64")),
			)
			notOk := jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad type %T"),
					jen.Id("values").Index(jen.Id(idx)),
				)),
			)
			assign := dest.Clone().Op("=").Add(
				jen.Qual(modelPath, uIntID).Call(jen.Id("v").Dot("Int64")), //nolint:gosec
			)

			return []jen.Code{vAssertion, notOk, assign}
		},
		releaseFn: releaseNullFn(col, "NullInt64"),
	}
}

// pointerUUIDField handles `*model.UUID` (nullable FK). Scans into
// model.UUID directly; treats NilUUID as the "no value" sentinel so a
// NULL column maps to a nil pointer on the destination.
func pointerUUIDField(col string, dest *jen.Statement) scannerField {
	return scannerField{
		columnName: col,
		valueExpr:  jen.Qual(conditionPath, "AcquireUUID").Call(),
		assignFn: func(idx string) []jen.Code {
			vAssertion := jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
				jen.Id("values").Index(jen.Id(idx)).Assert(jen.Op("*").Qual(modelPath, uuid)),
			)
			notOk := jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad type %T"),
					jen.Id("values").Index(jen.Id(idx)),
				)),
			)
			branch := jen.If(jen.Op("*").Id("v").Op("==").Qual(modelPath, "NilUUID")).Block(
				dest.Clone().Op("=").Nil(),
			).Else().Block(
				jen.Id("tmp").Op(":=").Op("*").Id("v"),
				dest.Clone().Op("=").Op("&").Id("tmp"),
			)

			return []jen.Code{vAssertion, notOk, branch}
		},
		releaseFn: releaseUUIDFn(col),
	}
}

// pointerUIntIDField handles `*model.UIntID` (nullable FK). Scans into
// sql.NullInt64 since model.UIntID has no Scanner interface itself.
func pointerUIntIDField(col string, dest *jen.Statement) scannerField {
	return scannerField{
		columnName: col,
		valueExpr:  jen.Qual(conditionPath, "AcquireNullInt64").Call(),
		assignFn: func(idx string) []jen.Code {
			vAssertion := jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
				jen.Id("values").Index(jen.Id(idx)).Assert(jen.Op("*").Qual("database/sql", "NullInt64")),
			)
			notOk := jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad type %T"),
					jen.Id("values").Index(jen.Id(idx)),
				)),
			)
			branch := jen.If(jen.Id("v").Dot("Valid")).Block(
				jen.Id("tmp").Op(":=").Qual(modelPath, uIntID).Call(jen.Id("v").Dot("Int64")), //nolint:gosec
				dest.Clone().Op("=").Op("&").Id("tmp"),
			).Else().Block(
				dest.Clone().Op("=").Nil(),
			)

			return []jen.Code{vAssertion, notOk, branch}
		},
		releaseFn: releaseNullFn(col, "NullInt64"),
	}
}

func uuidField(col string, dest *jen.Statement) scannerField {
	return scannerField{
		columnName: col,
		valueExpr:  jen.Qual(conditionPath, "AcquireUUID").Call(),
		assignFn: func(idx string) []jen.Code {
			vAssertion := jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
				jen.Id("values").Index(jen.Id(idx)).Assert(jen.Op("*").Qual(modelPath, uuid)),
			)
			notOk := jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad type %T"),
					jen.Id("values").Index(jen.Id(idx)),
				)),
			)
			assign := dest.Clone().Op("=").Op("*").Id("v")

			return []jen.Code{vAssertion, notOk, assign}
		},
		releaseFn: releaseUUIDFn(col),
	}
}

// releaseUUIDFn emits the body of one case in ReleaseValues for a pooled
// *model.UUID cell.
func releaseUUIDFn(col string) scannerAssignFunc {
	return func(idx string) []jen.Code {
		return []jen.Code{
			jen.If(
				jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
					jen.Id("values").Index(jen.Id(idx)).Assert(
						jen.Op("*").Qual(modelPath, uuid),
					),
				),
				jen.Id("ok"),
			).Block(
				jen.Qual(conditionPath, "ReleaseUUID").Call(jen.Id("v")),
			),
		}
	}
}

func timeField(col string, dest *jen.Statement) scannerField {
	return scannerField{
		columnName: col,
		valueExpr:  jen.Qual(conditionPath, "AcquireNullTime").Call(),
		assignFn: func(idx string) []jen.Code {
			vAssertion := jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
				jen.Id("values").Index(jen.Id(idx)).Assert(jen.Op("*").Qual("database/sql", "NullTime")),
			)
			notOk := jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad type %T"),
					jen.Id("values").Index(jen.Id(idx)),
				)),
			)
			validGuard := jen.If(jen.Id("v").Dot("Valid")).Block(
				dest.Clone().Op("=").Id("v").Dot("Time"),
			)

			return []jen.Code{vAssertion, notOk, validGuard}
		},
		releaseFn: releaseNullFn(col, "NullTime"),
	}
}

func nullableTypeField(col string, dest *jen.Statement, typeV Type) scannerField {
	typeQual := jen.Qual(typeV.Pkg().Path(), typeV.Name())

	// Use the pooled Acquire helper for the well-known nullable types.
	// gorm.DeletedAt and the sql.Null* family all have helpers in
	// condition/scanner_pool.go. Fall back to plain `new(T)` for anything
	// else (rare — non-Null custom-nullable types).
	acquireFn := poolAcquireFor(typeV)
	releaseFn := poolReleaseFor(typeV)

	valueExpr := jen.New(typeQual.Clone())
	if acquireFn != "" {
		valueExpr = jen.Qual(conditionPath, acquireFn).Call()
	}

	return scannerField{
		columnName: col,
		valueExpr:  valueExpr,
		assignFn: func(idx string) []jen.Code {
			vAssertion := jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
				jen.Id("values").Index(jen.Id(idx)).Assert(jen.Op("*").Add(typeQual.Clone())),
			)
			notOk := jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad type %T"),
					jen.Id("values").Index(jen.Id(idx)),
				)),
			)
			assign := dest.Clone().Op("=").Op("*").Id("v")

			return []jen.Code{vAssertion, notOk, assign}
		},
		releaseFn: releaseFn,
	}
}

// poolAcquireFor returns the name of the condition.AcquireX function for a
// well-known nullable type, or "" if no pool exists. The map matches
// scanner_pool.go.
func poolAcquireFor(typeV Type) string {
	switch typeV.String() {
	case "database/sql.NullBool":
		return "AcquireNullBool"
	case "database/sql.NullString":
		return "AcquireNullString"
	case "database/sql.NullInt16":
		return "AcquireNullInt16"
	case "database/sql.NullInt32":
		return "AcquireNullInt32"
	case "database/sql.NullInt64":
		return "AcquireNullInt64"
	case "database/sql.NullByte":
		return "AcquireNullByte"
	case "database/sql.NullFloat64":
		return "AcquireNullFloat64"
	case "database/sql.NullTime":
		return "AcquireNullTime"
	case "gorm.io/gorm.DeletedAt":
		return "AcquireDeletedAt"
	}

	return ""
}

func poolReleaseFor(typeV Type) scannerAssignFunc {
	releaseName := ""

	switch typeV.String() {
	case "database/sql.NullBool":
		releaseName = "ReleaseNullBool"
	case "database/sql.NullString":
		releaseName = "ReleaseNullString"
	case "database/sql.NullInt16":
		releaseName = "ReleaseNullInt16"
	case "database/sql.NullInt32":
		releaseName = "ReleaseNullInt32"
	case "database/sql.NullInt64":
		releaseName = "ReleaseNullInt64"
	case "database/sql.NullByte":
		releaseName = "ReleaseNullByte"
	case "database/sql.NullFloat64":
		releaseName = "ReleaseNullFloat64"
	case "database/sql.NullTime":
		releaseName = "ReleaseNullTime"
	case "gorm.io/gorm.DeletedAt":
		releaseName = "ReleaseDeletedAt"
	default:
		return nil
	}

	typeQual := jen.Qual(typeV.Pkg().Path(), typeV.Name())

	return func(idx string) []jen.Code {
		return []jen.Code{
			jen.If(
				jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
					jen.Id("values").Index(jen.Id(idx)).Assert(
						jen.Op("*").Add(typeQual.Clone()),
					),
				),
				jen.Id("ok"),
			).Block(
				jen.Qual(conditionPath, releaseName).Call(jen.Id("v")),
			),
		}
	}
}

// customScannerField handles user-defined types that implement sql.Scanner.
// The allocated destination is a NullableScanner wrapper so NULL columns are
// silently skipped (mirrors gorm's reflective path; without this, types
// whose Scan doesn't accept nil — e.g. MultiString — fail on NULL columns).
//
// AssignValues unwraps via the .Inner field and asserts to the concrete
// pointer to copy the value into dest.
func customScannerField(col string, dest *jen.Statement, typeV Type) scannerField {
	typeQual := jen.Qual(typeV.Pkg().Path(), typeV.Name())

	return scannerField{
		columnName: col,
		valueExpr: jen.Op("&").Qual(conditionPath, "NullableScanner").Values(jen.Dict{
			jen.Id("Inner"): jen.New(typeQual.Clone()),
		}),
		assignFn: func(idx string) []jen.Code {
			wrapAssertion := jen.List(jen.Id("w"), jen.Id("ok")).Op(":=").Add(
				jen.Id("values").Index(jen.Id(idx)).Assert(jen.Op("*").Qual(conditionPath, "NullableScanner")),
			)
			notOkWrap := jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad wrapper type %T"),
					jen.Id("values").Index(jen.Id(idx)),
				)),
			)
			innerAssertion := jen.List(jen.Id("v"), jen.Id("ok2")).Op(":=").Add(
				jen.Id("w").Dot("Inner").Assert(jen.Op("*").Add(typeQual.Clone())),
			)
			notOkInner := jen.If(jen.Op("!").Id("ok2")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("cql scanner "+col+": bad inner type %T"),
					jen.Id("w").Dot("Inner"),
				)),
			)
			assign := dest.Clone().Op("=").Op("*").Id("v")

			return []jen.Code{wrapAssertion, notOkWrap, innerAssertion, notOkInner, assign}
		},
	}
}

// emitScannerVar writes the package-level var <Model>Scanner with ScanValues
// and AssignValues closures.
func (sg ScannerGenerator) emitScannerVar(file *File, objectQual *jen.Statement, varName string, fields []scannerField) {
	// Build ScanValues body: a switch with one case per column.
	scanCases := make([]jen.Code, 0, len(fields)+1)
	for _, f := range fields {
		scanCases = append(scanCases,
			jen.Case(jen.Lit(f.columnName)).Block(
				jen.Id("values").Index(jen.Id("i")).Op("=").Add(f.valueExpr),
			),
		)
	}

	scanCases = append(scanCases,
		jen.Default().Block(
			jen.Id("values").Index(jen.Id("i")).Op("=").Qual(conditionPath, "AcquireNullSink").Call(),
		),
	)

	scanValuesFn := jen.Func().Params(
		jen.Id("columns").Index().String(),
	).Params(
		jen.Index().Any(),
		jen.Error(),
	).Block(
		jen.Id("values").Op(":=").Make(jen.Index().Any(), jen.Len(jen.Id("columns"))),
		jen.For(jen.List(jen.Id("i"), jen.Id("c")).Op(":=").Range().Id("columns")).Block(
			jen.Switch(jen.Id("c")).Block(scanCases...),
		),
		jen.Return(jen.Id("values"), jen.Nil()),
	)

	// Build AssignValues body.
	assignCases := make([]jen.Code, 0, len(fields))
	for _, f := range fields {
		assignCases = append(assignCases,
			jen.Case(jen.Lit(f.columnName)).Block(f.assignFn("i")...),
		)
	}

	assignValuesFn := jen.Func().Params(
		jen.Id("dest").Op("*").Add(objectQual.Clone()),
		jen.Id("columns").Index().String(),
		jen.Id("values").Index().Any(),
	).Error().Block(
		jen.For(jen.List(jen.Id("i"), jen.Id("c")).Op(":=").Range().Id("columns")).Block(
			jen.Switch(jen.Id("c")).Block(assignCases...),
		),
		jen.Return(jen.Nil()),
	)

	// Build ReleaseValues body — returns pooled cells. Skips fields that
	// have no releaseFn (custom Scanner types, []byte, etc.).
	releaseCases := make([]jen.Code, 0, len(fields)+1)

	for _, f := range fields {
		if f.releaseFn == nil {
			continue
		}

		releaseCases = append(releaseCases,
			jen.Case(jen.Lit(f.columnName)).Block(f.releaseFn("i")...),
		)
	}

	// Always release the NullSink default-case cells.
	releaseCases = append(releaseCases,
		jen.Default().Block(
			jen.If(
				jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Add(
					jen.Id("values").Index(jen.Id("i")).Assert(jen.Op("*").Qual(conditionPath, "NullSink")),
				),
				jen.Id("ok"),
			).Block(
				jen.Qual(conditionPath, "ReleaseNullSink").Call(jen.Id("v")),
			),
		),
	)

	releaseValuesFn := jen.Func().Params(
		jen.Id("columns").Index().String(),
		jen.Id("values").Index().Any(),
	).Block(
		jen.For(jen.List(jen.Id("i"), jen.Id("c")).Op(":=").Range().Id("columns")).Block(
			jen.Switch(jen.Id("c")).Block(releaseCases...),
		),
	)

	file.Add(
		jen.Var().Id(varName).Op("=").Op("&").Qual(conditionPath, "Scanner").Types(objectQual.Clone()).Values(
			jen.Dict{
				jen.Id("ScanValues"):    scanValuesFn,
				jen.Id("AssignValues"):  assignValuesFn,
				jen.Id("ReleaseValues"): releaseValuesFn,
			},
		),
	)
}

// emitRelationScanners emits a `var <model><Relation>JoinScanner = &condition.RelationScanner[Parent, Child]{...}`
// for each BelongsTo / HasOne relation on the model. The generated
// JoinCondition builder (emitted by conditionsGenerator.generateJoin)
// references this var by the convention name from relationScannerVarName,
// so the runtime can register a fast-scan mounter when the user calls
// .Preload() on the relation.
//
// Skipped silently when:
//   - the relation's child type isn't a CQL model
//   - the relation is HasMany (`*[]Seller`-style) — covered by HasMany
//     loaders in a future phase, not by joined-preload scanners
func (sg ScannerGenerator) emitRelationScanners(file *File, destPkg string) {
	fields, err := getFields(sg.objectType)
	if err != nil {
		return
	}

	for _, f := range fields {
		// Walk through embeddings to find any top-level relation fields.
		// Embedded structs (UUIDModel) don't carry relations.
		if f.Embedded {
			continue
		}

		sg.emitRelationScannerForField(file, destPkg, f)
	}
}

// emitRelationScannerForField inspects one field and emits the relation
// scanner if it's a BelongsTo / HasOne to another CQL model, or a HasMany
// loader if the field is a (pointer-to-)slice of CQL models.
func (sg ScannerGenerator) emitRelationScannerForField(file *File, destPkg string, field Field) {
	// HasMany: (*)[]Child or (*)[]*Child where Child is a CQL model.
	if hasManyChild, ok := hasManyChildType(field.Type.Type); ok {
		sg.emitHasManyLoader(file, destPkg, field, hasManyChild)

		return
	}

	// BelongsTo / HasOne: (*)Child where Child is a CQL model.
	isPointer := false

	ft := field.Type.Type
	if ptr, ok := ft.(*types.Pointer); ok {
		isPointer = true
		ft = ptr.Elem()
	}

	named, ok := ft.(*types.Named)
	if !ok {
		return
	}

	childType := Type{Type: named}
	if _, err := childType.CQLModelStruct(); err != nil {
		return // not a CQL model — not a relation
	}

	parentQual := jen.Qual(
		getRelativePackagePath(destPkg, sg.objectType),
		sg.objectType.Name(),
	)
	childQual := jen.Qual(
		getRelativePackagePath(destPkg, childType),
		childType.Name(),
	)

	varName := relationScannerVarName(sg.object.Name(), field.Name)
	childScannerVar := strcase.ToCamel(childType.Name()) + "Scanner"

	// Mount/SetNil bodies differ for pointer vs value relations.
	var mountBody, setNilBody []jen.Code

	if isPointer {
		// pointer: assign the child pointer; nil on reset.
		mountBody = []jen.Code{jen.Id("p").Dot(field.Name).Op("=").Id("c")}
		setNilBody = []jen.Code{jen.Id("p").Dot(field.Name).Op("=").Nil()}
	} else {
		// value: dereference the child; SetNil is a no-op since the
		// zero-valued struct already matches gorm's LEFT-JOIN-miss
		// behavior (probed in TestProbe_ValueRelationNoMatch).
		mountBody = []jen.Code{jen.Id("p").Dot(field.Name).Op("=").Op("*").Id("c")}
		setNilBody = nil
	}

	mountFn := jen.Func().Params(
		jen.Id("p").Op("*").Add(parentQual.Clone()),
		jen.Id("c").Op("*").Add(childQual.Clone()),
	).Block(mountBody...)

	values := jen.Dict{
		jen.Id("RelationField"): jen.Lit(field.Name),
		jen.Id("ChildScanner"):  jen.Id(childScannerVar),
		jen.Id("Mount"):         mountFn,
	}

	if setNilBody != nil {
		setNilFn := jen.Func().Params(
			jen.Id("p").Op("*").Add(parentQual.Clone()),
		).Block(setNilBody...)
		values[jen.Id("SetNil")] = setNilFn
	}

	file.Add(
		jen.Var().Id(varName).Op("=").Op("&").Qual(conditionPath, "RelationScanner").Types(
			parentQual.Clone(),
			childQual.Clone(),
		).Values(values),
	)
}

// relationScannerVarName is the convention used by both ScannerGenerator
// (to emit the var) and ConditionsGenerator (to reference it from the
// generated JoinCondition builder). e.g. relationScannerVarName("Phone",
// "Brand") => "phoneBrandJoinScanner".
func relationScannerVarName(modelName, relationName string) string {
	return strcase.ToCamel(modelName) + strcase.ToPascal(relationName) + "JoinScanner"
}

// hasManyLoaderVarName matches the convention used by ScannerGenerator
// (to emit the var) and ConditionsGenerator (to reference it from
// createCollection). e.g. hasManyLoaderVarName("Company", "Sellers") =>
// "companySellersHasManyLoader".
func hasManyLoaderVarName(modelName, relationName string) string {
	return strcase.ToCamel(modelName) + strcase.ToPascal(relationName) + "HasManyLoader"
}

// hasManyChildType reports whether t is a slice (possibly pointer-wrapped)
// of a CQL model — i.e. a HasMany relation field. Returns the child's
// Type when so.
func hasManyChildType(t types.Type) (Type, bool) {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	slice, ok := t.(*types.Slice)
	if !ok {
		return Type{}, false
	}

	elem := slice.Elem()
	if ptr, ok := elem.(*types.Pointer); ok {
		elem = ptr.Elem()
	}

	if _, ok := elem.(*types.Named); !ok {
		return Type{}, false
	}

	childType := Type{Type: elem}
	if _, err := childType.CQLModelStruct(); err != nil {
		return Type{}, false
	}

	return childType, true
}

// emitHasManyLoader writes a HasManyLoader[Parent, Child] var. Generated
// code wires:
//   - ParentID: extract parent's primary key (assumes embedded UUIDModel /
//     UIntModel — pulls dest.ID)
//   - ChildFK: extract child's FK to parent, with nullable-pointer
//     dereferencing for `*model.UUID` / `*model.UIntID` foreign keys
//   - Mount: assign the slice onto the parent's collection field
//     (handles `*[]Child` / `[]Child` / `*[]*Child` / `[]*Child` shapes)
//   - BuildQuery: SELECT children WHERE fk IN (parent_ids), with nested
//     conditions appended for nested preload chains
func (sg ScannerGenerator) emitHasManyLoader(file *File, destPkg string, field Field, childType Type) {
	parentQual := jen.Qual(
		getRelativePackagePath(destPkg, sg.objectType),
		sg.objectType.Name(),
	)
	childQual := jen.Qual(
		getRelativePackagePath(destPkg, childType),
		childType.Name(),
	)

	varName := hasManyLoaderVarName(sg.object.Name(), field.Name)

	// Detect parent PK type (UUID or UIntID) from the parent's embedded
	// base model. Used to type the IN-list values.
	parentIsUUID := sg.parentPKIsUUID()

	// FK column on the child Go field (e.g. "CompanyID"). gorm convention
	// is <ParentName>ID unless overridden by the foreignKey tag. cql-gen
	// already resolves this for Collection emission via
	// field.getRelatedTypeFKAttribute — we use the same value.
	fkGoField := field.getRelatedTypeFKAttribute(sg.objectType.Name())

	// IN-list builder: typed value of model.UUID or model.UIntID per id.
	var idTypedValueCall *jen.Statement

	var valueOfType *jen.Statement

	if parentIsUUID {
		idTypedValueCall = jen.Qual(conditionPath, "UUID")
		valueOfType = jen.Qual(conditionPath, "ValueOfType").Types(jen.Qual(modelPath, uuid))
	} else {
		idTypedValueCall = jen.Qual(conditionPath, "UIntID")
		valueOfType = jen.Qual(conditionPath, "ValueOfType").Types(jen.Qual(modelPath, uIntID))
	}

	// Determine whether the slice is by-value or by-pointer, and whether
	// the field is a pointer-to-slice. The Mount body needs to know:
	//   `Sellers *[]Seller`        → p.Sellers = &valuesSlice
	//   `Sellers []Seller`         →  p.Sellers = valuesSlice
	//   `Sellers *[]*Seller`       → p.Sellers = &childrenPtrs
	//   `Sellers []*Seller`        →  p.Sellers = childrenPtrs
	fieldType := field.Type.Type

	sliceIsPointer := false
	if ptr, ok := fieldType.(*types.Pointer); ok {
		sliceIsPointer = true
		fieldType = ptr.Elem()
	}

	slice, _ := fieldType.(*types.Slice)
	elemIsPointer := false

	if _, ok := slice.Elem().(*types.Pointer); ok {
		elemIsPointer = true
	}

	// Mount body.
	mountBody := []jen.Code{}
	if elemIsPointer {
		// children IS []*Child — assign directly (or take address).
		if sliceIsPointer {
			mountBody = append(mountBody,
				jen.Id("p").Dot(field.Name).Op("=").Op("&").Id("children"),
			)
		} else {
			mountBody = append(mountBody,
				jen.Id("p").Dot(field.Name).Op("=").Id("children"),
			)
		}
	} else {
		// Need to deref children []*Child into a []Child.
		mountBody = append(mountBody,
			jen.Id("values").Op(":=").Make(jen.Index().Add(childQual.Clone()), jen.Len(jen.Id("children"))),
			jen.For(jen.List(jen.Id("i"), jen.Id("c")).Op(":=").Range().Id("children")).Block(
				jen.Id("values").Index(jen.Id("i")).Op("=").Op("*").Id("c"),
			),
		)

		if sliceIsPointer {
			mountBody = append(mountBody,
				jen.Id("p").Dot(field.Name).Op("=").Op("&").Id("values"),
			)
		} else {
			mountBody = append(mountBody,
				jen.Id("p").Dot(field.Name).Op("=").Id("values"),
			)
		}
	}

	// ChildFK body — depends on whether child FK is a pointer.
	childFKBody := childFKExtractor(childType, fkGoField, parentIsUUID)

	// ParentID body — emit `return p.ID` (assumes embedded base model).
	parentIDFn := jen.Func().Params(jen.Id("p").Op("*").Add(parentQual.Clone())).
		Params(jen.Any()).Block(
		jen.Return(jen.Id("p").Dot("ID")),
	)

	childFKFn := jen.Func().Params(jen.Id("c").Op("*").Add(childQual.Clone())).
		Params(jen.Any()).Block(childFKBody...)

	mountFn := jen.Func().Params(
		jen.Id("p").Op("*").Add(parentQual.Clone()),
		jen.Id("children").Index().Op("*").Add(childQual.Clone()),
	).Block(mountBody...)

	// BuildQuery body — typed IN-list + NewQuery[Child].
	buildQueryFn := jen.Func().Params(
		jen.Id("tx").Op("*").Qual("gorm.io/gorm", "DB"),
		jen.Id("parentIDs").Index().Any(),
		jen.Id("nested").Index().Qual(conditionPath, cqlCondition).Types(childQual.Clone()),
	).Params(
		jen.Op("*").Qual(conditionPath, "Query").Types(childQual.Clone()),
		jen.Error(),
	).Block(
		jen.Id("typedIDs").Op(":=").Make(jen.Index().Add(valueOfType.Clone()), jen.Len(jen.Id("parentIDs"))),
		jen.For(jen.List(jen.Id("i"), jen.Id("id")).Op(":=").Range().Id("parentIDs")).Block(
			jen.Id("typedIDs").Index(jen.Id("i")).Op("=").Add(idTypedValueCall.Clone()).Call(
				jen.Id("id").Assert(jen.Qual(modelPath, parentPKTypeName(parentIsUUID))),
			),
		),
		jen.Id("conds").Op(":=").Append(
			jen.Index().Qual(conditionPath, cqlCondition).Types(childQual.Clone()).Values(
				jen.Id(strcase.ToPascal(childType.Name())).Dot(fkGoField).Dot("Is").Call().Dot("In").Call(jen.Id("typedIDs").Op("...")),
			),
			jen.Id("nested").Op("..."),
		),
		jen.Return(
			jen.Qual(conditionPath, "NewQuery").Types(childQual.Clone()).Call(
				jen.Id("tx"), jen.Id("conds").Op("..."),
			),
			jen.Nil(),
		),
	)

	values := jen.Dict{
		jen.Id("CollectionField"): jen.Lit(field.Name),
		jen.Id("ParentID"):        parentIDFn,
		jen.Id("ChildFK"):         childFKFn,
		jen.Id("Mount"):           mountFn,
		jen.Id("BuildQuery"):      buildQueryFn,
	}

	file.Add(
		jen.Var().Id(varName).Op("=").Op("&").Qual(conditionPath, "HasManyLoader").Types(
			parentQual.Clone(),
			childQual.Clone(),
		).Values(values),
	)
}

// parentPKTypeName returns the model.<Type> for the parent's PK.
func parentPKTypeName(isUUID bool) string {
	if isUUID {
		return uuid
	}

	return uIntID
}

// parentPKIsUUID returns true when the parent's embedded base model uses
// UUID as its primary key.
func (sg ScannerGenerator) parentPKIsUUID() bool {
	fields, err := getFields(sg.objectType)
	if err != nil {
		return false
	}

	for _, f := range fields {
		if !f.Embedded {
			continue
		}

		switch f.TypeString() {
		case modelPath + "." + uuidModel, modelPath + "." + uuidModelWithTimestamps:
			return true
		case modelPath + "." + uIntModel, modelPath + "." + uIntModelWithTimestamps:
			return false
		}
	}

	return true // default to UUID
}

// childFKExtractor builds the ChildFK function body. For a pointer FK
// (e.g. `*model.UUID`), returns nil sentinel when nil; else dereferences.
// For a value FK (e.g. `model.UUID`), returns directly.
func childFKExtractor(childType Type, fkGoField string, parentIsUUID bool) []jen.Code {
	// Inspect the child struct to find the FK field's exact type.
	childStruct, err := childType.CQLModelStruct()
	if err != nil {
		return []jen.Code{jen.Return(jen.Nil())}
	}

	for i := 0; i < childStruct.NumFields(); i++ {
		f := childStruct.Field(i)
		if f.Name() != fkGoField {
			continue
		}

		ft := f.Type()
		if ptr, ok := ft.(*types.Pointer); ok {
			// Nullable FK — guard against nil.
			_ = ptr

			var nilSentinel *jen.Statement
			if parentIsUUID {
				nilSentinel = jen.Qual(modelPath, "NilUUID")
			} else {
				nilSentinel = jen.Qual(modelPath, "NilUIntID")
			}

			return []jen.Code{
				jen.If(jen.Id("c").Dot(fkGoField).Op("==").Nil()).Block(
					jen.Return(nilSentinel),
				),
				jen.Return(jen.Op("*").Id("c").Dot(fkGoField)),
			}
		}

		// Value FK — non-nullable, return directly.
		return []jen.Code{jen.Return(jen.Id("c").Dot(fkGoField))}
	}

	return []jen.Code{jen.Return(jen.Nil())}
}

// emitInitRewiring writes an init() that overwrites each Field on the model's
// conditions struct with a scanner-aware constructor call. This piggybacks on
// the existing conditionsGenerator output: we look up the var name (the model
// name) and reassign each field. Re-uses the same NewXxxField constructor
// resolution that condition.go already performs.
//
// Crucially, ConditionsGenerator and ScannerGenerator are decoupled: this
// init() runs after both vars are declared and rebinds Brand.X = NewField(...,
// scannerPtr) so all subsequent uses see the scanner.
func (sg ScannerGenerator) emitInitRewiring(file *File, scannerVar string) {
	fields, err := getFields(sg.objectType)
	if err != nil {
		return
	}

	// Build the flattened list with NamePrefix / ColumnPrefix tracking,
	// mirroring conditionsGenerator's ForEachField → generateForEmbeddedField.
	flatConditions := []condBinding{}

	for _, f := range fields {
		flatConditions = append(flatConditions, sg.flattenConditionBindings(file.destPkg, f, "", "")...)
	}

	stmts := []jen.Code{}

	for _, b := range flatConditions {
		stmts = append(stmts, b.rebindStmt(file.destPkg, sg.objectType, scannerVar))
	}

	// Wire parent scanner onto each HasMany Collection so a query whose
	// only condition is Collection.Preload() can still take the fast
	// path (the resolveScanner walk needs SOMETHING to ask for a scanner;
	// the Collection's wrapping condition is the only candidate).
	for _, f := range fields {
		if f.Embedded {
			continue
		}

		if _, ok := hasManyChildType(f.Type.Type); !ok {
			continue
		}

		stmts = append(stmts,
			jen.Id(sg.objectType.Name()).Dot(f.Name).Op("=").
				Id(sg.objectType.Name()).Dot(f.Name).Dot("WithParentScanner").Call(jen.Id(scannerVar)),
		)
	}

	file.Add(jen.Func().Id("init").Params().Block(stmts...))
}

// condBinding mirrors how Condition is built so we can regenerate the
// constructor call with a scanner argument appended.
type condBinding struct {
	field        Field
	namePrefix   string
	columnPrefix string
	param        *JenParam
}

func (sg ScannerGenerator) flattenConditionBindings(destPkg string, field Field, namePrefix, columnPrefix string) []condBinding {
	if field.Embedded {
		embeddedStruct, ok := field.Type.Underlying().(*types.Struct)
		if !ok {
			return nil
		}

		subFields, err := getStructFields(embeddedStruct)
		if err != nil {
			return nil
		}

		if !isBaseModel(field.TypeString()) {
			newNamePrefix := namePrefix + field.Name
			newColumnPrefix := columnPrefix + field.Tags.getEmbeddedPrefix()

			out := []condBinding{}

			for _, sub := range subFields {
				out = append(out, sg.flattenConditionBindings(destPkg, sub, newNamePrefix, newColumnPrefix)...)
			}

			return out
		}

		// Base model — sub fields carry no prefix.
		out := []condBinding{}
		for _, sub := range subFields {
			out = append(out, sg.flattenConditionBindings(destPkg, sub, namePrefix, columnPrefix)...)
		}

		return out
	}

	// If the field is a pointer to anything, unwrap so the resulting
	// Field carries wasPointer=true. That makes IsNullable() return true
	// and pickConstructor select the Nullable* variant, matching what
	// conditionsGenerator emits for the same field. Without this, the
	// init() rewiring assigns a non-nullable constructor result into a
	// NullableXField slot and the build fails.
	unwrapped := unwrapPointerLeaf(field)

	// Relation fields (cql model named types) are NOT emitted as Fields
	// on the conditions struct — they're handled via a builder method
	// (e.g. owned.Owner(...)). Skip them here so init() doesn't try to
	// rewire a non-existent struct field.
	if _, err := unwrapped.Type.CQLModelStruct(); err == nil {
		return nil
	}

	// Slices of cql models (`Sellers *[]Seller`) are HasMany collections,
	// also handled via a separate Collection (not a Field). Skip.
	if sliceOfCQLModel(field.Type.Type) {
		return nil
	}

	param := NewJenParam()
	classifyParam(unwrapped, param)

	return []condBinding{{
		field:        unwrapped,
		namePrefix:   namePrefix,
		columnPrefix: columnPrefix,
		param:        param,
	}}
}

// sliceOfCQLModel reports whether t is a slice (or pointer-to-slice) whose
// element type is itself a cql model struct. Used to skip HasMany relation
// fields during init()-rewiring.
func sliceOfCQLModel(t types.Type) bool {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	slice, ok := t.(*types.Slice)
	if !ok {
		return false
	}

	elem := slice.Elem()
	if ptr, ok := elem.(*types.Pointer); ok {
		elem = ptr.Elem()
	}

	if _, ok := elem.(*types.Named); !ok {
		return false
	}

	_, err := (Type{Type: elem}).CQLModelStruct()

	return err == nil
}

// unwrapPointerLeaf walks *T -> T at the leaf, marking wasPointer so
// downstream nullability checks behave as the conditions generator does.
func unwrapPointerLeaf(field Field) Field {
	if ptr, ok := field.Type.Type.(*types.Pointer); ok {
		return field.ChangeType(ptr.Elem(), true)
	}

	return field
}

// rebindStmt emits one assignment: <Model>.<Path> = condition.NewXField(...).
func (b condBinding) rebindStmt(destPkg string, objectType Type, scannerVar string) jen.Code {
	objectQual := jen.Qual(
		getRelativePackagePath(destPkg, objectType),
		objectType.Name(),
	)

	fieldName := b.namePrefix + b.field.Name

	fieldColumn := jen.Lit("")
	if col := b.field.getColumnName(); col != "" {
		fieldColumn = jen.Lit(col)
	}

	fieldColumnPrefix := jen.Lit(b.columnPrefix)

	var newFieldQual *jen.Statement

	switch {
	case b.param.isString:
		newFieldQual = pickConstructor(b.field, cqlNullableStringField, cqlStringField, cqlStringField).Types(objectQual)
	case b.param.isBool:
		newFieldQual = pickConstructor(b.field, cqlNullableBoolField, cqlBoolField, cqlBoolField).Types(objectQual)
	case b.param.isNumeric:
		newFieldQual = pickConstructor(b.field, cqlNullableNumericField, cqlNumericField, cqlNumericField).Types(
			objectQual,
			b.param.GenericType(),
		)
	default:
		newFieldQual = pickConstructor(b.field, cqlNullableField, cqlUpdatableField, cqlField).Types(
			objectQual,
			b.param.GenericType(),
		)
	}

	return jen.Id(objectType.Name()).Dot(fieldName).Op("=").Add(
		newFieldQual.Call(
			jen.Lit(b.field.Name),
			fieldColumn,
			fieldColumnPrefix,
			jen.Id(scannerVar),
		),
	)
}

func pickConstructor(field Field, nullableType, updatableType, notNullableType string) *jen.Statement {
	switch {
	case field.IsNullable():
		return jen.Qual(conditionPath, cqlNewField+nullableType)
	case field.IsUpdatable():
		return jen.Qual(conditionPath, cqlNewField+updatableType)
	default:
		return jen.Qual(conditionPath, cqlNewField+notNullableType)
	}
}

// classifyParam mirrors condition.go's generate() switch enough to set the
// isString / isBool / isNumeric / generic-type flags so we pick the right
// NewXxxField constructor.
func classifyParam(field Field, param *JenParam) {
	switch ft := field.GetType().(type) {
	case *types.Basic:
		param.ToBasicKind(ft)
	case *types.Pointer:
		inner := Field{Name: field.Name, Type: Type{Type: ft.Elem(), wasPointer: true}, Tags: field.Tags}
		classifyParam(inner, param)
	case *types.Named:
		switch {
		case field.Type.IsSQLNullableType():
			param.SQLToBasicType(field.Type)
		case field.Type.IsGormCustomType() || field.TypeString() == "time.Time" || field.IsModelID():
			param.ToCustomType("", field.Type)
		}
	case *types.Slice:
		if elem, ok := ft.Elem().(*types.Basic); ok {
			param.ToSlice()
			param.ToBasicKind(elem)
		}
	}
}

func fileNameForScanner(dir, objectName string) string {
	return dir + "/" + strcase.ToSnake(objectName) + "_scanner.go"
}
