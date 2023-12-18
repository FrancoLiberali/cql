package orm_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/orm"
)

func TestParseCorrectUUID(t *testing.T) {
	uuidString := uuid.New().String()
	uuid, err := orm.ParseUUID(uuidString)
	assert.Nil(t, err)
	assert.Equal(t, uuidString, uuid.String())
}

func TestParseIncorrectUUID(t *testing.T) {
	uid, err := orm.ParseUUID("not uuid")
	assert.Error(t, err)
	assert.Equal(t, orm.NilUUID, uid)
}
