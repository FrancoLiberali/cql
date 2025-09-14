package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testSameModel() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testSameModelWithoutIndex() {
	cql.Query(
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedModel() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithJoinedWithCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithJoinedWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand().Preload(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithJoinedWithConditionsWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		).Preload(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedModelInVariable() {
	value := conditions.Brand.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(value),
	).Find()
}

func testNotJoinedInSameLine() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testNotJoinedInDifferentLines() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(
			conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedInVariable() {
	value := conditions.City.Name

	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(value), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testNotJoinedWithTrue() {
	cql.Query[models.Brand](
		db,
		cql.True[models.Brand](),
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedInsideConnector() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
		cql.And(
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
	).Find()
}

func testNotJoinedInsideConnector() {
	cql.Query[models.Brand](
		db,
		cql.And(
			conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testJoinedInsideJoinCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
	).Find()
}

func testNotJoinedInsideJoinCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testJoinedInSecondCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
	).Find()
}

func testNotJoinedInSecondCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testNotJoinedInsideNestedJoinCondition() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			),
		),
	).Find()
}

func testJoinedInsideNestedJoinConditionWithMainModel() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.Is().Eq(conditions.Child.Name),
			),
		),
	).Find()
}

func testJoinedInsideNestedJoinConditionWithPreviousJoin() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.Is().Eq(conditions.Parent1.Name),
			),
		),
	).Find()
}

func testNotJoinedWithJoinedWithConditionBefore() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name), // want "github.com/FrancoLiberali/cql/test/models.Brand is not joined by the query"
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		),
	).Find()
}

func testJoinedWithDifferentRelationNameWithConditionsUsesConditionName() {
	cql.Query[models.Bicycle](
		db,
		conditions.Bicycle.Owner(
			conditions.Person.Name.Is().Eq(cql.String("asd")),
		),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithDifferentRelationNameWithConditionsWithPreloadUsesConditionName() {
	cql.Query[models.Bicycle](
		db,
		conditions.Bicycle.Owner(
			conditions.Person.Name.Is().Eq(cql.String("asd")),
		).Preload(),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithDifferentRelationNameWithoutConditions() {
	cql.Query[models.Bicycle](
		db,
		conditions.Bicycle.Owner(),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithDifferentRelationNameWithoutConditionsWithPreload() {
	cql.Query[models.Bicycle](
		db,
		conditions.Bicycle.Owner().Preload(),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithAppearance() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name.Appearance(0)),
	).Find()
}

func testJoinedWithAppearanceVariable() {
	value := conditions.Brand.Name.Appearance(0)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(value),
	).Find()
}

func testNotJoinedWithAppearance() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name.Appearance(0)), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedWithAppearanceVariable() {
	value := conditions.City.Name.Appearance(0)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testJoinedWithFunction() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Concat(cql.String("asd"))),
		),
	).Find()
}

func testJoinedWithFunctionVariable() {
	value := conditions.Phone.Name.Concat(cql.String("asd"))

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value),
		),
	).Find()
}

func testJoinedWithFunctionOverVariable() {
	value := conditions.Phone.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value.Concat(cql.String("asd"))),
		),
	).Find()
}

func testNotJoinedWithFunction() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name.Concat(cql.String("asd"))), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedWithFunctionVariable() {
	value := conditions.City.Name.Concat(cql.String("asd"))

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedWithFunctionOverVariable() {
	value := conditions.City.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value.Concat(cql.String("asd"))), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedWithTwoFunctions() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name.Concat(cql.String("asd")).Concat(cql.String("asd"))), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testMultipleArgumentsFirstNotJoined() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			conditions.Brand.Name,
		),
	).Find()
}

func testMultipleArgumentsFirstNotJoinedWithVariable() {
	value := conditions.City.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			value, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			conditions.Brand.Name,
		),
	).Find()
}

