package cmd

import (
	"github.com/ditrit/verdeter"

	"github.com/FrancoLiberali/cql/cql-cli/cmd/gen"
	"github.com/FrancoLiberali/cql/cql-cli/cmd/log"
	"github.com/FrancoLiberali/cql/cql-cli/cmd/version"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:     "badaas-cli",
	Short:   "the badaas command line client",
	Long:    `badaas-cli is the command line tool that makes it possible to configure and run your badaas applications easily.`,
	Version: version.Version,
})

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddSubCommand(gen.GenCmd)

	err := rootCmd.GKey(
		log.VerboseKey, verdeter.IsBool, "v",
		"Verbose logging",
	)
	if err != nil {
		panic(err)
	}

	rootCmd.SetDefault(log.VerboseKey, false)
}
