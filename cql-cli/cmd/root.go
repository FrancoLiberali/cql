package cmd

import (
	"github.com/ditrit/verdeter"

	"github.com/FrancoLiberali/cql/cql-cli/cmd/gen"
	"github.com/FrancoLiberali/cql/cql-cli/cmd/log"
	"github.com/FrancoLiberali/cql/cql-cli/cmd/version"
)

// TODO cli can be simplified, it is just a command
// rootCmd represents the base command when called without any subcommands
var rootCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:     "cql-cli",
	Short:   "the cql command line client",
	Long:    `cql-cli is the command line tool that makes it possible to use cql in your project.`,
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
