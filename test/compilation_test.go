package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	Name  string
	Code  string
	Error string
}

func TestQueryCompilationErrors(t *testing.T) {
	t.Parallel()

	queryMethods := []string{"cql.Query", "cql.Delete", "cql.Update"}

	tests := []testCase{
		{
			Name: "wrong name of condition",
			Code: `
	_ = %s[models.Product](
		context.Background(),
		db,
		conditions.ProductNotExists.Int.Is().Eq(cql.Int(1)),
	)`,
			Error: `undefined: conditions.ProductNotExists`,
		},
		{
			Name: "wrong name of property",
			Code: `
		_ = %s[models.Product](
			context.Background(),
			db,
			conditions.Product.IntNotExists.Is().Eq(cql.Int(1)),
		)`,
			Error: `conditions.Product.IntNotExists undefined (type conditions.productConditions has no field or method IntNotExists)`,
		},
		{
			Name: "basic wrong type in value",
			Code: `
		_ = %s[models.Product](
			context.Background(),
			db,
			conditions.Product.Int.Is().Eq("1"),
		)`,
			Error: `cannot use "1" (constant of type string) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Is().Eq: string does not implement condition.ValueOfType[float64] (missing method GetValue)`,
		},
		{
			Name: "Use wrong type in value",
			Code: `
		_ = %s[models.Product](
			context.Background(),
			db,
			conditions.Product.Int.Is().Eq(cql.Int("1")),
		)`,
			Error: `cannot use "1" (untyped string constant) as int value in argument to cql.Int`,
		},
		{
			Name: "Compare with wrong type",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Eq(cql.String("1")),
			)`,
			Error: `cannot use cql.String("1") (value of struct type condition.Value[string]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Is().Eq: condition.Value[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "Compare with wrong type for multiple values operator",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Between(
					cql.Int(1),
					cql.String("1"),
				),
			)`,
			Error: `cannot use cql.String("1") (value of struct type condition.Value[string]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Is().Between: condition.Value[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "Compare with wrong type for list of values operator",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().In(
					cql.Int(1),
					cql.String("1"),
				),
			)`,
			Error: `cannot use cql.String("1") (value of struct type condition.Value[string]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Is().In: condition.Value[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "Use condition of another model",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Sale.Code.Is().Eq(cql.Int(1)),
			)`,
			Error: `cannot use conditions.Sale.Code.Is().Eq(cql.Int(1)) (value of interface type condition.WhereCondition[models.Sale]) as condition.Condition[models.Product] value in argument to %s[models.Product]: condition.WhereCondition[models.Sale] does not implement condition.Condition[models.Product] (wrong type for method interfaceVerificationMethod)`,
		},
		{
			Name: "Use condition of another model inside join",
			Code: `
			_ = %s[models.Sale](
				context.Background(),
				db,
				conditions.Sale.Seller(
					conditions.Sale.Code.Is().Eq(cql.Int(1)),
				),
			)`,
			Error: `cannot use conditions.Sale.Code.Is().Eq(cql.Int(1)) (value of interface type condition.WhereCondition[models.Sale]) as condition.Condition[models.Seller] value in argument to conditions.Sale.Seller: condition.WhereCondition[models.Sale] does not implement condition.Condition[models.Seller] (wrong type for method interfaceVerificationMethod)`,
		},
		{
			Name: "Use condition of another model inside logical operator",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				cql.Or(conditions.Sale.Code.Is().Eq(cql.Int(1))),
			)`,
			Error: `cannot use cql.Or(conditions.Sale.Code.Is().Eq(cql.Int(1))) (value of interface type condition.WhereCondition[models.Sale]) as condition.Condition[models.Product] value in argument to %s[models.Product]: condition.WhereCondition[models.Sale] does not implement condition.Condition[models.Product] (wrong type for method interfaceVerificationMethod)`,
		},
		{
			Name: "Use condition of another model inside logical operator multiple",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				cql.Or[models.Product](
					conditions.Product.Int.Is().Eq(cql.Int(1)),
					conditions.Sale.Code.Is().Eq(cql.Int(1)),
				),
			)`,
			Error: `cannot use conditions.Sale.Code.Is().Eq(cql.Int(1)) (value of interface type condition.WhereCondition[models.Sale]) as condition.WhereCondition[models.Product] value in argument to cql.Or[models.Product]: condition.WhereCondition[models.Sale] does not implement condition.WhereCondition[models.Product] (wrong type for method interfaceVerificationMethod)`,
		},
		{
			Name: "Use condition of another model inside slice operator",
			Code: `
			_ = %s[models.Company](
				context.Background(),
				db,
				conditions.Company.Sellers.Any(
					conditions.Sale.Code.Is().Eq(cql.Int(1)),
				),
			)`,
			Error: `cannot use conditions.Sale.Code.Is().Eq(cql.Int(1)) (value of interface type condition.WhereCondition[models.Sale]) as condition.WhereCondition[models.Seller] value in argument to conditions.Company.Sellers.Any: condition.WhereCondition[models.Sale] does not implement condition.WhereCondition[models.Seller] (wrong type for method interfaceVerificationMethod)`,
		},
		{
			Name: "Condition with field of another type",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Eq(conditions.Product.ID),
			)`,
			Error: `cannot use conditions.Product.ID (variable of struct type condition.Field[models.Product, model.UUID]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Is().Eq: condition.Field[models.Product, model.UUID] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "Use operator not present for field type",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().True(),
			)`,
			Error: `conditions.Product.Int.Is().True undefined (type condition.NumericFieldIs[models.Product] has no field or method True)`,
		},
		{
			Name: "Use custom operator not present for field type",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Custom(
					condition.Like("_a!_").Escape('!'),
				),
			)`,
			Error: `cannot use condition.Like("_a!_").Escape('!') (value of struct type condition.ValueOperator[string]) as condition.Operator[float64] value in argument to conditions.Product.Int.Is().Custom: condition.ValueOperator[string] does not implement condition.Operator[float64] (wrong type for method InterfaceVerificationMethod)`,
		},
		{
			Name: "Use function not present for field type",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Concat(cql.String("asd")).Is().Eq(cql.Int(1)),
			)`,
			Error: `conditions.Product.Int.Concat undefined (type condition.NumericField[models.Product, int] has no field or method Concat)`,
		},
		{
			Name: "Use function with incorrect value type",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Plus(cql.String("asd")).Is().Eq(cql.Int(1)),
			)`,
			Error: `cannot use cql.String("asd") (value of struct type condition.Value[string]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Plus: condition.Value[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "Use function dynamic with incorrect value type",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Plus(conditions.Product.String).Is().Eq(cql.Int(1)),
			)`,
			Error: `cannot use conditions.Product.String (variable of struct type condition.StringField[models.Product]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Plus: condition.StringField[models.Product] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "Use function not present for field type inside comparison",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Eq(conditions.Product.Int.Concat(cql.String("asd"))),
			)`,
			Error: `conditions.Product.Int.Concat undefined (type condition.NumericField[models.Product, int] has no field or method Concat)`,
		},
		{
			Name: "Use function with incorrect value type inside comparison",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Eq(conditions.Product.Int.Plus(cql.String("asd"))),
			)`,
			Error: `cannot use cql.String("asd") (value of struct type condition.Value[string]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Plus: condition.Value[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "Use function dynamic with incorrect value type inside comparison",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Eq(conditions.Product.Int.Plus(conditions.Product.String)),
			)`,
			Error: `cannot use conditions.Product.String (variable of struct type condition.StringField[models.Product]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Plus: condition.StringField[models.Product] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "Use function with not same type of numeric value for logical operator",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Eq(conditions.Product.Int.Or(cql.Float64(1))),
			)`,
			Error: `cannot use cql.Float64(1) (value of struct type condition.NumericValue[float64]) as condition.NumericOfType[int] value in argument to conditions.Product.Int.Or: condition.NumericValue[float64] does not implement condition.NumericOfType[int] (wrong type for method GetNumericValue)`,
		},
		{
			Name: "Use function with not same type of numeric value for logical operator dynamic",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Eq(conditions.Product.Int.Or(conditions.Product.Float)),
			)`,
			Error: `cannot use conditions.Product.Float (variable of struct type condition.NumericField[models.Product, float64]) as condition.NumericOfType[int] value in argument to conditions.Product.Int.Or: condition.NumericField[models.Product, float64] does not implement condition.NumericOfType[int] (wrong type for method GetNumericValue)`,
		},
		{
			Name: "Use function with not int type of numeric value for shift operator",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Eq(conditions.Product.Int.ShiftLeft(cql.Float64(1))),
			)`,
			Error: `cannot use cql.Float64(1) (value of struct type condition.NumericValue[float64]) as condition.NumericOfType[int] value in argument to conditions.Product.Int.ShiftLeft: condition.NumericValue[float64] does not implement condition.NumericOfType[int] (wrong type for method GetNumericValue)`,
		},
		{
			Name: "Use function with not int type of numeric value for shift operator dynamic",
			Code: `
			_ = %s[models.Product](
				context.Background(),
				db,
				conditions.Product.Int.Is().Eq(conditions.Product.Int.ShiftLeft(conditions.Product.Float)),
			)`,
			Error: `cannot use conditions.Product.Float (variable of struct type condition.NumericField[models.Product, float64]) as condition.NumericOfType[int] value in argument to conditions.Product.Int.ShiftLeft: condition.NumericField[models.Product, float64] does not implement condition.NumericOfType[int] (wrong type for method GetNumericValue)`,
		},
	}

	for _, testCase := range tests {
		for _, method := range queryMethods {
			internalTestCase := testCase

			t.Run(method+"_"+internalTestCase.Name, func(t *testing.T) {
				t.Parallel()

				internalTestCase.Code = fmt.Sprintf(internalTestCase.Code, method)

				if strings.Contains(internalTestCase.Error, "%s") {
					internalTestCase.Error = fmt.Sprintf(internalTestCase.Error, method)
				}

				executeTest(t, "", internalTestCase)
			})
		}
	}
}

