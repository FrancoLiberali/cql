package models

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

type Company struct {
	model.UUIDModelWithTimestamps

	Name    string
	Sellers *[]Seller // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

func (m Company) Equal(other Company) bool {
	return m.ID == other.ID
}

type CompanyNoTimestamps struct {
	model.UUIDModel

	Name    string
	Sellers *[]SellerNoTimestamps // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

func (m CompanyNoTimestamps) Equal(other CompanyNoTimestamps) bool {
	return m.ID == other.ID
}

type MultiString []string

func (s *MultiString) Scan(src interface{}) error {
	switch typedSrc := src.(type) {
	case string:
		*s = strings.Split(typedSrc, ",")
		return nil
	case []byte:
		str := string(typedSrc)
		*s = strings.Split(str, ",")

		return nil
	default:
		return fmt.Errorf("failed to scan multistring field - source is not a string, is %T", src)
	}
}

func (s MultiString) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil //nolint:nilnil // is necessary
	}

	return strings.Join(s, ","), nil
}

func (MultiString) GormDataType() string {
	return "text"
}

func (MultiString) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlserver":
		return "varchar(255)"
	default:
		return "text"
	}
}

func (s MultiString) GetValue() MultiString {
	return s
}

func (s MultiString) ToSQL(_ *condition.CQLQuery) (string, []any, error) {
	return "", []any{s}, nil
}

type ToBeEmbedded struct {
	EmbeddedInt int
}

type ToBeGormEmbedded struct {
	Int int
}

type Product struct {
	model.UUIDModelWithTimestamps

	String      string `gorm:"column:string_something_else"`
	Int         int
	IntPointer  *int
	Float       float64
	NullFloat   sql.NullFloat64
	Bool        bool
	NullBool    sql.NullBool
	ByteArray   []byte
	MultiString MultiString
	ToBeEmbedded
	GormEmbedded ToBeGormEmbedded `gorm:"embedded;embeddedPrefix:gorm_embedded_"`
	String2      string
}

func (m Product) Equal(other Product) bool {
	return m.ID == other.ID
}

type ProductNoTimestamps struct {
	model.UUIDModel

	String      string `gorm:"column:string_something_else"`
	Int         int
	IntPointer  *int
	Float       float64
	NullFloat   sql.NullFloat64
	Bool        bool
	NullBool    sql.NullBool
	ByteArray   []byte
	MultiString MultiString
	ToBeEmbedded
	GormEmbedded ToBeGormEmbedded `gorm:"embedded;embeddedPrefix:gorm_embedded_"`
	String2      string
}

func (m ProductNoTimestamps) Equal(other ProductNoTimestamps) bool {
	return m.ID == other.ID
}

// AllTypes exercises the fast scanner across every supported scalar type —
// as a value and as a nullable pointer — so a round-trip test can prove each
// generated scan/assign path materializes the value correctly.
type AllTypes struct {
	model.UUIDModel

	ValInt     int
	ValInt8    int8
	ValInt16   int16
	ValInt32   int32
	ValInt64   int64
	ValUint    uint
	ValUint8   uint8
	ValUint16  uint16
	ValUint32  uint32
	ValUint64  uint64
	ValFloat32 float32
	ValFloat64 float64
	ValBool    bool
	ValString  string
	ValTime    time.Time

	PtrInt     *int
	PtrInt8    *int8
	PtrInt16   *int16
	PtrInt32   *int32
	PtrInt64   *int64
	PtrUint    *uint
	PtrUint8   *uint8
	PtrUint16  *uint16
	PtrUint32  *uint32
	PtrUint64  *uint64
	PtrFloat32 *float32
	PtrFloat64 *float64
	PtrBool    *bool
	PtrString  *string
	PtrTime    *time.Time

	// database/sql nullable wrappers used directly as fields — these
	// exercise the sql.NullInt16 / NullInt32 / NullByte scan pools that the
	// sized int/uint value fields (which scan via NullInt64) never touch.
	NullInt16 sql.NullInt16
	NullInt32 sql.NullInt32
	NullByte  sql.NullByte

	// named scalars (underlying int): stored as their underlying kind, so the
	// scanner reads them via NullInt64 and casts back to Color — value and
	// nullable-pointer variants.
	Favorite Color
	PtrColor *Color
}

// Color is a named scalar type: no custom Scan/Value, so gorm stores it as its
// underlying int. The fast scanner classifies it via its underlying kind and
// casts the scanned value back to Color, so a plain named scalar takes the
// fast path (no gorm fallback).
type Color int

const (
	ColorRed Color = iota + 1
	ColorGreen
	ColorBlue
)

