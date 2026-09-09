package tests

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/viper"
	"gotest.tools/assert"

	"github.com/FrancoLiberali/cql/cql-gen/cmd"
)

const chunkSize = 100000

func TestUIntModel(t *testing.T) {
	doTest(t, "./uintmodel", []Comparison{
		{
			Have:            "uint_model_conditions.go",
			Expected:        "./results/uintmodel.go",
			ScannerHave:     "uint_model_scanner.go",
			ScannerExpected: "./results/uintmodel_scanner.go",
		},
	})
	CheckFileNotExists(t, "./uintmodel/cql.go")
}

func TestUIntModelWithTimestamp(t *testing.T) {
	doTest(t, "./uintmodelwithtimestamp", []Comparison{
		{
			Have:            "u_int_model_with_timestamp_conditions.go",
			Expected:        "./results/uintmodelwithtimestamp.go",
			ScannerHave:     "u_int_model_with_timestamp_scanner.go",
			ScannerExpected: "./results/uintmodelwithtimestamp_scanner.go",
		},
	})
	CheckFileNotExists(t, "./uintmodelwithtimestamp/cql.go")
}

func TestUIntModelWithTimestamps(t *testing.T) {
	doTest(t, "./uintmodelwithtimestamps", []Comparison{
		{
			Have:            "u_int_model_with_timestamps_conditions.go",
			Expected:        "./results/uintmodelwithtimestamps.go",
			ScannerHave:     "u_int_model_with_timestamps_scanner.go",
			ScannerExpected: "./results/uintmodelwithtimestamps_scanner.go",
		},
	})
	CheckFileNotExists(t, "./uintmodelwithtimestamps/cql.go")
}

func TestUUIDModel(t *testing.T) {
	doTest(t, "./uuidmodel", []Comparison{
		{
			Have:            "uuid_model_conditions.go",
			Expected:        "./results/uuidmodel.go",
			ScannerHave:     "uuid_model_scanner.go",
			ScannerExpected: "./results/uuidmodel_scanner.go",
		},
	})
	CheckFileNotExists(t, "./uuidmodel/cql.go")
}

func TestUUIDModelWithTimestamp(t *testing.T) {
	doTest(t, "./uuidmodelwithtimestamp", []Comparison{
		{
			Have:            "uuid_model_with_timestamp_conditions.go",
			Expected:        "./results/uuidmodelwithtimestamp.go",
			ScannerHave:     "uuid_model_with_timestamp_scanner.go",
			ScannerExpected: "./results/uuidmodelwithtimestamp_scanner.go",
		},
	})
	CheckFileNotExists(t, "./uuidmodelwithtimestamp/cql.go")
}

func TestUUIDModelWithTimestamps(t *testing.T) {
	doTest(t, "./uuidmodelwithtimestamps", []Comparison{
		{
			Have:            "uuid_model_with_timestamps_conditions.go",
			Expected:        "./results/uuidmodelwithtimestamps.go",
			ScannerHave:     "uuid_model_with_timestamps_scanner.go",
			ScannerExpected: "./results/uuidmodelwithtimestamps_scanner.go",
		},
	})
	CheckFileNotExists(t, "./uuidmodelwithtimestamps/cql.go")
}

// BasicTypes is the happy-path fixture for every Go basic type that
// round-trips through SQL via gorm. Both conditions and scanner are emitted
// and compared against goldens.
func TestBasicTypes(t *testing.T) {
	doTest(t, "./basictypes", []Comparison{
		{
			Have:            "basic_types_conditions.go",
			Expected:        "./results/basictypes.go",
			ScannerHave:     "basic_types_scanner.go",
			ScannerExpected: "./results/basictypes_scanner.go",
		},
	})
	CheckFileNotExists(t, "./basictypes/cql.go")
}

