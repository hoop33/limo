package cmd

import (
	"fmt"

	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

// RenameCmd renames a tag
var RenameCmd = &cobra.Command{
	Use:     "rename <tag> <name>",
	Aliases: []string{"mv"},
	Short:   "Rename a tag",
	Long:    "Rename tag <tag> to <name>",
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) < 2 {
			output.Fatal("You must specify a tag and a new name")
		}

		db, err := getDatabase()
		if err != nil {
			output.Fatal(err.Error())
		}

		tag, err := model.FindTagByName(db, args[0])
		if err != nil {
			output.Fatal(err.Error())
		}

		if tag == nil {
			output.Fatal(fmt.Sprintf("Tag '%s' not found", args[0]))
		}

		err = tag.Rename(db, args[1])
		if err != nil {
			output.Fatal(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(RenameCmd)
}
