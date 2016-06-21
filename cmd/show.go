package cmd

import (
	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

// ShowCmd shows the version
var ShowCmd = &cobra.Command{
	Use:   "show <star>",
	Short: "Show <star>",
	Long:  "Show details about the specified star",
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

		for _, star := range stars {
			output.Star(&star)
			output.Info("")
		}
	},
}

func init() {
	RootCmd.AddCommand(ShowCmd)
}