// UnsupportedBasicTypes is the negative-path fixture: a model containing
// only Go basic types that no SQL database can persist (uintptr, complex*).
// Verified by the gorm probe in /condition/scanner_gorm_probe_test.go that
// gorm itself rejects these at AutoMigrate, so cql-gen fails loudly at
// codegen rather than letting the error surface at first DB call.
//
// Users who genuinely need complex-shaped columns can define a wrapper
// implementing sql.Scanner + driver.Valuer (see TestCustomType) — that path
// bypasses this check.
func TestUnsupportedBasicTypes(t *testing.T) {
	defer func() {
		// Scanner validation runs before conditions are written, so no
		// stale file should be on disk — but defensively clean up if
		// behavior ever changes.
		RemoveFile("unsupported_basic_types_conditions.go")

		r := recover()
		if r == nil {
			t.Fatal("expected cql-gen to panic on UnsupportedBasicTypes")
		}

		err, _ := r.(error)
		if err == nil {
			t.Fatalf("expected error panic, got %T: %v", r, r)
		}

		// The first field in the model is UIntptr, so the panic should
		// name it (panics on the first offending field encountered).
		const want = "field UIntptr"
		if !containsSub(err.Error(), want) {
			t.Fatalf("panic message %q should contain %q", err.Error(), want)
		}
	}()

	doTest(t, "./unsupportedbasictypes", nil)
}

// containsSub is strings.Contains spelled out to avoid reshuffling this
// file's import block. Used by TestUnsupportedBasicTypes' panic check.
func containsSub(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}

	return false
}

func TestBasicPointers(t *testing.T) {
	doTest(t, "./basicpointers", []Comparison{
		{Have: "basic_pointers_conditions.go", Expected: "./results/basicpointers.go"},
	})
	CheckFileNotExists(t, "./basicpointers/cql.go")
}

func TestBasicSlices(t *testing.T) {
	doTest(t, "./basicslices", []Comparison{
		{Have: "basic_slices_conditions.go", Expected: "./results/basicslices.go"},
	})
	CheckFileNotExists(t, "./basicslices/cql.go")
}

func TestBasicSlicesPointer(t *testing.T) {
	doTest(t, "./basicslicespointer", []Comparison{
		{Have: "basic_slices_pointer_conditions.go", Expected: "./results/basicslicespointer.go"},
	})
	CheckFileNotExists(t, "./basicslicespointer/cql.go")
}

// TestGoEmbedded covers the negative case for Go-style anonymous embeds:
// both top-level Int and ToBeEmbedded.Int resolve to column "int". cql-gen's
// scanner generator refuses to emit because the duplicate switch case
// wouldn't compile AND gorm would reject the table at AutoMigrate. The
// conditions file is still emitted (cql-gen historically accepts ambiguous
// models there).
func TestGoEmbedded(t *testing.T) {
	defer func() {
		RemoveFile("go_embedded_conditions.go")

		r := recover()
		if r == nil {
			t.Fatal("expected cql-gen to panic on GoEmbedded (column \"int\" collision)")
		}

		err, _ := r.(error)
		if err == nil {
			t.Fatalf("expected error panic, got %T: %v", r, r)
		}

		const want = "column \"int\""
		if !containsSub(err.Error(), want) {
			t.Fatalf("panic message %q should contain %q", err.Error(), want)
		}
	}()

	doTest(t, "./goembedded", []Comparison{
		{Have: "go_embedded_conditions.go", Expected: "./results/goembedded.go"},
	})
}

// TestGormEmbedded covers the negative case for gorm:"embedded": two fields
// (top-level Int and GormEmbeddedNoPrefix.Int with no prefix) resolve to the
// same column "int". cql-gen now refuses to generate the scanner for such a
// model because (a) the duplicate case wouldn't compile and (b) gorm itself
// would reject the table at AutoMigrate. The conditions file is still
// emitted (cql-gen historically accepts ambiguous models there).
func TestGormEmbedded(t *testing.T) {
	defer func() {
		RemoveFile("gorm_embedded_conditions.go")

		r := recover()
		if r == nil {
			t.Fatal("expected cql-gen to panic on GormEmbedded (column \"int\" collision)")
		}

		err, _ := r.(error)
		if err == nil {
			t.Fatalf("expected error panic, got %T: %v", r, r)
		}

		const want = "column \"int\""
		if !containsSub(err.Error(), want) {
			t.Fatalf("panic message %q should contain %q", err.Error(), want)
		}
	}()

	doTest(t, "./gormembedded", []Comparison{
		{Have: "gorm_embedded_conditions.go", Expected: "./results/gormembedded.go"},
	})
}

