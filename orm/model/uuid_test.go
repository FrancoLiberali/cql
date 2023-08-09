package model_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/orm/model"
)

func TestParseCorrectUUID(t *testing.T) {
	uuidString := uuid.New().String()
	uuid, err := model.ParseUUID(uuidString)
	assert.Nil(t, err)
	assert.Equal(t, uuidString, uuid.String())
}

func TestParseIncorrectUUID(t *testing.T) {
	uid, err := model.ParseUUID("not uuid")
	assert.Error(t, err)
	assert.Equal(t, model.NilUUID, uid)
}
