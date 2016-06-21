package cmd

import (
	"fmt"

	"github.com/hoop33/limo/model"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

// OpenCmd opens a star's URL in your browser
var OpenCmd = &cobra.Command{
	Use:   "open <star>",
	Short: "Open a star's URL",
	Long:  "Open a star's URL in your default browser",
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) == 0 {
			output.Fatal("You must specify a star")
		}

		db, err := getDatabase()
		if err != nil {
			output.Fatal(err.Error())
		}

		stars, err := model.FuzzyFindStarsWithName(db, args[0])
		if err != nil {
			output.Fatal(err.Error())
		}

		if len(stars) == 0 {
			output.Fatal(fmt.Sprintf("No stars match '%s'", args[0]))
		}

		if len(stars) > 1 {
			output.Error(fmt.Sprintf("Star '%s' ambiguous:\n", args[0]))
			for _, star := range stars {
				output.StarLine(&star)
			}
			output.Fatal("Narrow your search")
		}

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