func executeTest(t *testing.T, extraCode string, testCase testCase) {
	t.Helper()

	code := `
package main

import (
	"context"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/models"
	"github.com/FrancoLiberali/cql/test/conditions"
)

var db *cql.DB
` + extraCode + `

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
}

func TestGroupByCompilationErrors(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		{
			Name: "aggregation do not exist for value type",
			Code: `
		_ = cql.Query[models.Product](
			context.Background(),
			db,
		).GroupBy(
			conditions.Product.Int,
		).SelectValue(
			conditions.Product.Int.Aggregate().All(), "aggregation1",
		)`,
			Error: `conditions.Product.Int.Aggregate().All undefined (type condition.NumericFieldAggregation has no field or method All)`,
		},
		{
			Name: "having not compared with correct type of value",
			Code: `
		_ = cql.Query[models.Product](
			context.Background(),
			db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Int.Aggregate().Max().Eq(cql.String("13")),
		).SelectValue(
			conditions.Product.Int.Aggregate().Max(), "aggregation1",
		)`,
			Error: `cannot use cql.String("13") (value of struct type condition.Value[string]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Aggregate().Max().Eq: condition.Value[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "having not compared with correct type of another aggregation",
			Code: `
		_ = cql.Query[models.Product](
			context.Background(),
			db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Int.Aggregate().Max().Eq(conditions.Product.String.Aggregate().Min()),
		).SelectValue(
			conditions.Product.Int.Aggregate().Max(), "aggregation1",
		)`,
			Error: ` cannot use conditions.Product.String.Aggregate().Min() (value of struct type condition.AggregationResult[string]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Aggregate().Max().Eq: condition.AggregationResult[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			executeTest(t, "", testCase)
		})
	}
}

func TestUpdateCompilationErrors(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		{
			Name: "set value of wrong type",
			Code: `
		_, _ = cql.Update[models.Product](
			context.Background(),
			db,
			conditions.Product.Bool.Is().False(),
		).Set(
			conditions.Product.Int.Set().Eq(cql.String("1")),
		)`,
			Error: `cannot use cql.String("1") (value of struct type condition.Value[string]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Set().Eq: condition.Value[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "set field of wrong type",
			Code: `
		_, _ = cql.Update[models.Product](
			context.Background(),
			db,
			conditions.Product.Bool.Is().False(),
		).Set(
			conditions.Product.Int.Set().Eq(conditions.Product.String),
		)`,
			Error: `cannot use conditions.Product.String (variable of struct type condition.StringField[models.Product]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Set().Eq: condition.StringField[models.Product] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "set multiple value of wrong type",
			Code: `
		_, _ = cql.Update[models.Product](
			context.Background(),
			db,
			conditions.Product.Bool.Is().False(),
		).SetMultiple(
			conditions.Product.Int.Set().Eq(cql.String("1")),
		)`,
			Error: `cannot use cql.String("1") (value of struct type condition.Value[string]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Set().Eq: condition.Value[string] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "set field of wrong type",
			Code: `
		_, _ = cql.Update[models.Product](
			context.Background(),
			db,
			conditions.Product.Bool.Is().False(),
		).SetMultiple(
			conditions.Product.Int.Set().Eq(conditions.Product.String),
		)`,
			Error: `cannot use conditions.Product.String (variable of struct type condition.StringField[models.Product]) as condition.ValueOfType[float64] value in argument to conditions.Product.Int.Set().Eq: condition.StringField[models.Product] does not implement condition.ValueOfType[float64] (wrong type for method GetValue)`,
		},
		{
			Name: "set can not be used after a function",
			Code: `
		_, _ = cql.Update[models.Product](
			context.Background(),
			db,
			conditions.Product.Bool.Is().False(),
		).Set(
			conditions.Product.Int.Plus(1).Set().Eq(cql.Int(1)),
		)`,
			Error: `conditions.Product.Int.Plus(1).Set undefined (type condition.NotUpdatableNumericField[models.Product, int] has no field or method Set)`,
		},
		{
			Name: "set null can not be used for not nullable types",
			Code: `
		_, _ = cql.Update[models.Product](
			context.Background(),
			db,
			conditions.Product.Bool.Is().False(),
		).Set(
			conditions.Product.Int.Set().Null(),
		)`,
			Error: `conditions.Product.Int.Set().Null undefined (type condition.NumericFieldSet[models.Product, int] has no field or method Null)`,
		},
		{
			Name: "returning model must be the same as query",
			Code: `
		productsReturned := []models.Seller{}

		_, _ = cql.Update[models.Product](
			context.Background(),
			db,
			conditions.Product.Bool.Is().False(),
		).Returning(
			&productsReturned,
		).Set(
			conditions.Product.Int.Set().Eq(cql.Int(1)),
		)`,
			Error: `cannot use &productsReturned (value of type *[]models.Seller) as *[]models.Product value in argument to cql.Update[models.Product](context.Background(), db, conditions.Product.Bool.Is().False()).Returning`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			executeTest(t, "", testCase)
		})
	}
}

func TestSelectCompilationErrors(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		{
			Name: "value into different destinations",
			Code: `
				_, _ = cql.Select(
					cql.Query[models.Product](
						context.Background(),
						db,
					),
					cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
						result.Int = int(value)
					}),
					cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt2) {
						result.Int = int(value)
					}),
				)
			`,
			Error: `in call to cql.Select, type *cql.ValueIntoSelection[float64, ResultInt2] of cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt2) {…}) does not match inferred type condition.Selection[ResultInt] for condition.Selection[TResults]`,
		},
		{
			Name: "value not the same time of the query",
			Code: `
				_, _ = cql.Select(
					cql.Query[models.Product](
						context.Background(),
						db,
					),
					cql.ValueInto(conditions.Product.Int, func(value string, result *ResultInt) {
					}),
				)
			`,
			Error: `in call to cql.ValueInto, type func(value string, result *ResultInt) of func(value string, result *ResultInt) {} does not match inferred type func(float64, *TResults) for func(TValue, *TResults)`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			executeTest(t, `
type ResultInt struct {
	Int          int
}

type ResultInt2 struct {
	Int          int
}
`, testCase)
		})
	}
}

