package uuidmodelwithtimestamp

import (
	"time"

	"github.com/FrancoLiberali/cql/model"
)

type UUIDModelWithTimestamp struct {
	model.UUIDModel

	CreatedAt time.Time
}
