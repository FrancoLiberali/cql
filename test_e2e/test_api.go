package main

import (
	"github.com/Masterminds/semver/v3"

	"github.com/ditrit/badaas"
)

func main() {
	badaas.BaDaaS.AddModules(
		badaas.InfoModule,
		badaas.AuthModule,
	).Provide(
		NewAPIVersion,
	).Start()
}

func NewAPIVersion() *semver.Version {
	return semver.MustParse("0.0.0-unreleased")
}
