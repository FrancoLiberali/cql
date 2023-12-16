package utils

import (
	"io/fs"
	"log"
	"os"
	"testing"
)

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
