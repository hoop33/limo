package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

// ShowCmd shows the version
var ShowCmd = &cobra.Command{
	Use:     "show <star>",
	Short:   "Show a star's details",
	Long:    "Show details about the star identified by <star>.",
	Example: fmt.Sprintf("  %s show limo", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) == 0 {
			output.Fatal("You must specify a star")
		}

		db, err := getDatabase()
		fatalOnError(err)

		stars, err := model.FuzzyFindStarsByName(db, args[0])
		fatalOnError(err)

		for _, star := range stars {
			err = star.LoadTags(db)
			if err != nil {
				output.Error(err.Error())
			} else {
				output.Star(&star)
				output.Info("")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(ShowCmd)
}
