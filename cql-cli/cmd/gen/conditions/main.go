package conditions

import (
	"errors"
	"fmt"
	"go/types"
	"os"
	"path/filepath"

	"github.com/ditrit/verdeter"
	"github.com/ettle/strcase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/packages"

	"github.com/FrancoLiberali/cql/cql-cli/cmd/log"
)

var GenConditionsCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "conditions",
	Short: "Generate conditions to query your objects using badaas-orm",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
	Run:   generateConditions,
	Args:  cobra.MinimumNArgs(1),
})

const (
	DestPackageKey = "dest_package"
	cqlPath        = "github.com/FrancoLiberali/cql"
)

func init() {
	err := GenConditionsCmd.LKey(
		DestPackageKey, verdeter.IsStr, "d",
		"Destination package (not used if ran with go generate)",
	)
	if err != nil {
		panic(err)
	}
}

// GenConditionsCmd Run func
func generateConditions(_ *cobra.Command, args []string) {
	log.SetLevel()
	// Inspect package and use type checker to infer imported types
	pkgs := loadPackages(args)

	// Get the package of the file with go:generate comment or in command params
	destPkg := os.Getenv("GOPACKAGE")
	if destPkg == "" {
		destPkg = viper.GetString(DestPackageKey)
		if destPkg == "" {
			panic(errors.New("config --dest_package or use go generate"))
		}
	}

	// Generate conditions for each package
	for i, pkg := range pkgs {
		generateConditionsForPkg(destPkg, args[i], pkg)
	}
}

// Generates a file with conditions for each Badaas model in the package
func generateConditionsForPkg(destPkg string, pkgPath string, pkg *packages.Package) {
	log.Logger.Infof("Generating conditions for types in package %q", pkg.Types.Name())

	relationGettersFile := NewFile(pkg.Types.Name(), filepath.Join(pkgPath, "badaas-orm.go"))

	for _, name := range pkg.Types.Scope().Names() {
		object := getObject(pkg, name)
		if object != nil {
			generateConditionsForObject(destPkg, object)
			_ = NewRelationGettersGenerator(object).Into(relationGettersFile)
		}
	}

	err := relationGettersFile.Save()
	if err != nil {
		panic(err)
	}
}

func generateConditionsForObject(destPkg string, object types.Object) {
	file := NewFile(
		destPkg,
		strcase.ToSnake(object.Name())+"_conditions.go",
	)

	err := NewConditionsGenerator(object).Into(file)
	if err != nil {
		// object is not a Badaas model, do not generate conditions
		return
	}

	err = file.Save()
	if err != nil {
		panic(err)
	}
}

// Load package information from paths
func loadPackages(paths []string) []*packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes}

	pkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		panic(fmt.Errorf("loading packages for inspection: %w", err))
	}

	// print compilation errors of source packages
	packages.PrintErrors(pkgs)

	return pkgs
}

// Get object by name in the package
func getObject(pkg *packages.Package, name string) types.Object {
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		panic(fmt.Errorf("%s not found in declared types of %s",
			name, pkg))
	}

	// Generate only if it is a declared type
	object, ok := obj.(*types.TypeName)
	if !ok {
		return nil
	}

	return object
}
