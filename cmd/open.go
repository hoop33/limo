package cmd

import (
	"fmt"

	"github.com/lucmski/limo/config"
	"github.com/lucmski/limo/model"
	"github.com/spf13/cobra"
)

var homepage = false

// OpenCmd opens a star's URL in your browser
var OpenCmd = &cobra.Command{
	Use:     "open <star>",
	Short:   "Open a star's URL",
	Long:    "Open a star's URL in your default browser.",
	Example: fmt.Sprintf("  %s open limo", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) == 0 {
			output.Fatal("You must specify a star")
		}

		db, err := getDatabase()
		fatalOnError(err)

		stars, err := model.FuzzyFindStarsByName(db, args[0])
		fatalOnError(err)

		checkOneStar(args[0], stars)

		err = stars[0].OpenInBrowser(homepage)
		fatalOnError(err)
	},
}

func init() {
	OpenCmd.Flags().BoolVarP(&homepage, "homepage", "H", false, "open home page instead of URL")
	RootCmd.AddCommand(OpenCmd)
}
