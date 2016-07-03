package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/skratchdot/open-golang/open"
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

		// The page to open
		var page string

		if homepage && stars[0].Homepage != nil && *stars[0].Homepage != "" {
			page = *stars[0].Homepage
		} else if stars[0].URL != nil && *stars[0].URL != "" {
			page = *stars[0].URL
		} else {
			output.Fatal("No URL for star")
		}

		output.Info(fmt.Sprintf("Opening %s...", page))
		fatalOnError(open.Start(page))
	},
}

func init() {
	OpenCmd.Flags().BoolVarP(&homepage, "homepage", "H", false, "open home page instead of URL")
	RootCmd.AddCommand(OpenCmd)
}
