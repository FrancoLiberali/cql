package testutils

import (
	"io/fs"
	"log"
	"os"
	"testing"
)

func CheckFileNotExists(t *testing.T, name string) {
	_, err := os.Stat(name)
	if err == nil {
		t.Error(err, "Should have been an error")
	}
}

func CheckFileExists(t *testing.T, name string) fs.FileInfo {
	stat, err := os.Stat(name)
	if err != nil {
		t.Error(err)
	}

	return stat
}

func RemoveFile(name string) {
	err := os.RemoveAll(name)
	if err != nil {
		log.Fatal(err)
	}
}
