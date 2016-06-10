package config

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	configDirectoryPath = "./tmp"
}

func TestDefaultDatabasePathIsSetWhenConfigIsEmpty(t *testing.T) {
	rmdirConfig()
	config, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, config.DatabasePath)
}

func TestCanSetDatabasePath(t *testing.T) {
	rmdirConfig()
	config, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}
	config.DatabasePath = "foo"
	err = config.WriteConfig()
	if err != nil {
		t.Fatal(err)
	}
	contents, err := ioutil.ReadFile(configFilePath())
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, strings.IndexAny(string(contents), "foo") >= 0)
}

func mkdirConfig() {
	err := os.MkdirAll(configDirectoryPath, 0700)
	if err != nil {
		panic(err)
	}
}

func rmdirConfig() {
	err := os.RemoveAll(configDirectoryPath)
	if err != nil {
		panic(err)
	}
}
