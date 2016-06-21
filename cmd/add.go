package cmd

import (
	"fmt"

	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

var adders = map[string]func([]string){
	"star": addStar,
	"tag":  addTag,
}

// AddCmd adds stars and tags
var AddCmd = &cobra.Command{
	Use:   "add <star|tag> values...",
	Short: "Add stars or tags",
	Long:  "Add a star or a tag",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			getOutput().Fatal("You must specify star or tag and values")
		}

		if fn, ok := adders[args[0]]; ok {
			fn(args[1:])
		} else {
			getOutput().Fatal(fmt.Sprintf("'%s' not valid", args[0]))
		}
	},
}

func addStar(values []string) {
}

func addTag(values []string) {
	output := getOutput()

	db, err := getDatabase()
	if err != nil {
		output.Fatal(err.Error())
	}

	for _, value := range values {
		tag, created, err := model.FindOrCreateTagByName(db, value)
		if err != nil {
			output.Error(err.Error())
		} else {
			if created {
				output.Info(fmt.Sprintf("Created tag '%s'", tag.Name))
			} else {
				output.Error(fmt.Sprintf("Tag '%s' already exists", tag.Name))
			}
		}
	}
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