// WithUnsupportedColumn holds a column the fast scanner genuinely can't
// classify: a []string persisted through gorm's json serializer. There's no
// static Scan/Value and the value lives in a JSON text column, so no scanner
// is generated and its queries fall back to gorm's reflective scan — the
// safety net exercised by TestUnsupportedColumnFallsBackToGorm.
type WithUnsupportedColumn struct {
	model.UUIDModel

	Name string
	Tags []string `gorm:"serializer:json"`
}

type University struct {
	model.UUIDModel

	Name string
}

func (m University) Equal(other University) bool {
	return m.ID == other.ID
}

type Seller struct {
	model.UUIDModel

	Name      string
	Company   *Company
	CompanyID *model.UUID // Company HasMany Sellers (Company 0..1 -> 0..* Seller)

	University   *University
	UniversityID *model.UUID
}

func (m Seller) Equal(other Seller) bool {
	return m.Name == other.Name
}

type SellerNoTimestamps struct {
	model.UUIDModel

	Name                  string
	CompanyNoTimestamps   *CompanyNoTimestamps
	CompanyNoTimestampsID *model.UUID // Company HasMany Sellers (Company 0..1 -> 0..* Seller)

	University   *University
	UniversityID *model.UUID
}

func (m SellerNoTimestamps) Equal(other SellerNoTimestamps) bool {
	return m.Name == other.Name
}

type Sale struct {
	model.UUIDModelWithTimestamps

	Code        int
	Description string

	// Sale belongsTo Product (Sale 0..* -> 1 Product)
	Product   Product
	ProductID model.UUID

	// Sale belongsTo Seller (Sale 0..* -> 0..1 Seller)
	Seller   *Seller
	SellerID *model.UUID
}

func (m Sale) Equal(other Sale) bool {
	return m.ID == other.ID
}

type SaleNoTimestamps struct {
	model.UUIDModel

	Code        int
	Description string

	// Sale belongsTo Product (Sale 0..* -> 1 Product)
	Product   ProductNoTimestamps
	ProductID model.UUID

	// Sale belongsTo Seller (Sale 0..* -> 0..1 Seller)
	Seller   *SellerNoTimestamps
	SellerID *model.UUID
}

func (m SaleNoTimestamps) Equal(other SaleNoTimestamps) bool {
	return m.ID == other.ID
}

type Country struct {
	model.UUIDModel

	Name    string
	Capital City // Country HasOne City (Country 1 -> 1 City)
}

type City struct {
	model.UUIDModel

	Name      string
	Country   *Country
	CountryID model.UUID // Country HasOne City (Country 1 -> 1 City)
}

func (m Country) Equal(other Country) bool {
	return m.Name == other.Name
}

func (m City) Equal(other City) bool {
	return m.Name == other.Name
}

type Person struct {
	model.UUIDModel

	Name string `gorm:"unique;type:VARCHAR(255)"`
}

func (m Person) TableName() string {
	return "persons_and_more_name"
}

type Bicycle struct {
	model.UUIDModel

	Name string
	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Owner     Person `gorm:"references:Name;foreignKey:OwnerName"`
	OwnerName string
}

func (m Bicycle) Equal(other Bicycle) bool {
	return m.Name == other.Name
}

type Brand struct {
	model.UIntModel

	Name string
}

func (m Brand) Equal(other Brand) bool {
	return m.Name == other.Name
}

type Phone struct {
	model.UIntModelWithTimestamps

	Name string
	// Phone belongsTo Brand (Phone 0..* -> 1 Brand)
	Brand   Brand
	BrandID uint
}

func (m Phone) Equal(other Phone) bool {
	return m.Name == other.Name
}

type PhoneNoTimestamps struct {
	model.UIntModel

	Name string
	// Phone belongsTo Brand (Phone 0..* -> 1 Brand)
	Brand   Brand
	BrandID uint
}

func (m PhoneNoTimestamps) Equal(other PhoneNoTimestamps) bool {
	return m.Name == other.Name
}

type ParentParent struct {
	model.UUIDModel

	Name   string
	Number int
}

func (m ParentParent) Equal(other ParentParent) bool {
	return m.ID == other.ID
}

type Parent1 struct {
	model.UUIDModel

	Name string

	ParentParent   ParentParent
	ParentParentID model.UUID
}

func (m Parent1) Equal(other Parent1) bool {
	return m.ID == other.ID
}

type Parent2 struct {
	model.UUIDModel

	ParentParent   ParentParent
	ParentParentID model.UUID
}

func (m Parent2) Equal(other Parent2) bool {
	return m.ID == other.ID
}

type Child struct {
	model.UUIDModel

	Name   string
	Number int

	Parent1   Parent1
	Parent1ID model.UUID

	Parent2   Parent2
	Parent2ID model.UUID
}

func (m Child) Equal(other Child) bool {
	return m.ID == other.ID
}
