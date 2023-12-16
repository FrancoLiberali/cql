package gen

import (
	"github.com/ditrit/badaas-cli/cmd/gen/conditions"
	"github.com/ditrit/verdeter"
)

var GenCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "gen",
	Short: "Files and configurations generator",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
})

func init() {
	GenCmd.AddSubCommand(genDockerCmd)
	GenCmd.AddSubCommand(conditions.GenConditionsCmd)
}
