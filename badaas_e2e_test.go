package main

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
	"net/http"
	"os"
	"testing"
)

type TestContext struct {
	statusCode int
	json       map[string]interface{}
	httpClient *http.Client
}

var opts = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	t := &TestContext{}
	t.httpClient = &http.Client{}

	ctx.Step(`^I request "(.+)"$`, t.requestGET)
	ctx.Step(`^I expect status code is "(\d+)"$`, t.assertStatusCode)
	ctx.Step(`^I expect response field "(.+)" is "(.+)"$`, t.assertResponseFieldIsEquals)
}
