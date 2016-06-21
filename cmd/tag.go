package cmd

import (
	"fmt"

	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

// TagCmd tags a star
var TagCmd = &cobra.Command{
	Use:   "tag <star> <tag>...",
	Short: "Tag <star> with <tag>s",
	Long:  "Tag the specified star with the specified tags, creating tags as necessary",
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) < 2 {
			output.Fatal("You must specify a star and at least one tag")
		}

		db, err := getDatabase()
		if err != nil {
			output.Fatal(err.Error())
		}

		stars, err := model.FuzzyFindStarsWithName(db, args[0])
		if err != nil {
			output.Fatal(err.Error())
		}

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
