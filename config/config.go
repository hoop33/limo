package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/cep21/xdgbasedir"

	"gopkg.in/yaml.v2"
)

var configDirectoryPath string

// ServiceConfig contains configuration information for a service
type ServiceConfig struct {
	Token string
	User  string
}

// OutputConfig sontains configuration information for an output
type OutputConfig struct {
	SpinnerIndex    int `yaml:"spinnerIndex"`
	SpinnerInterval int `yaml:"spinnerInterval"`
	SpinnerColor    string `yaml:"spinnerColor"`
}

// Config contains configuration information
type Config struct {
	DatabasePath string                    `yaml:"databasePath"`
	IndexPath    string                    `yaml:"indexPath"`
	Services     map[string]*ServiceConfig `yaml:"services"`
	Outputs      map[string]*OutputConfig  `yaml:"outputs"`
}

// GetService returns the configuration information for a service
func (config *Config) GetService(name string) *ServiceConfig {
	if config.Services == nil {
		config.Services = make(map[string]*ServiceConfig)
	}

	service := config.Services[name]
	if service == nil {
		service = &ServiceConfig{}
		config.Services[name] = service
	}
	return service
}

// GetOutput returns the configuration information for an output
func (config *Config) GetOutput(name string) *OutputConfig {
	if config.Outputs == nil {
		config.Outputs = make(map[string]*OutputConfig)
	}

	output := config.Outputs[name]
	if output == nil {
		output = &OutputConfig{}
		config.Outputs[name] = output
	}
	return output
}

// ReadConfig reads the configuration information
func ReadConfig() (*Config, error) {
	file := configFilePath()

	var config Config
	if _, err := os.Stat(file); err == nil {
		// Read and unmarshal file only if it exists
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(f, &config)
		if err != nil {
			return nil, err
		}
	}

	// Set default database path
	if config.DatabasePath == "" {
		config.DatabasePath = path.Join(configDirectoryPath, fmt.Sprintf("%s.db", ProgramName))
	}

	// Set default search index path
	if config.IndexPath == "" {
		config.IndexPath = path.Join(configDirectoryPath, fmt.Sprintf("%s.idx", ProgramName))
	}
	return &config, nil
}

// WriteConfig writes the configuration information
func (config *Config) WriteConfig() error {
	err := os.MkdirAll(configDirectoryPath, 0700)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFilePath(), data, 0600)
}

func configFilePath() string {
	return path.Join(configDirectoryPath, fmt.Sprintf("%s.yaml", ProgramName))
}

func init() {
	baseDir, err := xdgbasedir.ConfigHomeDirectory()
	if err != nil {
		log.Fatal("Can't find XDG BaseDirectory")
	} else {
		configDirectoryPath = path.Join(baseDir, ProgramName)
	}
}
