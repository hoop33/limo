package cmd

import (
	"fmt"

	"github.com/lucmski/limo/config"
	"github.com/lucmski/limo/model"
	"github.com/spf13/cobra"
)

// UntagCmd tags a star
var UntagCmd = &cobra.Command{
	Use:     "untag <star> [tag]...",
	Short:   "Untag a star",
	Long:    "Untag the star identified by <star> with the tags specified by [tag], or all if [tag] not specified.",
	Example: fmt.Sprintf("  %s untag limo gui", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) == 0 {
			output.Fatal("You must specify a star (and optionally a tag)")
		}

		db, err := getDatabase()
		fatalOnError(err)

		stars, err := model.FuzzyFindStarsByName(db, args[0])
		fatalOnError(err)

		checkOneStar(args[0], stars)

		output.StarLine(&stars[0])

		if len(args) == 1 {
			// Untag all
			fatalOnError(stars[0].RemoveAllTags(db))
			output.Info(fmt.Sprintf("Removed all tags"))
		} else {
			fatalOnError(stars[0].LoadTags(db))

			for _, tagName := range args[1:] {
				tag, err := model.FindTagByName(db, tagName)
				if err != nil {
					output.Error(err.Error())
				} else if tag == nil {
					output.Error(fmt.Sprintf("Tag '%s' does not exist", tagName))
				} else if !stars[0].HasTag(tag) {
					output.Error(fmt.Sprintf("'%s' isn't tagged with '%s'", *stars[0].FullName, tagName))
				} else {
					err = stars[0].RemoveTag(db, tag)
					if err != nil {
						output.Error(err.Error())
					} else {
						output.Info(fmt.Sprintf("Removed tag '%s'", tag.Name))
					}
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(UntagCmd)
}
