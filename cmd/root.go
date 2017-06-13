package cmd

import (
	"fmt"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/hoop33/limo/output"
	"github.com/hoop33/limo/service"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

var configuration *config.Config
var db *gorm.DB
var index bleve.Index

var options struct {
	language string
	output   string
	service  string
	tag      string
	verbose  bool
}

// RootCmd is the root command for limo
var RootCmd = &cobra.Command{
	Use:   "limo",
	Short: "A CLI for managing starred repositories",
	Long: `limo allows you to manage your starred repositories on GitHub, GitLab, and Bitbucket.
You can tag, display, and search your starred repositories.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	flags := RootCmd.PersistentFlags()
	flags.StringVarP(&options.language, "language", "l", "", "language")
	flags.StringVarP(&options.output, "output", "o", "color", "output type")
	flags.StringVarP(&options.service, "service", "s", "github", "service")
	flags.StringVarP(&options.tag, "tag", "t", "", "tag")
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose output")
}

func getConfiguration() (*config.Config, error) {
	if configuration == nil {
		var err error
		if configuration, err = config.ReadConfig(); err != nil {
			return nil, err
		}
	}
	return configuration, nil
}

func getDatabase() (*gorm.DB, error) {
	if db == nil {
		cfg, err := getConfiguration()
		if err != nil {
			return nil, err
		}
		db, err = model.InitDB(cfg.DatabasePath, options.verbose)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func getIndex() (bleve.Index, error) {
	if index == nil {
		cfg, err := getConfiguration()
		if err != nil {
			return nil, err
		}
		index, err = model.InitIndex(cfg.IndexPath)
		if err != nil {
			return nil, err
		}
	}
	return index, nil
}

func getOutput() output.Output {
	output := output.ForName(options.output)
	oc, err := getConfiguration()
	if err == nil {
		output.Configure(oc.GetOutput(options.output))
	}
	return output
}

func getService() (service.Service, error) {
	return service.ForName(options.service)
}

func checkOneStar(name string, stars []model.Star) {
	output := getOutput()

	if len(stars) == 0 {
		output.Fatal(fmt.Sprintf("No stars match '%s'", name))
	}

	if len(stars) > 1 {
		output.Error(fmt.Sprintf("Star '%s' ambiguous:\n", name))
		for _, star := range stars {
			output.StarLine(&star)
		}
		output.Fatal("Narrow your search")
	}
}

func fatalOnError(err error) {
	if err != nil {
		getOutput().Fatal(err.Error())
	}
}
