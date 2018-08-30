package cmd

import (
	"fmt"
	"os"
	"strings"

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
	insecure bool
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
	flags.BoolVarP(&options.insecure, "insecure", "i", false, "skip certificate verification")
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
	o := output.ForName(options.output)
	oc, err := getConfiguration()
	if err == nil {
		o.Configure(oc.GetOutput(options.output))
	}
	return o
}

func getService() (service.Service, error) {
	return service.ForName(options.service, options.insecure)
}

func checkOneStar(name string, stars []model.Star) {
	o := getOutput()

	if len(stars) == 0 {
		o.Fatal(fmt.Sprintf("No stars match '%s'", name))
	}

	if len(stars) > 1 {
		o.Error(fmt.Sprintf("Star '%s' ambiguous:\n", name))
		for _, star := range stars {
			o.StarLine(&star)
		}
		o.Fatal("Narrow your search")
	}
}

func fatalOnError(err error) {
	if err != nil {
		getOutput().Fatal(err.Error())
	}
}

func parseServiceOwnerRepo(values []string) (string, string, string) {
	// `values` can be:
	// * Full URL (e.g., https://github.com/hoop33/limo)
	// * Owner/Repo (e.g., hoop33/limo)
	// * Owner Repo (e.g., hoop33 limo)
	var serviceName, owner, repo string

	if len(values) == 1 {
		values = strings.Split(strings.TrimPrefix(values[0], "https://"), "/")
	}
	if len(values) == 3 {
		serviceName, values = values[0], values[1:]
	}
	if len(values) == 2 {
		owner, repo = values[0], values[1]
	}

	// Drop the TLD from the service name
	if n := strings.LastIndex(serviceName, "."); n != -1 {
		serviceName = serviceName[:n]
	}
	return serviceName, owner, repo
}
