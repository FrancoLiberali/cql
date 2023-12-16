package cmd

import (
	"github.com/ditrit/verdeter"
)

// rootCommand represents the base command when called without any subcommands
var rootCommand = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "badaas-cli",
	Short: "the BadAas controller",
	Long:  `badaas-cli is the command line tool that makes it possible to configure and run a badaas application`,
})

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCommand.
func Execute() {
	rootCommand.Execute()
}