func TestInsertCompilationErrors(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		{
			Name: "no other conflict can be called after DoNothing",
			Code: `
				_, _ = cql.Insert(
					context.Background(),
					db,
					product,
				).OnConflict().DoNothing().OnConflict().DoNothing().Exec()
			`,
			Error: `cql.Insert(context.Background(), db, product).OnConflict().DoNothing().OnConflict undefined (type *condition.InsertExec[models.Product] has no field or method OnConflict)`,
		},
		{
			Name: "no other conflict can be called after UpdateAll",
			Code: `
				_, _ = cql.Insert(
					context.Background(),
					db,
					product,
				).OnConflict().UpdateAll().OnConflict().UpdateAll().Exec()
			`,
			Error: `cql.Insert(context.Background(), db, product).OnConflict().UpdateAll().OnConflict undefined (type *condition.InsertExec[models.Product] has no field or method OnConflict)`,
		},
		{
			Name: "no other conflict can be called after Update",
			Code: `
				_, _ = cql.Insert(
					context.Background(),
					db,
					product,
				).OnConflict().Update().OnConflict().Update().Exec()
			`,
			Error: `cql.Insert(context.Background(), db, product).OnConflict().Update().OnConflict undefined (type *condition.InsertExec[models.Product] has no field or method OnConflict)`,
		},
		{
			Name: "no other conflict can be called after Set",
			Code: `
				_, _ = cql.Insert(
					context.Background(),
					db,
					product,
				).OnConflict().Set().OnConflict().Set().Exec()
			`,
			Error: `cql.Insert(context.Background(), db, product).OnConflict().Set().OnConflict undefined (type *condition.InsertOnConflictSet[models.Product] has no field or method OnConflict)`,
		},
		{
			Name: "no other conflict can be called after Where",
			Code: `
				_, _ = cql.Insert(
					context.Background(),
					db,
					product,
				).OnConflict().Set().Where().OnConflict().Set().Exec()
			`,
			Error: `cql.Insert(context.Background(), db, product).OnConflict().Set().Where().OnConflict undefined (type *condition.InsertExec[models.Product] has no field or method OnConflict)`,
		},
		{
			Name: "on conflict on field of different model",
			Code: `
				_, _ = cql.Insert(
					context.Background(),
					db,
					product,
				).OnConflictOn(conditions.City.ID).Update(conditions.Product.Int).Exec()
			`,
			Error: `cannot use conditions.City.ID (variable of struct type condition.Field[models.City, model.UUID]) as condition.FieldOfModel[models.Product] value in argument to cql.Insert(context.Background(), db, product).OnConflictOn: condition.Field[models.City, model.UUID] does not implement condition.FieldOfModel[models.Product] (wrong type for method getModel)`,
		},
		{
			Name: "update field of different model",
			Code: `
				_, _ = cql.Insert(
					context.Background(),
					db,
					product,
				).OnConflictOn(conditions.Product.ID).Update(conditions.City.ID).Exec()
			`,
			Error: `cannot use conditions.City.ID (variable of struct type condition.Field[models.City, model.UUID]) as condition.FieldOfModel[models.Product] value in argument to cql.Insert(context.Background(), db, product).OnConflictOn(conditions.Product.ID).Update: condition.Field[models.City, model.UUID] does not implement condition.FieldOfModel[models.Product] (wrong type for method getModel)`,
		},
		{
			Name: "set field of different model",
			Code: `
				_, _ = cql.Insert(
					context.Background(),
					db,
					product,
				).OnConflictOn(conditions.Product.ID).Set(
					conditions.City.Name.Set().Eq(cql.String("asd")),
				).Exec()
			`,
			Error: `cannot use conditions.City.Name.Set().Eq(cql.String("asd")) (value of type *condition.Set[models.City]) as *condition.Set[models.Product] value in argument to cql.Insert(context.Background(), db, product).OnConflictOn(conditions.Product.ID).Set`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			executeTest(t, `
var product = &models.Product{}
`, testCase)
		})
	}
}