func testMultipleArgumentsSecondNotJoined() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			conditions.Brand.Name,
			conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testMultipleArgumentsSecondNotJoinedWithVariable() {
	value := conditions.City.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			conditions.Brand.Name,
			value, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testJoinedWithFunctionOnLeftSide() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Concat(conditions.Phone.Name).Is().Eq(conditions.Phone.Name.Concat(cql.String("asd"))),
		),
	).Find()
}

func testJoinedWithFunctionVariableOnLeftSide() {
	value := conditions.Brand.Name.Concat(conditions.Phone.Name)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			value.Is().Eq(conditions.Phone.Name.Concat(cql.String("asd"))),
		),
	).Find()
}

func testJoinedWithFunctionOverVariableOnLeftSide() {
	value := conditions.Brand.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			value.Concat(conditions.Phone.Name).Is().Eq(conditions.Phone.Name.Concat(cql.String("asd"))),
		),
	).Find()
}

func testJoinedWithFunctionWithVariableOnLeftSide() {
	value := conditions.Phone.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Concat(value).Is().Eq(conditions.Phone.Name.Concat(cql.String("asd"))),
		),
	).Find()
}

func testNotJoinedWithFunctionOnLeftSide() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Concat(conditions.City.Name).Is().Eq(conditions.Phone.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedWithFunctionOverVariableOnLeftSide() {
	value := conditions.City.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Concat(value).Is().Eq(conditions.Phone.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testJoinedWithFunctionOnRightSide() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Concat(conditions.Phone.Name)),
		),
	).Find()
}

func testJoinedWithFunctionVariableOnRightSide() {
	value := conditions.Brand.Name.Concat(conditions.Phone.Name)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value),
		),
	).Find()
}

func testJoinedWithFunctionOverVariableOnRightSide() {
	value := conditions.Brand.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Concat(value)),
		),
	).Find()
}

func testNotJoinedWithFunctionOnRightSide() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Concat(conditions.City.Name)), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedWithMultipleFunctionOnRightSideFirst() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Concat(
				conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			).Concat(
				conditions.Phone.Name,
			)),
		),
	).Find()
}

func testNotJoinedWithMultipleFunctionOnRightSideSecond() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Concat(
				conditions.Phone.Name,
			).Concat(
				conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			)),
		),
	).Find()
}

func testNotJoinedWithMultipleFunctionOnRightSideTwice() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Concat(
				conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			).Concat(
				conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			)),
		),
	).Find()
}

func testNotJoinedWithFunctionOverVariableOnRightSide() {
	value := conditions.City.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Concat(value)), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testJoinedConditionInVariable() {
	value := conditions.Phone.Name.Is().Eq(conditions.Phone.Name)

	cql.Query[models.Phone](
		db,
		value,
	).Find()
}

func testNotJoinedConditionInVariable() {
	value := conditions.Phone.Name.Is().Eq(conditions.City.Name) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		db,
		value,
	).Find()
}

func testJoinedConditionInList() {
	values := []condition.Condition[models.Phone]{
		conditions.Phone.Name.Is().Eq(conditions.Phone.Name),
	}

	cql.Query[models.Phone](
		db,
		values...,
	).Find()
}

func testNotJoinedConditionInList() {
	values := []condition.Condition[models.Phone]{
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	}

	cql.Query[models.Phone](
		db,
		values...,
	).Find()
}

func testJoinedConditionInListWithAppend() {
	values := []condition.Condition[models.Phone]{}

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.Phone.Name),
	)

	cql.Query[models.Phone](
		db,
		values...,
	).Find()
}

func testNotJoinedConditionInListWithAppend() {
	values := []condition.Condition[models.Phone]{}

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		values...,
	).Find()
}

func testNotJoinedConditionInListWithAppendSecond() {
	values := []condition.Condition[models.Phone]{}

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.Phone.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		values...,
	).Find()
}

func testNotJoinedConditionInListWithAppendMultiple() {
	values := []condition.Condition[models.Phone]{}

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.Phone.Name),
	)

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		values...,
	).Find()
}