// TestGormEmbeddedScan exercises the happy path for gorm:"embedded": two
// embedded fields use distinct prefixes so no column collides, and the
// scanner emits assignments through the right Go access path
// (dest.Foo.Int / dest.Bar.Int / dest.Top). Regression for the bug where
// embedded leaves were misrouted to dest.Int.
func TestGormEmbeddedScan(t *testing.T) {
	doTest(t, "./gormembeddedscan", []Comparison{
		{
			Have:            "gorm_embedded_scan_conditions.go",
			Expected:        "./results/gormembeddedscan.go",
			ScannerHave:     "gorm_embedded_scan_scanner.go",
			ScannerExpected: "./results/gormembeddedscan_scanner.go",
		},
	})
	CheckFileNotExists(t, "./gormembeddedscan/cql.go")
}

func TestCustomType(t *testing.T) {
	doTest(t, "./customtype", []Comparison{
		{
			Have:            "custom_type_conditions.go",
			Expected:        "./results/customtype.go",
			ScannerHave:     "custom_type_scanner.go",
			ScannerExpected: "./results/customtype_scanner.go",
		},
	})
	CheckFileNotExists(t, "./customtype/cql.go")
}

func TestColumnDefinition(t *testing.T) {
	doTest(t, "./columndefinition", []Comparison{
		{
			Have:            "column_definition_conditions.go",
			Expected:        "./results/columndefinition.go",
			ScannerHave:     "column_definition_scanner.go",
			ScannerExpected: "./results/columndefinition_scanner.go",
		},
	})
	CheckFileNotExists(t, "./columndefinition/cql.go")
}

func TestNullableTypes(t *testing.T) {
	doTest(t, "./nullabletypes", []Comparison{
		{
			Have:            "nullable_types_conditions.go",
			Expected:        "./results/nullabletypes.go",
			ScannerHave:     "nullable_types_scanner.go",
			ScannerExpected: "./results/nullabletypes_scanner.go",
		},
	})
	CheckFileNotExists(t, "./nullabletypes/cql.go")
}

func TestBelongsTo(t *testing.T) {
	doTest(t, "./belongsto", []Comparison{
		{
			Have:            "owner_conditions.go",
			Expected:        "./results/belongsto_owner.go",
			ScannerHave:     "owner_scanner.go",
			ScannerExpected: "./results/belongsto_owner_scanner.go",
		},
		{
			Have:            "owned_conditions.go",
			Expected:        "./results/belongsto_owned.go",
			ScannerHave:     "owned_scanner.go",
			ScannerExpected: "./results/belongsto_owned_scanner.go",
		},
		{Have: "./belongsto/cql.go", Expected: "./belongsto/cql_result.go"},
	})
}

func TestHasOne(t *testing.T) {
	doTest(t, "./hasone", []Comparison{
		{
			Have:            "country_conditions.go",
			Expected:        "./results/hasone_country.go",
			ScannerHave:     "country_scanner.go",
			ScannerExpected: "./results/hasone_country_scanner.go",
		},
		{
			Have:            "city_conditions.go",
			Expected:        "./results/hasone_city.go",
			ScannerHave:     "city_scanner.go",
			ScannerExpected: "./results/hasone_city_scanner.go",
		},
		{Have: "./hasone/cql.go", Expected: "./hasone/cql_result.go"},
	})
}

func TestHasMany(t *testing.T) {
	doTest(t, "./hasmany", []Comparison{
		{
			Have:            "company_conditions.go",
			Expected:        "./results/hasmany_company.go",
			ScannerHave:     "company_scanner.go",
			ScannerExpected: "./results/hasmany_company_scanner.go",
		},
		{
			Have:            "seller_conditions.go",
			Expected:        "./results/hasmany_seller.go",
			ScannerHave:     "seller_scanner.go",
			ScannerExpected: "./results/hasmany_seller_scanner.go",
		},
		{Have: "./hasmany/cql.go", Expected: "./hasmany/cql_result.go"},
	})
}

