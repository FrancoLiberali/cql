package conditions

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/viper"
	"gotest.tools/assert"

	"github.com/FrancoLiberali/cql/cql-cli/cmd/testutils"
)

const chunkSize = 100000

func TestUIntModel(t *testing.T) {
	doTest(t, "./tests/uintmodel", []Comparison{
		{Have: "uint_model_conditions.go", Expected: "./tests/results/uintmodel.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/uintmodel/cql.go")
}

func TestUUIDModel(t *testing.T) {
	doTest(t, "./tests/uuidmodel", []Comparison{
		{Have: "uuid_model_conditions.go", Expected: "./tests/results/uuidmodel.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/uuidmodel/cql.go")
}

func TestBasicTypes(t *testing.T) {
	doTest(t, "./tests/basictypes", []Comparison{
		{Have: "basic_types_conditions.go", Expected: "./tests/results/basictypes.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/basictypes/cql.go")
}

func TestBasicPointers(t *testing.T) {
	doTest(t, "./tests/basicpointers", []Comparison{
		{Have: "basic_pointers_conditions.go", Expected: "./tests/results/basicpointers.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/basicpointers/cql.go")
}

func TestBasicSlices(t *testing.T) {
	doTest(t, "./tests/basicslices", []Comparison{
		{Have: "basic_slices_conditions.go", Expected: "./tests/results/basicslices.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/basicslices/cql.go")
}

func TestBasicSlicesPointer(t *testing.T) {
	doTest(t, "./tests/basicslicespointer", []Comparison{
		{Have: "basic_slices_pointer_conditions.go", Expected: "./tests/results/basicslicespointer.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/basicslicespointer/cql.go")
}

func TestGoEmbedded(t *testing.T) {
	doTest(t, "./tests/goembedded", []Comparison{
		{Have: "go_embedded_conditions.go", Expected: "./tests/results/goembedded.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/goembedded/cql.go")
}

func TestGormEmbedded(t *testing.T) {
	doTest(t, "./tests/gormembedded", []Comparison{
		{Have: "gorm_embedded_conditions.go", Expected: "./tests/results/gormembedded.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/gormembedded/cql.go")
}

func TestCustomType(t *testing.T) {
	doTest(t, "./tests/customtype", []Comparison{
		{Have: "custom_type_conditions.go", Expected: "./tests/results/customtype.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/customtype/cql.go")
}

func TestColumnDefinition(t *testing.T) {
	doTest(t, "./tests/columndefinition", []Comparison{
		{Have: "column_definition_conditions.go", Expected: "./tests/results/columndefinition.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/columndefinition/cql.go")
}

func TestNullableTypes(t *testing.T) {
	doTest(t, "./tests/nullabletypes", []Comparison{
		{Have: "nullable_types_conditions.go", Expected: "./tests/results/nullabletypes.go"},
	})
	testutils.CheckFileNotExists(t, "./tests/nullabletypes/cql.go")
}

func TestBelongsTo(t *testing.T) {
	doTest(t, "./tests/belongsto", []Comparison{
		{Have: "owner_conditions.go", Expected: "./tests/results/belongsto_owner.go"},
		{Have: "owned_conditions.go", Expected: "./tests/results/belongsto_owned.go"},
		{Have: "./tests/belongsto/cql.go", Expected: "./tests/belongsto/cql_result.go"},
	})
}

func TestHasOne(t *testing.T) {
	doTest(t, "./tests/hasone", []Comparison{
		{Have: "country_conditions.go", Expected: "./tests/results/hasone_country.go"},
		{Have: "city_conditions.go", Expected: "./tests/results/hasone_city.go"},
		{Have: "./tests/hasone/cql.go", Expected: "./tests/hasone/cql_result.go"},
	})
}

func TestHasMany(t *testing.T) {
	doTest(t, "./tests/hasmany", []Comparison{
		{Have: "company_conditions.go", Expected: "./tests/results/hasmany_company.go"},
		{Have: "seller_conditions.go", Expected: "./tests/results/hasmany_seller.go"},
		{Have: "./tests/hasmany/cql.go", Expected: "./tests/hasmany/cql_result.go"},
	})
}

func TestHasManyWithPointers(t *testing.T) {
	doTest(t, "./tests/hasmanywithpointers", []Comparison{
		{Have: "company_with_pointers_conditions.go", Expected: "./tests/results/hasmanywithpointers_company.go"},
		{Have: "seller_in_pointers_conditions.go", Expected: "./tests/results/hasmanywithpointers_seller.go"},
		{Have: "./tests/hasmanywithpointers/cql.go", Expected: "./tests/hasmanywithpointers/cql_result.go"},
	})
}

func TestSelfReferential(t *testing.T) {
	doTest(t, "./tests/selfreferential", []Comparison{
		{Have: "employee_conditions.go", Expected: "./tests/results/selfreferential.go"},
		{Have: "./tests/selfreferential/cql.go", Expected: "./tests/selfreferential/cql_result.go"},
	})
}

func TestMultiplePackage(t *testing.T) {
	doTest(t, "./tests/multiplepackage/package1", []Comparison{
		{Have: "package1_conditions.go", Expected: "./tests/results/multiplepackage_package1.go"},
		{Have: "./tests/multiplepackage/package1/cql.go", Expected: "./tests/multiplepackage/package1/cql_result.go"},
	})
	doTest(t, "./tests/multiplepackage/package2", []Comparison{
		{Have: "package2_conditions.go", Expected: "./tests/results/multiplepackage_package2.go"},
	})
}

func TestOverrideForeignKey(t *testing.T) {
	doTest(t, "./tests/overrideforeignkey", []Comparison{
		{Have: "bicycle_conditions.go", Expected: "./tests/results/overrideforeignkey_bicycle.go"},
		{Have: "person_conditions.go", Expected: "./tests/results/overrideforeignkey_person.go"},
		{Have: "./tests/overrideforeignkey/cql.go", Expected: "./tests/overrideforeignkey/cql_result.go"},
	})
}

func TestOverrideReferences(t *testing.T) {
	doTest(t, "./tests/overridereferences", []Comparison{
		{Have: "phone_conditions.go", Expected: "./tests/results/overridereferences_phone.go"},
		{Have: "brand_conditions.go", Expected: "./tests/results/overridereferences_brand.go"},
		{Have: "./tests/overridereferences/cql.go", Expected: "./tests/overridereferences/cql_result.go"},
	})
}

func TestOverrideForeignKeyInverse(t *testing.T) {
	doTest(t, "./tests/overrideforeignkeyinverse", []Comparison{
		{Have: "user_conditions.go", Expected: "./tests/results/overrideforeignkeyinverse_user.go"},
		{Have: "credit_card_conditions.go", Expected: "./tests/results/overrideforeignkeyinverse_credit_card.go"},
		{Have: "./tests/overrideforeignkeyinverse/cql.go", Expected: "./tests/overrideforeignkeyinverse/cql_result.go"},
	})
}

func TestOverrideReferencesInverse(t *testing.T) {
	doTest(t, "./tests/overridereferencesinverse", []Comparison{
		{Have: "computer_conditions.go", Expected: "./tests/results/overridereferencesinverse_computer.go"},
		{Have: "processor_conditions.go", Expected: "./tests/results/overridereferencesinverse_processor.go"},
		{Have: "./tests/overridereferencesinverse/cql.go", Expected: "./tests/overridereferencesinverse/cql_result.go"},
	})
}

type Comparison struct {
	Have     string
	Expected string
}

func doTest(t *testing.T, sourcePkg string, comparisons []Comparison) {
	viper.Set(DestPackageKey, "conditions")
	generateConditions(nil, []string{sourcePkg})

	for _, comparison := range comparisons {
		checkFilesEqual(t, comparison.Have, comparison.Expected)
	}
}

func checkFilesEqual(t *testing.T, file1, file2 string) {
	stat1 := testutils.CheckFileExists(t, file1)
	stat2 := testutils.CheckFileExists(t, file2)

	// do inputs at least have the same size?
	assert.Equal(t, stat1.Size(), stat2.Size(), "File lens are not equal")

	// long way: compare contents
	f1, err := os.Open(file1)
	if err != nil {
		t.Error(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		t.Error(err)
	}
	defer f2.Close()

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

	testutils.RemoveFile(file1)
}
