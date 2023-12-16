package main

import (
	"github.com/Masterminds/semver/v3"

	"github.com/ditrit/badaas"
	"github.com/ditrit/badaas/controllers"
)

func main() {
	badaas.BaDaaS.AddModules(
		controllers.InfoControllerModule,
		controllers.AuthControllerModule,
	).Provide(
		NewAPIVersion,
	).Start()
}

func NewAPIVersion() *semver.Version {
	return semver.MustParse("0.0.0-unreleased")
}
