package gormprobes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type probeBrand struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func (probeBrand) TableName() string { return "brands" }

type probePhone struct {
	ID      uint `gorm:"primarykey"`
	Name    string
	BrandID uint
	Brand   probeBrand
}

func (probePhone) TableName() string { return "phones" }

// TestLeftJoinNoMatchIsAllNull locks in the signal the fast scanner keys on:
// when a LEFT JOIN finds no matching child row, every joined column comes back
// NULL. scanOne treats an all-NULL joined block as "no child" and skips
// mounting it onto the parent.
func TestLeftJoinNoMatchIsAllNull(t *testing.T) {
	t.Parallel()

	db := openProbeDB(t)
	require.NoError(t, db.AutoMigrate(&probeBrand{}, &probePhone{}))

	// brand_id 999 matches no brand row.
	require.NoError(t, db.Create(&probePhone{Name: "p_no_brand", BrandID: 999}).Error)

	rows, err := db.Raw(`
		SELECT phones.*,
		       Brand.id   AS "Brand__id",
		       Brand.name AS "Brand__name"
		  FROM phones
		  LEFT JOIN brands AS Brand ON Brand.id = phones.brand_id
	`).Rows()
	require.NoError(t, err)

	defer rows.Close()

	require.True(t, rows.Next(), "expected one phone row")

	var (
		pid      uint
		pname    string
		pbrand   uint
		bidVal   any
		bnameVal any
	)

	require.NoError(t, rows.Scan(&pid, &pname, &pbrand, &bidVal, &bnameVal))

	require.Nil(t, bidVal, "unmatched LEFT JOIN must yield NULL joined columns")
	require.Nil(t, bnameVal, "unmatched LEFT JOIN must yield NULL joined columns")
}

type probeBrandValue struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func (probeBrandValue) TableName() string { return "brands_value" }

type probePhoneValueRelation struct {
	ID      uint `gorm:"primarykey"`
	Name    string
	BrandID uint
	Brand   probeBrandValue // value, not pointer
}

func (probePhoneValueRelation) TableName() string { return "phones_value" }

// TestValueRelationNoMatchStaysZero confirms gorm leaves a value-type relation
// as the zero struct (not an error, not a sentinel) when a preload finds no
// match. The scanner's Mount is therefore a safe no-op for value relations on
// an all-NULL joined block.
func TestValueRelationNoMatchStaysZero(t *testing.T) {
	t.Parallel()

	db := openProbeDB(t)
	require.NoError(t, db.AutoMigrate(&probeBrandValue{}, &probePhoneValueRelation{}))

	require.NoError(t, db.Create(&probePhoneValueRelation{Name: "no_brand", BrandID: 999}).Error)

	var phone probePhoneValueRelation
	require.NoError(t, db.Preload("Brand").First(&phone).Error)

	require.Zero(t, phone.Brand.ID, "value relation must stay zero on no match")
	require.Empty(t, phone.Brand.Name, "value relation must stay zero on no match")
}
