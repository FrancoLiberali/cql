package router

import (
	"testing"
)

func TestSetupRouter(t *testing.T) {
	router := SetupRouter()

	if router == nil {
		t.Error("function SetupRouter should return instantiate router")
	}
}
