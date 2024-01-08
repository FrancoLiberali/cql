package main

import (
	"flag"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/FrancoLiberali/cql/cql-lint/pkg/analyzer"
)

func main() {
	// Don't use it: just to not crash on -unsafeptr flag from go vet
	flag.Bool("unsafeptr", false, "")

	singlechecker.Main(analyzer.Analyzer)
}
