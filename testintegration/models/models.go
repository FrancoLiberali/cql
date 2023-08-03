package models

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/ditrit/badaas/orm"
)

type Employee struct {
	orm.UUIDModel

	Name   string
	Boss   *Employee // Self-Referential Has One (Employee 0..* -> 0..1 Employee)
	BossID *orm.UUID
}

func (m Employee) Equal(other Employee) bool {
	return m.Name == other.Name
}

type Company struct {
	orm.UUIDModel

	Name    string
	Sellers []Seller // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
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

type ToBeEmbedded struct {
	EmbeddedInt int
}

type ToBeGormEmbedded struct {
	Int int
}

type Product struct {
	orm.UUIDModel

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

type Seller struct {
	orm.UUIDModel

	Name      string
	Company   *Company
	CompanyID *orm.UUID // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

type Sale struct {
	orm.UUIDModel

	Code        int
	Description string

	// Sale belongsTo Product (Sale 0..* -> 1 Product)
	Product   Product
	ProductID orm.UUID

	// Sale HasOne Seller (Sale 0..* -> 0..1 Seller)
	Seller   *Seller
	SellerID *orm.UUID
}

func (m Sale) Equal(other Sale) bool {
	return m.ID == other.ID
}

func (m Seller) Equal(other Seller) bool {
	return m.Name == other.Name
}

type Country struct {
	orm.UUIDModel

	Name    string
	Capital City // Country HasOne City (Country 1 -> 1 City)
}

type City struct {
	orm.UUIDModel

	Name      string
	Country   *Country
	CountryID orm.UUID // Country HasOne City (Country 1 -> 1 City)
}

func (m Country) Equal(other Country) bool {
	return m.Name == other.Name
}

func (m City) Equal(other City) bool {
	return m.Name == other.Name
}

type Person struct {
	orm.UUIDModel

	Name string `gorm:"unique;type:VARCHAR(255)"`
}

func (m Person) TableName() string {
	return "persons_and_more_name"
}

type Bicycle struct {
	orm.UUIDModel

	Name string
	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Owner     Person `gorm:"references:Name;foreignKey:OwnerName"`
	OwnerName string
}

func (m Bicycle) Equal(other Bicycle) bool {
	return m.Name == other.Name
}

type Brand struct {
	orm.UIntModel

	Name string
}

func (m Brand) Equal(other Brand) bool {
	return m.Name == other.Name
}

type Phone struct {
	orm.UIntModel

	Name string
	// Phone belongsTo Brand (Phone 0..* -> 1 Brand)
	Brand   Brand
	BrandID uint
}

func (m Phone) Equal(other Phone) bool {
	return m.Name == other.Name
}

type ParentParent struct {
	orm.UUIDModel

	Name string
}

func (m ParentParent) Equal(other ParentParent) bool {
	return m.ID == other.ID
}

type Parent1 struct {
	orm.UUIDModel

	ParentParent   ParentParent
	ParentParentID orm.UUID
}

func (m Parent1) Equal(other Parent1) bool {
	return m.ID == other.ID
}

type Parent2 struct {
	orm.UUIDModel

	ParentParent   ParentParent
	ParentParentID orm.UUID
}

func (m Parent2) Equal(other Parent2) bool {
	return m.ID == other.ID
}

type Child struct {
	orm.UUIDModel

	Parent1   Parent1
	Parent1ID orm.UUID

	Parent2   Parent2
	Parent2ID orm.UUID
}

func (m Child) Equal(other Child) bool {
	return m.ID == other.ID
}
