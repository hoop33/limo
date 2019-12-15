package cmd

import (
	"fmt"

	"github.com/lucmski/limo/config"
	"github.com/lucmski/limo/model"
	"github.com/spf13/cobra"
)

// TagCmd tags a star
var TagCmd = &cobra.Command{
	Use:     "tag <star> <tag>...",
	Short:   "Tag a star",
	Long:    "Tag the star identified by <star> with the tags specified by <tag>, creating tags as necessary.",
	Example: fmt.Sprintf("  %s tag limo git cli", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) < 2 {
			output.Fatal("You must specify a star and at least one tag")
		}

		db, err := getDatabase()
		fatalOnError(err)

		stars, err := model.FuzzyFindStarsByName(db, args[0])
		fatalOnError(err)

		checkOneStar(args[0], stars)

		output.StarLine(&stars[0])
		for _, tagName := range args[1:] {
			tag, _, err := model.FindOrCreateTagByName(db, tagName)
			if err != nil {
				output.Error(err.Error())
			} else {
				err = stars[0].AddTag(db, tag)
				if err != nil {
					output.Error(err.Error())
				} else {
					output.Info(fmt.Sprintf("Added tag '%s'", tag.Name))
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(TagCmd)
}
