package cmd

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDockerFilesCreateFilesWhenDestinationFolderNotExists(t *testing.T) {
	generateDockerFiles(nil, nil)
	checkFilesExist(t)
	teardown()
}

func TestGenerateDockerFilesOverwriteFilesWhenDestinationFolderExists(t *testing.T) {
	destDir := filepath.Join("badaas", "docker", "db")
	err := os.MkdirAll(destDir, os.ModePerm)
	assert.Nil(t, err)

	destFile := filepath.Join(destDir, "docker-compose.yml")
	err = os.WriteFile(destFile, []byte("hello"), 0o0600)
	assert.Nil(t, err)

	generateDockerFiles(nil, nil)
	checkFilesExist(t)

	fileContent, err := os.ReadFile(destFile)
	assert.Nil(t, err)
	assert.NotEqual(t, string(fileContent), "hello")

	teardown()
}

func TestCopyDirCreatesDestinationDirIfItDoesNotExist(t *testing.T) {
	copyDir(
		filepath.Join("docker", "db"),
		filepath.Join("badaas", "docker", "db"),
	)
	checkDockerDBFilesExist(t)
	teardown()
}

func TestCopyDirCopyFilesIfTheDestinationFolderAlreadyExists(t *testing.T) {
	destDir := filepath.Join("badaas", "docker", "db")
	err := os.MkdirAll(destDir, os.ModePerm)
	assert.Nil(t, err)

	copyDir(filepath.Join("docker", "db"), destDir)
	checkDockerDBFilesExist(t)
	teardown()
}

func TestCopyDirPanicsIfStatOnDestDirIsNotPossible(t *testing.T) {
	assertPanic(t, func() {
		copyDir(filepath.Join("docker", "db"), "\000")
	}, "error running stat")
}

func TestCopyDirPanicsIfDestDirCreationFails(t *testing.T) {
	assertPanic(t, func() {
		copyDir(filepath.Join("docker", "db"), "")
	}, "error creating directory")
}

func TestCopyDirPanicsIfReadOnEmbedFileSystemIsNotPossible(t *testing.T) {
	assertPanic(t, func() {
		copyDir("not_exists", filepath.Join("badaas", "docker", "db"))
	}, "error reading source directory")
}

func TestCopyDirPanicsIfDestDirIsNotADirectory(t *testing.T) {
	err := os.WriteFile("file.txt", []byte("hello"), 0o0600)
	assert.Nil(t, err)

	assertPanic(t, func() {
		copyDir(filepath.Join("docker", "db"), "file.txt")
	}, "destination path file.txt is not a directory")
}

func TestCopyFilePanicsWhenDestPathDoesNotExist(t *testing.T) {
	assertPanic(t, func() {
		copyFile(
			filepath.Join("docker", "db", "docker-compose.yml"),
			filepath.Join("badaas", "docker", "db", "docker-compose.yml"),
		)
	}, "error writing on destination file")
}

func TestCopyFileWorksWhenDestPathAlreadyExistsButNotTheFile(t *testing.T) {
	destDir := filepath.Join("badaas", "docker", "db")
	err := os.MkdirAll(destDir, os.ModePerm)
	assert.Nil(t, err)

	copyFile(
		filepath.Join("docker", "db", "docker-compose.yml"),
		filepath.Join(destDir, "docker-compose.yml"),
	)
	checkDockerDBFilesExist(t)
	teardown()
}

func TestCopyFileWorksWhenDestPathAndFileAlready(t *testing.T) {
	destDir := filepath.Join("badaas", "docker", "db")
	err := os.MkdirAll(destDir, os.ModePerm)
	assert.Nil(t, err)
	destFile := filepath.Join(destDir, "docker-compose.yml")
	err = os.WriteFile(destFile, []byte("hello"), 0o0600)
	assert.Nil(t, err)

	copyFile(filepath.Join("docker", "db", "docker-compose.yml"), destFile)

	checkDockerDBFilesExist(t)
	fileContent, err := os.ReadFile(destFile)
	assert.Nil(t, err)
	assert.NotEqual(t, string(fileContent), "hello")

	teardown()
}

func TestCopyFilePanicsIfReadOnEmbedFileSystemIsNotPossible(t *testing.T) {
	assertPanic(t, func() {
		copyFile(
			filepath.Join("docker", "db", "not_exists"),
			filepath.Join("badaas", "docker", "db"),
		)
	}, "error reading source file")
}

func TestCopyFilePanicsIfDestPathIsADirectory(t *testing.T) {
	err := os.MkdirAll("badaas", os.ModePerm)
	assert.Nil(t, err)

	assertPanic(t, func() {
		copyFile(filepath.Join("docker", "db", "docker-compose.yml"), "badaas/")
	}, "error writing on destination file")

	teardown()
}

func TestCopyFilePanicsIfWriteOnDestPathIsNotPossible(t *testing.T) {
	assertPanic(t, func() {
		copyFile(filepath.Join("docker", "db", "docker-compose.yml"), "/badaas.txt")
	}, "permission denied")
}

func assertPanic(t *testing.T, functionShouldPanic func(), errorMessage string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The function did not panic")
		} else {
			err := r.(error)
			assert.ErrorContains(t, err, errorMessage)
		}
	}()
	functionShouldPanic()
}

func checkFilesExist(t *testing.T) {
	checkFileExists(t, ".dockerignore")
	checkFileExists(t, "Makefile")
	checkFileExists(t, filepath.Join("badaas", "config", "badaas.yml"))
	checkFileExists(t, filepath.Join("badaas", "docker", "api", "docker-compose.yml"))
	checkFileExists(t, filepath.Join("badaas", "docker", "api", "Dockerfile"))
	checkDockerDBFilesExist(t)
}

func checkDockerDBFilesExist(t *testing.T) {
	checkFileExists(t, filepath.Join("badaas", "docker", "db", "docker-compose.yml"))
}

func checkFileExists(t *testing.T, name string) {
	if _, err := os.Stat(name); err != nil {
		t.Error(err)
	}
}

func teardown() {
	remove(".dockerignore")
	remove("Makefile")
	remove("file.txt")
	remove("badaas")
}

func remove(name string) {
	err := os.RemoveAll(name)
	if err != nil {
		log.Fatal(err)
	}
}
