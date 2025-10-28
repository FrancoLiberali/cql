package uintmodelwithtimestamp

import (
	"time"

	"github.com/FrancoLiberali/cql/model"
)

type UIntModelWithTimestamp struct {
	model.UIntModel

	CreatedAt time.Time
}
