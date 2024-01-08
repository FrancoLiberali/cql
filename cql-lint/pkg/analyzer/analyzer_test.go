package analyzer_test

import (
	"testing"

	"github.com/FrancoLiberali/cql/cql-lint/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "a")
}
