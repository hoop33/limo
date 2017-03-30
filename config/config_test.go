package config

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"fmt"

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
	assert.True(t, strings.ContainsAny(string(contents), "database-path-foo"))
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
	assert.True(t, strings.ContainsAny(string(contents), "index-path-foo"))
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

func TestReadConfigFileReadsFileWhenExists(t *testing.T) {
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

	cfg2, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "foo", cfg2.DatabasePath)

	rmdirConfig()
}

func TestReadConfigDoesNotPanicForMalformedConfigurationFile(t *testing.T) {
	rmdirConfig()
	mkdirConfig()

	contents := "{this is not a yaml file}"
	err := ioutil.WriteFile(fmt.Sprintf("%s/limo.yaml", configDirectoryPath), []byte(contents), 0700)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, cfg.DatabasePath == "")

	rmdirConfig()
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