// TestHasManyUintParent covers the has-many loader generator branches specific
// to a UInt-keyed parent (parentPKIsUUID=false, parentPKTypeName -> UIntID,
// childFKExtractor -> NilUIntID) that the UUID-parent fixtures don't reach.
func TestHasManyUintParent(t *testing.T) {
	doTest(t, "./hasmanyuint", []Comparison{
		{
			Have:            "company_uint_conditions.go",
			Expected:        "./results/hasmanyuint_company.go",
			ScannerHave:     "company_uint_scanner.go",
			ScannerExpected: "./results/hasmanyuint_company_scanner.go",
		},
		{
			Have:            "seller_uint_conditions.go",
			Expected:        "./results/hasmanyuint_seller.go",
			ScannerHave:     "seller_uint_scanner.go",
			ScannerExpected: "./results/hasmanyuint_seller_scanner.go",
		},
	})
}

func TestHasManyWithPointers(t *testing.T) {
	doTest(t, "./hasmanywithpointers", []Comparison{
		{
			Have:            "company_with_pointers_conditions.go",
			Expected:        "./results/hasmanywithpointers_company.go",
			ScannerHave:     "company_with_pointers_scanner.go",
			ScannerExpected: "./results/hasmanywithpointers_company_scanner.go",
		},
		{
			Have:            "seller_in_pointers_conditions.go",
			Expected:        "./results/hasmanywithpointers_seller.go",
			ScannerHave:     "seller_in_pointers_scanner.go",
			ScannerExpected: "./results/hasmanywithpointers_seller_scanner.go",
		},
		{Have: "./hasmanywithpointers/cql.go", Expected: "./hasmanywithpointers/cql_result.go"},
	})
}

func TestSelfReferential(t *testing.T) {
	doTest(t, "./selfreferential", []Comparison{
		{
			Have:            "employee_conditions.go",
			Expected:        "./results/selfreferential.go",
			ScannerHave:     "employee_scanner.go",
			ScannerExpected: "./results/selfreferential_scanner.go",
		},
		{Have: "./selfreferential/cql.go", Expected: "./selfreferential/cql_result.go"},
	})
}

func TestMultiplePackage(t *testing.T) {
	doTest(t, "./multiplepackage/package1", []Comparison{
		{
			Have:            "package1_conditions.go",
			Expected:        "./results/multiplepackage_package1.go",
			ScannerHave:     "package1_scanner.go",
			ScannerExpected: "./results/multiplepackage_package1_scanner.go",
		},
		{Have: "./multiplepackage/package1/cql.go", Expected: "./multiplepackage/package1/cql_result.go"},
	})
	doTest(t, "./multiplepackage/package2", []Comparison{
		{
			Have:            "package2_conditions.go",
			Expected:        "./results/multiplepackage_package2.go",
			ScannerHave:     "package2_scanner.go",
			ScannerExpected: "./results/multiplepackage_package2_scanner.go",
		},
	})
}

func TestOverrideForeignKey(t *testing.T) {
	doTest(t, "./overrideforeignkey", []Comparison{
		{
			Have:            "bicycle_conditions.go",
			Expected:        "./results/overrideforeignkey_bicycle.go",
			ScannerHave:     "bicycle_scanner.go",
			ScannerExpected: "./results/overrideforeignkey_bicycle_scanner.go",
		},
		{
			Have:            "person_conditions.go",
			Expected:        "./results/overrideforeignkey_person.go",
			ScannerHave:     "person_scanner.go",
			ScannerExpected: "./results/overrideforeignkey_person_scanner.go",
		},
		{Have: "./overrideforeignkey/cql.go", Expected: "./overrideforeignkey/cql_result.go"},
	})
}

func TestOverrideReferences(t *testing.T) {
	doTest(t, "./overridereferences", []Comparison{
		{
			Have:            "phone_conditions.go",
			Expected:        "./results/overridereferences_phone.go",
			ScannerHave:     "phone_scanner.go",
			ScannerExpected: "./results/overridereferences_phone_scanner.go",
		},
		{
			Have:            "brand_conditions.go",
			Expected:        "./results/overridereferences_brand.go",
			ScannerHave:     "brand_scanner.go",
			ScannerExpected: "./results/overridereferences_brand_scanner.go",
		},
		{Have: "./overridereferences/cql.go", Expected: "./overridereferences/cql_result.go"},
	})
}

