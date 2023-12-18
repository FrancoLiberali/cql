package hasone

import (
	"github.com/ditrit/badaas/orm/model"
)

type Country struct {
	model.UUIDModel

	Capital City // Country HasOne City (Country 1 -> 1 City)
}

type City struct {
	model.UUIDModel

	Country   *Country
	CountryID model.UUID // Country HasOne City (Country 1 -> 1 City)
}
