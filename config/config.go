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

var ConfigDirectoryPath string

// ServiceConfig contains configuration information for a service
type ServiceConfig struct {
	Token string
	User  string
}

// OutputConfig sontains configuration information for an output
type OutputConfig struct {
	SpinnerIndex    int    `yaml:"spinnerIndex"`
	SpinnerInterval int    `yaml:"spinnerInterval"`
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
func (cfg *Config) GetService(name string) *ServiceConfig {
	if cfg.Services == nil {
		cfg.Services = make(map[string]*ServiceConfig)
	}

	service := cfg.Services[name]
	if service == nil {
		service = &ServiceConfig{}
		cfg.Services[name] = service
	}
	return service
}

// GetOutput returns the configuration information for an output
func (cfg *Config) GetOutput(name string) *OutputConfig {
	if cfg.Outputs == nil {
		cfg.Outputs = make(map[string]*OutputConfig)
	}

	output := cfg.Outputs[name]
	if output == nil {
		output = &OutputConfig{}
		cfg.Outputs[name] = output
	}
	return output
}

// ReadConfig reads the configuration information
func ReadConfig() (*Config, error) {
	file := configFilePath()

	var cfg Config
	if _, err := os.Stat(file); err == nil {
		// Read and unmarshal file only if it exists
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(f, &cfg)
		if err != nil {
			return nil, err
		}
	}

	// Set default database path
	if cfg.DatabasePath == "" {
		cfg.DatabasePath = path.Join(ConfigDirectoryPath, fmt.Sprintf("%s.db", ProgramName))
	}

	// Set default search index path
	if cfg.IndexPath == "" {
		cfg.IndexPath = path.Join(ConfigDirectoryPath, fmt.Sprintf("%s.idx", ProgramName))
	}
	return &cfg, nil
}

// WriteConfig writes the configuration information
func (cfg *Config) WriteConfig() error {
	err := os.MkdirAll(ConfigDirectoryPath, 0700)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFilePath(), data, 0600)
}

func configFilePath() string {
	return path.Join(ConfigDirectoryPath, fmt.Sprintf("%s.yaml", ProgramName))
}

func init() {
	baseDir, err := xdgbasedir.ConfigHomeDirectory()
	if err != nil {
		log.Fatal("Can't find XDG BaseDirectory")
	} else {
		ConfigDirectoryPath = path.Join(baseDir, ProgramName)
	}
}

func EnsureDir(path string) {
	d, err := os.Open(path)
	if err != nil {
		os.MkdirAll(path, os.FileMode(0755))
	}
	d.Close()
}
