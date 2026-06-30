package cql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

// TestJoinedPreloadFastScan_SameModelTwiceUnderDifferentRelations exercises
// what happens when the same model (ParentParent) is joined twice under
// different relation paths (Parent1.ParentParent and Parent2.ParentParent).
//
// Each occurrence gets its own alias chain via Table.DeliverTable:
//   - Parent1                     → alias "Parent1"
//   - Parent1__ParentParent       → alias "Parent1__ParentParent"
//   - Parent2                     → alias "Parent2"
//   - Parent2__ParentParent       → alias "Parent2__ParentParent"
//
// Two RelationScanner instances for ParentParent reuse the same underlying
// per-model parentParentScanner but mount onto different parents. The fast
// path's longest-prefix matching must route each "Parent{1,2}__ParentParent__*"
// column to the right child instance.
func TestJoinedPreloadFastScan_SameModelTwiceUnderDifferentRelations(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	mock.ExpectQuery(
		`SELECT children\.\*,.*"Parent1__id".*"Parent1__ParentParent__id".*"Parent2__id".*"Parent2__ParentParent__id".*FROM "children"`,
	).WillReturnRows(
		sqlmock.NewRows([]string{
			// main
			"id", "name", "number",
			"parent1_id", "parent2_id",
			// Parent1 join
			"Parent1__id", "Parent1__name", "Parent1__parent_parent_id",
			// Parent1.ParentParent nested
			"Parent1__ParentParent__id", "Parent1__ParentParent__name", "Parent1__ParentParent__number",
			// Parent2 join
			"Parent2__id", "Parent2__parent_parent_id",
			// Parent2.ParentParent nested
			"Parent2__ParentParent__id", "Parent2__ParentParent__name", "Parent2__ParentParent__number",
		}).AddRow(
			newUUID("11111111-1111-1111-1111-111111111111"),
			"the_child", 0,
			newUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			newUUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			// Parent1 = name=p1
			newUUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), "p1",
			newUUID("cccccccc-cccc-cccc-cccc-cccccccccccc"),
			// Parent1.ParentParent = name=GP_A
			newUUID("cccccccc-cccc-cccc-cccc-cccccccccccc"), "GP_A", 1,
			// Parent2
			newUUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			newUUID("dddddddd-dddd-dddd-dddd-dddddddddddd"),
			// Parent2.ParentParent = name=GP_B
			newUUID("dddddddd-dddd-dddd-dddd-dddddddddddd"), "GP_B", 2,
		),
	)

	children, err := Query[models.Child](
		context.Background(),
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent().Preload(),
		).Preload(),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent().Preload(),
		).Preload(),
	).Find()

	require.NoError(t, err)
	require.Len(t, children, 1)

	c := children[0]
	assert.Equal(t, "the_child", c.Name)

	// Each path's ParentParent must land on its own parent — NOT cross-routed.
	assert.Equal(t, "p1", c.Parent1.Name)
	assert.Equal(t, "GP_A", c.Parent1.ParentParent.Name)
	assert.Equal(t, 1, c.Parent1.ParentParent.Number)

	assert.Equal(t, "GP_B", c.Parent2.ParentParent.Name)
	assert.Equal(t, 2, c.Parent2.ParentParent.Number)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestJoinedPreloadFastScan_SelfReferentialDeep exercises a 3-level
// self-referential preload (Employee → Boss → Boss → Boss). Validates that:
//
//   - cql-gen's codegen doesn't infinite-recurse on the cycle (it emits
//     one RelationScanner per Go field, regardless of cycles)
//   - the runtime mounts an arbitrary-depth chain correctly (depth ordering
//     in scanOneRow processes the leaf Boss__Boss__Boss before its
//     parents)
//
// The TestPreloadSelfReferentialAtSecondLevel integration test already
// covers depth-2; this is the explicit deeper variant.
func TestJoinedPreloadFastScan_SelfReferentialDeep(t *testing.T) {
	db, mock, cleanup := joinedPreloadDB(t)
	defer cleanup()

	mock.ExpectQuery(
		`SELECT employees\.\*,.*"Boss__id".*"Boss__Boss__id".*"Boss__Boss__Boss__id".*FROM "employees"`,
	).WillReturnRows(
		sqlmock.NewRows([]string{
			"id", "name", "boss_id",
			"Boss__id", "Boss__name", "Boss__boss_id",
			"Boss__Boss__id", "Boss__Boss__name", "Boss__Boss__boss_id",
			"Boss__Boss__Boss__id", "Boss__Boss__Boss__name", "Boss__Boss__Boss__boss_id",
		}).AddRow(
			newUUID("11111111-1111-1111-1111-111111111111"), "junior", newUUID("22222222-2222-2222-2222-222222222222"),
			newUUID("22222222-2222-2222-2222-222222222222"), "mid", newUUID("33333333-3333-3333-3333-333333333333"),
			newUUID("33333333-3333-3333-3333-333333333333"), "senior", newUUID("44444444-4444-4444-4444-444444444444"),
			newUUID("44444444-4444-4444-4444-444444444444"), "cto", nil,
		),
	)

	employees, err := Query[models.Employee](
		context.Background(),
		db,
		conditions.Employee.Boss(
			conditions.Employee.Boss(
				conditions.Employee.Boss().Preload(),
			).Preload(),
		).Preload(),
	).Find()

	require.NoError(t, err)
	require.Len(t, employees, 1)

	e := employees[0]
	assert.Equal(t, "junior", e.Name)

	require.NotNil(t, e.Boss)
	assert.Equal(t, "mid", e.Boss.Name)

	require.NotNil(t, e.Boss.Boss)
	assert.Equal(t, "senior", e.Boss.Boss.Name)

	require.NotNil(t, e.Boss.Boss.Boss)
	assert.Equal(t, "cto", e.Boss.Boss.Boss.Name)
	// cto has no boss — top of the chain.
	assert.Nil(t, e.Boss.Boss.Boss.Boss)

	assert.NoError(t, mock.ExpectationsWereMet())
}
