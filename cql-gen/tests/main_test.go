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
		{Have: "uint_model_conditions.go", Expected: "./results/uintmodel.go"},
	})
	CheckFileNotExists(t, "./uintmodel/cql.go")
}

func TestUUIDModel(t *testing.T) {
	doTest(t, "./uuidmodel", []Comparison{
		{Have: "uuid_model_conditions.go", Expected: "./results/uuidmodel.go"},
	})
	CheckFileNotExists(t, "./uuidmodel/cql.go")
}

func TestBasicTypes(t *testing.T) {
	doTest(t, "./basictypes", []Comparison{
		{Have: "basic_types_conditions.go", Expected: "./results/basictypes.go"},
	})
	CheckFileNotExists(t, "./basictypes/cql.go")
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

func TestGoEmbedded(t *testing.T) {
	doTest(t, "./goembedded", []Comparison{
		{Have: "go_embedded_conditions.go", Expected: "./results/goembedded.go"},
	})
	CheckFileNotExists(t, "./goembedded/cql.go")
}

func TestGormEmbedded(t *testing.T) {
	doTest(t, "./gormembedded", []Comparison{
		{Have: "gorm_embedded_conditions.go", Expected: "./results/gormembedded.go"},
	})
	CheckFileNotExists(t, "./gormembedded/cql.go")
}

func TestCustomType(t *testing.T) {
	doTest(t, "./customtype", []Comparison{
		{Have: "custom_type_conditions.go", Expected: "./results/customtype.go"},
	})
	CheckFileNotExists(t, "./customtype/cql.go")
}

func TestColumnDefinition(t *testing.T) {
	doTest(t, "./columndefinition", []Comparison{
		{Have: "column_definition_conditions.go", Expected: "./results/columndefinition.go"},
	})
	CheckFileNotExists(t, "./columndefinition/cql.go")
}

func TestNullableTypes(t *testing.T) {
	doTest(t, "./nullabletypes", []Comparison{
		{Have: "nullable_types_conditions.go", Expected: "./results/nullabletypes.go"},
	})
	CheckFileNotExists(t, "./nullabletypes/cql.go")
}

func TestBelongsTo(t *testing.T) {
	doTest(t, "./belongsto", []Comparison{
		{Have: "owner_conditions.go", Expected: "./results/belongsto_owner.go"},
		{Have: "owned_conditions.go", Expected: "./results/belongsto_owned.go"},
		{Have: "./belongsto/cql.go", Expected: "./belongsto/cql_result.go"},
	})
}

func TestHasOne(t *testing.T) {
	doTest(t, "./hasone", []Comparison{
		{Have: "country_conditions.go", Expected: "./results/hasone_country.go"},
		{Have: "city_conditions.go", Expected: "./results/hasone_city.go"},
		{Have: "./hasone/cql.go", Expected: "./hasone/cql_result.go"},
	})
}

func TestHasMany(t *testing.T) {
	doTest(t, "./hasmany", []Comparison{
		{Have: "company_conditions.go", Expected: "./results/hasmany_company.go"},
		{Have: "seller_conditions.go", Expected: "./results/hasmany_seller.go"},
		{Have: "./hasmany/cql.go", Expected: "./hasmany/cql_result.go"},
	})
}

func TestHasManyWithPointers(t *testing.T) {
	doTest(t, "./hasmanywithpointers", []Comparison{
		{Have: "company_with_pointers_conditions.go", Expected: "./results/hasmanywithpointers_company.go"},
		{Have: "seller_in_pointers_conditions.go", Expected: "./results/hasmanywithpointers_seller.go"},
		{Have: "./hasmanywithpointers/cql.go", Expected: "./hasmanywithpointers/cql_result.go"},
	})
}

func TestSelfReferential(t *testing.T) {
	doTest(t, "./selfreferential", []Comparison{
		{Have: "employee_conditions.go", Expected: "./results/selfreferential.go"},
		{Have: "./selfreferential/cql.go", Expected: "./selfreferential/cql_result.go"},
	})
}

func TestMultiplePackage(t *testing.T) {
	doTest(t, "./multiplepackage/package1", []Comparison{
		{Have: "package1_conditions.go", Expected: "./results/multiplepackage_package1.go"},
		{Have: "./multiplepackage/package1/cql.go", Expected: "./multiplepackage/package1/cql_result.go"},
	})
	doTest(t, "./multiplepackage/package2", []Comparison{
		{Have: "package2_conditions.go", Expected: "./results/multiplepackage_package2.go"},
	})
}

func TestOverrideForeignKey(t *testing.T) {
	doTest(t, "./overrideforeignkey", []Comparison{
		{Have: "bicycle_conditions.go", Expected: "./results/overrideforeignkey_bicycle.go"},
		{Have: "person_conditions.go", Expected: "./results/overrideforeignkey_person.go"},
		{Have: "./overrideforeignkey/cql.go", Expected: "./overrideforeignkey/cql_result.go"},
	})
}

func TestOverrideReferences(t *testing.T) {
	doTest(t, "./overridereferences", []Comparison{
		{Have: "phone_conditions.go", Expected: "./results/overridereferences_phone.go"},
		{Have: "brand_conditions.go", Expected: "./results/overridereferences_brand.go"},
		{Have: "./overridereferences/cql.go", Expected: "./overridereferences/cql_result.go"},
	})
}

func TestOverrideForeignKeyInverse(t *testing.T) {
	doTest(t, "./overrideforeignkeyinverse", []Comparison{
		{Have: "user_conditions.go", Expected: "./results/overrideforeignkeyinverse_user.go"},
		{Have: "credit_card_conditions.go", Expected: "./results/overrideforeignkeyinverse_credit_card.go"},
		{Have: "./overrideforeignkeyinverse/cql.go", Expected: "./overrideforeignkeyinverse/cql_result.go"},
	})
}

func TestOverrideReferencesInverse(t *testing.T) {
	doTest(t, "./overridereferencesinverse", []Comparison{
		{Have: "computer_conditions.go", Expected: "./results/overridereferencesinverse_computer.go"},
		{Have: "processor_conditions.go", Expected: "./results/overridereferencesinverse_processor.go"},
		{Have: "./overridereferencesinverse/cql.go", Expected: "./overridereferencesinverse/cql_result.go"},
	})
}

type Comparison struct {
	Have     string
	Expected string
}

func doTest(t *testing.T, sourcePkg string, comparisons []Comparison) {
	viper.Set(cmd.DestPackageKey, "conditions")
	cmd.GenerateConditions(nil, []string{sourcePkg})

	for _, comparison := range comparisons {
		checkFilesEqual(t, comparison.Have, comparison.Expected)
	}
}

func checkFilesEqual(t *testing.T, file1, file2 string) {
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

	RemoveFile(file1)
}
