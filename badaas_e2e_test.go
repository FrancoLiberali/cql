package main

import (
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
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
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	t.httpClient = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Duration(5 * time.Second),
		Jar:       jar,
	}

	ctx.Step(`^I request "(.+)"$`, t.requestGET)
	ctx.Step(`^I expect status code is "(\d+)"$`, t.assertStatusCode)
	ctx.Step(`^I expect response field "(.+)" is "(.+)"$`, t.assertResponseFieldIsEquals)
	ctx.Step(`^I request "(.+)" with method "(.+)" with json$`, t.requestWithJson)
}
