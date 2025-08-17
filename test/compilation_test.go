package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompilationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name  string
		Code  string
		Error string
	}{
		{
			Name: "Use wrong type in value",
			Code: `
	_ = cql.Query[models.Product](
		db,
		conditions.Product.Int.Is().Eq(cql.Int("1")),
	)`,
			Error: `cannot use "1" (untyped string constant) as int value in argument to cql.Int`,
		},
		{
			Name: "Use condition of another model",
			Code: `
		_ = cql.Query[models.Product](
			db,
			conditions.Sale.Code.Is().Eq(cql.Int(1)),
		)`,
			Error: `cannot use conditions.Sale.Code.Is().Eq(cql.Int(1)) (value of type condition.WhereCondition[models.Sale]) as condition.Condition[models.Product] value in argument to cql.Query[models.Product]: condition.WhereCondition[models.Sale] does not implement condition.Condition[models.Product] (wrong type for method interfaceVerificationMethod)`,
		},
		{
			Name: "Compare with wrong type",
			Code: `
		_ = cql.Query[models.Product](
			db,
			conditions.Sale.Code.Is().Eq(cql.String("1")),
		)`,
			Error: `cannot use cql.String("1") (value of type condition.Value[string]) as condition.ValueOfType[float64] value in argument to conditions.Sale.Code.Is().Eq: condition.Value[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			code := `
package main

import (
	"gorm.io/gorm"
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/models"
	"github.com/FrancoLiberali/cql/test/conditions"
)

var db *gorm.DB

func main() {
` + testCase.Code + `
}
`

			tempDir := t.TempDir()

			f, err := os.CreateTemp(tempDir, "cql-test-*.go")
			require.NoError(t, err)

			// Write data to the temporary file
			_, err = f.WriteString(code)
			require.NoError(t, err)

			cmd := exec.Command("go", "build", "-o", f.Name()+".exe", f.Name()) //nolint:gosec // necessary for the test

			output, err := cmd.CombinedOutput()
			require.Error(t, err)

			assert.Contains(t, string(output), testCase.Error)
		})
	}
}