func TestOverrideForeignKeyInverse(t *testing.T) {
	doTest(t, "./overrideforeignkeyinverse", []Comparison{
		{
			Have:            "user_conditions.go",
			Expected:        "./results/overrideforeignkeyinverse_user.go",
			ScannerHave:     "user_scanner.go",
			ScannerExpected: "./results/overrideforeignkeyinverse_user_scanner.go",
		},
		{
			Have:            "credit_card_conditions.go",
			Expected:        "./results/overrideforeignkeyinverse_credit_card.go",
			ScannerHave:     "credit_card_scanner.go",
			ScannerExpected: "./results/overrideforeignkeyinverse_credit_card_scanner.go",
		},
		{Have: "./overrideforeignkeyinverse/cql.go", Expected: "./overrideforeignkeyinverse/cql_result.go"},
	})
}

func TestOverrideReferencesInverse(t *testing.T) {
	doTest(t, "./overridereferencesinverse", []Comparison{
		{
			Have:            "computer_conditions.go",
			Expected:        "./results/overridereferencesinverse_computer.go",
			ScannerHave:     "computer_scanner.go",
			ScannerExpected: "./results/overridereferencesinverse_computer_scanner.go",
		},
		{
			Have:            "processor_conditions.go",
			Expected:        "./results/overridereferencesinverse_processor.go",
			ScannerHave:     "processor_scanner.go",
			ScannerExpected: "./results/overridereferencesinverse_processor_scanner.go",
		},
		{Have: "./overridereferencesinverse/cql.go", Expected: "./overridereferencesinverse/cql_result.go"},
	})
}

// TestNamedScalarIsFastScanned asserts a named scalar (`type Color int`, no
// custom Scan/Value) takes the fast path: cql-gen emits both conditions (typed
// by the named type) and a scanner covering the value and pointer variants.
func TestNamedScalarIsFastScanned(t *testing.T) {
	doTest(t, "./namedscalar", []Comparison{
		{
			Have:            "with_named_scalar_conditions.go",
			Expected:        "./results/namedscalar.go",
			ScannerHave:     "with_named_scalar_scanner.go",
			ScannerExpected: "./results/namedscalar_scanner.go",
		},
	})
	CheckFileNotExists(t, "./namedscalar/cql.go")
}

// TestScannerTypesKitchenSink exercises the fast-scan generator paths the
// narrower fixtures miss: a []byte column, time.Time value + pointer, and
// database/sql nullable wrappers used directly as fields.
func TestScannerTypesKitchenSink(t *testing.T) {
	doTest(t, "./scannertypes", []Comparison{
		{
			Have:            "scanner_types_conditions.go",
			Expected:        "./results/scannertypes.go",
			ScannerHave:     "scanner_types_scanner.go",
			ScannerExpected: "./results/scannertypes_scanner.go",
		},
	})
	CheckFileNotExists(t, "./scannertypes/cql.go")
}

// TestUnsupportedColumnFallsBackToGorm asserts the safety net: a model with a
// column the fast scanner genuinely can't classify (a []string via gorm's json
// serializer) gets conditions but NO scanner, so its queries fall back to gorm
// instead of silently dropping the column.
func TestUnsupportedColumnFallsBackToGorm(t *testing.T) {
	doTest(t, "./unsupportedcolumn", []Comparison{
		{
			Have:            "with_unsupported_column_conditions.go",
			Expected:        "./results/unsupportedcolumn.go",
			ScannerHave:     "with_unsupported_column_scanner.go",
			ExpectNoScanner: true,
		},
	})
	CheckFileNotExists(t, "./unsupportedcolumn/cql.go")
}

