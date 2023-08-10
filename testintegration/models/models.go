package models

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ditrit/badaas/orm/model"
)

type Company struct {
	model.UUIDModel

	Name    string
	Sellers *[]Seller // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

func (m Company) Equal(other Company) bool {
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
		return nil, nil
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

type ToBeEmbedded struct {
	EmbeddedInt int
}

type ToBeGormEmbedded struct {
	Int int
}

type Product struct {
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
}

func (m Product) Equal(other Product) bool {
	return m.ID == other.ID
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

type Sale struct {
	model.UUIDModel

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

func (m Seller) Equal(other Seller) bool {
	return m.Name == other.Name
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
	model.UIntModel

	Name string
	// Phone belongsTo Brand (Phone 0..* -> 1 Brand)
	Brand   Brand
	BrandID uint
}

func (m Phone) Equal(other Phone) bool {
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
