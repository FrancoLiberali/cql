package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ditrit/verdeter"
)

// File system embed in the executable that will have the following files:
//
//go:embed docker/*
//go:embed config/*
var genEmbedFS embed.FS

// genCommand represents the badaas-cli gen command
var genCommand = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "gen",
	Short: "Generate files and configurations necessary to use BadAss",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
	Run:   generateDockerFiles,
})

// directory where the generated files will be saved
const destBadaasDir = "badaas"

func init() {
	rootCmd.AddSubCommand(genCommand)
}

// copies all docker and configurations related files from the embed file system to the destination folder
func generateDockerFiles(cmd *cobra.Command, args []string) {
	sourceDockerDir := "docker"

	copyDir(
		filepath.Join(sourceDockerDir, "db"),
		filepath.Join(destBadaasDir, "docker", "db"),
	)

	copyDir(
		filepath.Join(sourceDockerDir, "api"),
		filepath.Join(destBadaasDir, "docker", "api"),
	)

	copyFile(
		filepath.Join(sourceDockerDir, ".dockerignore"),
		".dockerignore",
	)

	copyFile(
		filepath.Join(sourceDockerDir, "Makefile"),
		"Makefile",
	)

	copyDir(
		"config",
		filepath.Join(destBadaasDir, "config"),
	)
}

// copies a file from the embed file system to the destination folder
func copyFile(sourcePath, destPath string) {
	fileContent, err := genEmbedFS.ReadFile(sourcePath)
	if err != nil {
		panic(fmt.Errorf("error reading source file %s: %w", sourcePath, err))
	}

	if err := os.WriteFile(destPath, fileContent, 0o0600); err != nil {
		panic(fmt.Errorf("error writing on destination file %s: %w", destPath, err))
	}
}

// copies a directory from the embed file system to the destination folder
func copyDir(sourceDir, destDir string) {
	files, err := genEmbedFS.ReadDir(sourceDir)
	if err != nil {
		panic(fmt.Errorf("error reading source directory %s: %w", sourceDir, err))
	}

	fileInfo, err := os.Stat(destDir)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Errorf("error running stat on %s: %w", destDir, err))
		}

		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			panic(fmt.Errorf("error creating directory %s: %w", destDir, err))
		}
	} else if !fileInfo.IsDir() {
		panic(fmt.Errorf("destination path %s is not a directory", destDir))
	}

	for _, file := range files {
		copyFile(
			filepath.Join(sourceDir, file.Name()),
			filepath.Join(destDir, file.Name()),
		)
	}
}