type Comparison struct {
	Have     string
	Expected string

	// ScannerHave is the path of the generated scanner file (e.g.
	// "uint_model_scanner.go"). When set, the test asserts on the scanner:
	//   - if ScannerExpected is non-empty, files must match
	//   - if ScannerExpected is empty AND ExpectNoScanner is true, the
	//     scanner file must NOT exist (the model has fields the scanner
	//     generator can't yet support, so it silently skipped emission).
	// Leave both ScannerHave and ExpectNoScanner unset to opt out of
	// scanner-level assertions for this comparison.
	ScannerHave     string
	ScannerExpected string
	ExpectNoScanner bool
}

func doTest(t *testing.T, sourcePkg string, comparisons []Comparison) {
	viper.Set(cmd.DestPackageKey, "conditions")

	defer cleanupGeneratedScanners(comparisons)

	cmd.GenerateConditions(nil, []string{sourcePkg})

	for _, comparison := range comparisons {
		checkFilesEqual(t, comparison.Have, comparison.Expected)

		switch {
		case comparison.ScannerHave != "" && comparison.ScannerExpected != "":
			checkFilesEqual(t, comparison.ScannerHave, comparison.ScannerExpected)
		case comparison.ScannerHave != "" && comparison.ExpectNoScanner:
			CheckFileNotExists(t, comparison.ScannerHave)
		}
	}
}

// cleanupGeneratedScanners removes per-model scanner files cql-gen wrote that
// weren't explicitly compared. checkFilesEqual deletes the Have file after
// success, but the scanner file is sibling-output that may not be claimed by
// the test (e.g. relation tests like TestBelongsTo). Leaving them behind
// breaks the next test run because the directory holds the tests package and
// a stray "package conditions" file conflicts.
func cleanupGeneratedScanners(comparisons []Comparison) {
	for _, c := range comparisons {
		if c.ScannerHave != "" {
			RemoveFile(c.ScannerHave) // safe if file doesn't exist
			continue
		}

		// Derive the scanner filename from the conditions filename:
		// foo_conditions.go -> foo_scanner.go
		const condSuffix = "_conditions.go"
		if len(c.Have) > len(condSuffix) && c.Have[len(c.Have)-len(condSuffix):] == condSuffix {
			scanner := c.Have[:len(c.Have)-len(condSuffix)] + "_scanner.go"
			RemoveFile(scanner)
		}
	}
}

const HEADER_SIZE = 49 // Code generated by cql-gen v0.0.10, DO NOT EDIT

func checkFilesEqual(t *testing.T, file1, file2 string) {
	// Golden-update mode: overwrite the expected file (file2) with the freshly
	// generated one (file1). Run `UPDATE_GOLDEN=1 go test ./...` after changing
	// the generator, then re-run without it to confirm.
	if os.Getenv("UPDATE_GOLDEN") != "" {
		if data, err := os.ReadFile(file1); err == nil {
			if err := os.WriteFile(file2, data, 0o600); err != nil {
				t.Fatal(err)
			}
		}

		RemoveFile(file1)

		return
	}

	stat1 := CheckFileExists(t, file1)
	stat2 := CheckFileExists(t, file2)

	// do inputs at least have the same size?
	assert.Equal(t, stat1.Size(), stat2.Size(), "File lens are not equal")

	// long way: compare contents
	f1, err := os.Open(file1)
	if err != nil {
		t.Error(err)
	}
	defer f1.Close()
	f1.Seek(HEADER_SIZE, 0) // ignore header

	f2, err := os.Open(file2)
	if err != nil {
		t.Error(err)
	}
	defer f2.Close()
	f2.Seek(HEADER_SIZE, 0) // ignore header

	b1 := make([]byte, chunkSize)
	b2 := make([]byte, chunkSize)

	for {
		n1, err1 := io.ReadFull(f1, b1)
		n2, err2 := io.ReadFull(f2, b2)

		assert.Assert(t, bytes.Equal(b1[:n1], b2[:n2]))

		if (err1 == io.EOF && err2 == io.EOF) || (err1 == io.ErrUnexpectedEOF && err2 == io.ErrUnexpectedEOF) {
			break
		}

		// some other error, like a dropped network connection or a bad transfer
		if err1 != nil {
			t.Error(err1)
		}

		if err2 != nil {
			t.Error(err2)
		}
	}

	RemoveFile(file1)
}
