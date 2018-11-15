package cmd

import (
	"context"
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/hoop33/limo/service"
	"github.com/spf13/cobra"
)

var deleters = map[string]func([]string){
	"star": deleteStar,
	"tag":  deleteTag,
}

// DeleteCmd renames a tag
var DeleteCmd = &cobra.Command{
	Use:     "delete <star|tag> <name>...",
	Aliases: []string{"rm"},
	Short:   "Delete stars or tags",
	Long:    "Delete stars or tags. Deleting a tag removes it from your local database. Deleting a star unstars the repository on the specified service.",
	Example: fmt.Sprintf("  %s delete tag frameworks\n  %s delete star https://github.com/hoop33/limo", config.ProgramName, config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			getOutput().Fatal("You must specify star or tag and values")
		}

		if fn, ok := deleters[args[0]]; ok {
			fn(args[1:])
		} else {
			getOutput().Fatal(fmt.Sprintf("'%s' not valid", args[0]))
		}
	},
}

func deleteStar(values []string) {
	// Get configuration
	cfg, err := getConfiguration()
	fatalOnError(err)

	sn, owner, repo := parseServiceOwnerRepo(values)
	if owner == "" || repo == "" {
		getOutput().Fatal("You must specify a valid git URL, owner/repo, or owner repo")
	}

	svc, err := getService(sn)
	fatalOnError(err)
	serviceName := service.Name(svc)

	star, err := svc.DeleteStar(context.Background(), cfg.GetService(serviceName).Token, owner, repo)
	fatalOnError(err)

	// Get the database
	db, err := getDatabase()
	fatalOnError(err)

	dbSvc, _, err := model.FindOrCreateServiceByName(db, serviceName)
	fatalOnError(err)

	dbStar, err := model.FindStarByRemoteIDAndService(db, star.RemoteID, dbSvc)
	fatalOnError(err)

	err = dbStar.Delete(db)
	fatalOnError(err)

	getOutput().Info("Deleted star")
}

func deleteTag(values []string) {
	output := getOutput()

	db, err := getDatabase()
	fatalOnError(err)

	for _, value := range values {
		tag, err := model.FindTagByName(db, value)
		if err != nil {
			output.Error(err.Error())
		} else {
			if tag == nil {
				output.Error(fmt.Sprintf("Tag '%s' not found", value))
			} else {
				err = tag.Delete(db)
				if err != nil {
					output.Error(err.Error())
				} else {
					output.Info(fmt.Sprintf("Deleted tag '%s'", tag.Name))
				}
			}
		}
	}
}

func init() {
	RootCmd.AddCommand(DeleteCmd)
}
