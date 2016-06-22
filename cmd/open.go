package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

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
		if err != nil {
			output.Fatal(err.Error())
		}

		stars, err := model.FuzzyFindStarsByName(db, args[0])
		if err != nil {
			output.Fatal(err.Error())
		}

		checkOneStar(args[0], stars)

		if stars[0].URL == nil || *stars[0].URL == "" {
			output.Fatal("No URL for star")
		}

		output.Info(fmt.Sprintf("Opening %s...", *stars[0].URL))
		err = open.Start(*stars[0].URL)
		if err != nil {
			output.Fatal(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(OpenCmd)
}
