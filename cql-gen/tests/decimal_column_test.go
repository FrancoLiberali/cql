package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"gotest.tools/assert"

	"github.com/FrancoLiberali/cql/cql-gen/cmd"
)

// TestDecimalColumn proves cql-gen emits DecimalField for a decimal-tagged,
// arithmetic-capable column, and leaves a non-decimal custom column as a plain
// Field.
func TestDecimalColumn(t *testing.T) {
	viper.Set(cmd.DestPackageKey, "conditions")

	cmd.GenerateConditions(nil, []string{"./decimalcolumn"})

	defer os.Remove("wallet_conditions.go")
	defer os.Remove("wallet_scanner.go")
	defer os.Remove("./decimalcolumn/cql.go")

	generated, err := os.ReadFile("wallet_conditions.go")
	assert.NilError(t, err)

	content := string(generated)

	// Balance: decimal-tagged + arithmetic type -> DecimalField.
	assert.Assert(t,
		strings.Contains(content, "condition.DecimalField[decimalcolumn.Wallet, decimalcolumn.Money]"),
		"expected Balance to be a DecimalField, generated:\n%s", content,
	)
	assert.Assert(t,
		strings.Contains(content, "condition.NewDecimalField[decimalcolumn.Wallet, decimalcolumn.Money]"),
		"expected NewDecimalField constructor, generated:\n%s", content,
	)

	// Note: custom Valuer, no arithmetic, no decimal tag -> stays a plain Field.
	assert.Assert(t,
		strings.Contains(content, "condition.UpdatableField[decimalcolumn.Wallet, decimalcolumn.Note]"),
		"expected Note to stay a plain Field, generated:\n%s", content,
	)
	assert.Assert(t,
		!strings.Contains(content, "decimalcolumn.Note]") || !strings.Contains(content, "DecimalField[decimalcolumn.Wallet, decimalcolumn.Note]"),
		"Note must not be a DecimalField",
	)
}
