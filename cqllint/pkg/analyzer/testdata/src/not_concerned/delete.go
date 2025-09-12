package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testDeleteSameModel() {
	cql.Delete[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedModel() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedWithJoinedWithCondition() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedWithJoinedWithPreload() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand().Preload(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedWithJoinedWithConditionsWithPreload() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		).Preload(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedModelInVariable() {
	value := conditions.Brand.Name

	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(value),
	).Exec()
}

func testDeleteNotJoinedInSameLine() {
	cql.Delete[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteNotJoinedInDifferentLines() {
	cql.Delete[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(
			conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteNotJoinedInVariable() {
	value := conditions.City.Name

	cql.Delete[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(value), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteNotJoinedWithTrue() {
	cql.Delete[models.Brand](
		db,
		cql.True[models.Brand](),
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedInsideConnector() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(),
		cql.And(
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
	).Exec()
}

func testDeleteNotJoinedInsideConnector() {
	cql.Delete[models.Brand](
		db,
		cql.And(
			conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteJoinedInsideJoinCondition() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
	).Exec()
}

func testDeleteNotJoinedInsideJoinCondition() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteJoinedInSecondCondition() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
	).Exec()
}

func testDeleteNotJoinedInSecondCondition() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteNotJoinedInsideNestedJoinCondition() {
	cql.Delete[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			),
		),
	).Exec()
}

func testDeleteJoinedInsideNestedJoinConditionWithMainModel() {
	cql.Delete[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.Is().Eq(conditions.Child.Name),
			),
		),
	).Exec()
}

func testDeleteJoinedInsideNestedJoinConditionWithPreviousJoin() {
	cql.Delete[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.Is().Eq(conditions.Parent1.Name),
			),
		),
	).Exec()
}

func testDeleteNotJoinedWithJoinedWithConditionBefore() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name), // want "github.com/FrancoLiberali/cql/test/models.Brand is not joined by the query"
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		),
	).Exec()
}

func testDeleteJoinedWithDifferentRelationNameWithConditionsUsesConditionName() {
	cql.Delete[models.Bicycle](
		db,
		conditions.Bicycle.Owner(
			conditions.Person.Name.Is().Eq(cql.String("asd")),
		),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedWithDifferentRelationNameWithConditionsWithPreloadUsesConditionName() {
	cql.Delete[models.Bicycle](
		db,
		conditions.Bicycle.Owner(
			conditions.Person.Name.Is().Eq(cql.String("asd")),
		).Preload(),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedWithDifferentRelationNameWithoutConditions() {
	cql.Delete[models.Bicycle](
		db,
		conditions.Bicycle.Owner(),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedWithDifferentRelationNameWithoutConditionsWithPreload() {
	cql.Delete[models.Bicycle](
		db,
		conditions.Bicycle.Owner().Preload(),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testDeleteJoinedWithAppearance() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name.Appearance(0)),
	).Exec()
}

func testDeleteJoinedWithAppearanceVariable() {
	value := conditions.Brand.Name.Appearance(0)

	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(value),
	).Exec()
}

func testDeleteNotJoinedWithAppearance() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name.Appearance(0)), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteNotJoinedWithAppearanceVariable() {
	value := conditions.City.Name.Appearance(0)

	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteJoinedWithFunction() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Concat(cql.String("asd"))),
		),
	).Exec()
}

func testDeleteJoinedWithFunctionVariable() {
	value := conditions.Phone.Name.Concat(cql.String("asd"))

	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value),
		),
	).Exec()
}

func testDeleteJoinedWithFunctionOverVariable() {
	value := conditions.Phone.Name

	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value.Concat(cql.String("asd"))),
		),
	).Exec()
}

func testDeleteNotJoinedWithFunction() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name.Concat(cql.String("asd"))), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteNotJoinedWithFunctionVariable() {
	value := conditions.City.Name.Concat(cql.String("asd"))

	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteNotJoinedWithFunctionOverVariable() {
	value := conditions.City.Name

	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(value.Concat(cql.String("asd"))), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteNotJoinedWithTwoFunctions() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name.Concat(cql.String("asd")).Concat(cql.String("asd"))), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteMultipleArgumentsFirstNotJoined() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			conditions.Brand.Name,
		),
	).Exec()
}

func testDeleteMultipleArgumentsFirstNotJoinedWithVariable() {
	value := conditions.City.Name

	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			value, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			conditions.Brand.Name,
		),
	).Exec()
}

func testDeleteMultipleArgumentsSecondNotJoined() {
	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			conditions.Brand.Name,
			conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteMultipleArgumentsSecondNotJoinedWithVariable() {
	value := conditions.City.Name

	cql.Delete[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			conditions.Brand.Name,
			value, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Exec()
}

func testDeleteJoinedConditionInVariable() {
	value := conditions.Phone.Name.Is().Eq(conditions.Phone.Name)

	cql.Delete[models.Phone](
		db,
		value,
	).Exec()
}

func testDeleteNotJoinedConditionInVariable() {
	value := conditions.Phone.Name.Is().Eq(conditions.City.Name) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Delete[models.Phone](
		db,
		value,
	).Exec()
}

func testDeleteJoinedConditionInList() {
	values := []condition.Condition[models.Phone]{
		conditions.Phone.Name.Is().Eq(conditions.Phone.Name),
	}

	cql.Delete[models.Phone](
		db,
		values...,
	).Exec()
}

func testDeleteNotJoinedConditionInList() {
	values := []condition.Condition[models.Phone]{
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	}

	cql.Delete[models.Phone](
		db,
		values...,
	).Exec()
}

func testDeleteJoinedConditionInListWithAppend() {
	values := []condition.Condition[models.Phone]{}

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.Phone.Name),
	)

	cql.Delete[models.Phone](
		db,
		values...,
	).Exec()
}

func testDeleteNotJoinedConditionInListWithAppend() {
	values := []condition.Condition[models.Phone]{}

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Delete[models.Phone](
		db,
		values...,
	).Exec()
}

func testDeleteNotJoinedConditionInListWithAppendSecond() {
	values := []condition.Condition[models.Phone]{}

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.Phone.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Delete[models.Phone](
		db,
		values...,
	).Exec()
}

func testDeleteNotJoinedConditionInListWithAppendMultiple() {
	values := []condition.Condition[models.Phone]{}

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.Phone.Name),
	)

	values = append(
		values,
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Delete[models.Phone](
		db,
		values...,
	).Exec()
}
