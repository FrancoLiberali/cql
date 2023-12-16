package hasone

import "github.com/ditrit/badaas/orm"

type Country struct {
	orm.UUIDModel

	Capital City // Country HasOne City (Country 1 -> 1 City)
}

type City struct {
	orm.UUIDModel

	CountryID orm.UUID // Country HasOne City (Country 1 -> 1 City)
}
