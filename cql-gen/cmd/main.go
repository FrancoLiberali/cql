package cmd

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

	"github.com/FrancoLiberali/cql/cql-gen/cmd/log"
	"github.com/FrancoLiberali/cql/cql-gen/cmd/version"
)

var GenConditionsCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:     "cql-gen",
	Short:   "Generate conditions to query your objects using cql",
	Long:    `cql-gen is the command line tool that makes it possible to use cql in your project, generating conditions to query your objects.`,
	Run:     GenerateConditions,
	Args:    cobra.MinimumNArgs(1),
	Version: version.Version,
})

const (
	DestPackageKey = "dest_package"
	cqlPath        = "github.com/FrancoLiberali/cql"
)

func init() {
	err := GenConditionsCmd.GKey(
		log.VerboseKey, verdeter.IsBool, "v",
		"Verbose logging",
	)
	if err != nil {
		panic(err)
	}

	GenConditionsCmd.SetDefault(log.VerboseKey, false)

	err = GenConditionsCmd.LKey(
		DestPackageKey, verdeter.IsStr, "d",
		"Destination package (not used if ran with go generate)",
	)
	if err != nil {
		panic(err)
	}
}

// GenConditionsCmd Run func
func GenerateConditions(_ *cobra.Command, args []string) {
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

// Generates a file with conditions for each cql model in the package
func generateConditionsForPkg(destPkg string, pkgPath string, pkg *packages.Package) {
	log.Logger.Infof("Generating conditions for types in package %q", pkg.Types.Name())

	relationGettersFile := NewFile(pkg.Types.Name(), filepath.Join(pkgPath, "cql.go"))

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
		// object is not a cql model, do not generate conditions
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
