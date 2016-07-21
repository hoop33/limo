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

func TestDefaultIndexPathIsSetWhenConfigIsEmpty(t *testing.T) {
	rmdirConfig()
	config, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, config.IndexPath)
}

func TestCanSetDatabasePath(t *testing.T) {
	rmdirConfig()
	config, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}
	config.DatabasePath = "database-path-foo"
	err = config.WriteConfig()
	if err != nil {
		t.Fatal(err)
	}
	contents, err := ioutil.ReadFile(configFilePath())
	if err != nil {
		t.Fatal(err)
	}
	rmdirConfig()
	assert.True(t, strings.IndexAny(string(contents), "database-path-foo") >= 0)
}

func TestCanSetIndexPath(t *testing.T) {
	rmdirConfig()
	config, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}
	config.IndexPath = "index-path-foo"
	err = config.WriteConfig()
	if err != nil {
		t.Fatal(err)
	}
	contents, err := ioutil.ReadFile(configFilePath())
	if err != nil {
		t.Fatal(err)
	}
	rmdirConfig()
	assert.True(t, strings.IndexAny(string(contents), "index-path-foo") >= 0)
}

func TestGetServiceReturnsEmptyWhenServiceDoesNotExist(t *testing.T) {
	rmdirConfig()
	config, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	svcCfg := config.GetService("foo")
	assert.Equal(t, "", svcCfg.Token)
}

func mkdirConfig() {
	if err := os.MkdirAll(configDirectoryPath, 0700); err != nil {
		panic(err)
	}
}

func rmdirConfig() {
	if err := os.RemoveAll(configDirectoryPath); err != nil {
		panic(err)
	}
}
