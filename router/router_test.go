package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	router := SetupRouter()
	assert.NotNil(t, router)
}
