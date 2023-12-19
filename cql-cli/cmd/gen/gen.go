package gen

import (
	"github.com/ditrit/verdeter"

	"github.com/FrancoLiberali/cql/cql-cli/cmd/gen/conditions"
)

var GenCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "gen",
	Short: "Files and configurations generator",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
})

func init() {
	GenCmd.AddSubCommand(conditions.GenConditionsCmd)
}
