package hasone

import (
	"github.com/FrancoLiberali/cql/model"
)

type Country struct {
	model.UUIDModel

	Capital City // Country HasOne City (Country 1 -> 1 City)
}

type City struct {
	model.UUIDModel

	Country   *Country
	CountryID *model.UUID `gorm:"not null"` // Country HasOne City (Country 1 -> 1 City)
}
