package cmd

import (
	"fmt"

	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

var listers = map[string]func(){
	"languages": listLanguages,
	"stars":     listStars,
	"tags":      listTags,
	"trending":  listTrending,
}

// ListCmd lists stars, tags, or trending
var ListCmd = &cobra.Command{
	Use:   "list <languages|stars|tags|trending>",
	Short: "List languages, stars, trending, or tags",
	Long:  "List languages, stars, trending, or tags that match your specified criteria",
	Run: func(cmd *cobra.Command, args []string) {
		var target string
		if len(args) == 0 {
			target = "stars"
		} else {
			target = args[0]
		}

		if fn, ok := listers[target]; ok {
			fn()
		} else {
			getOutput().Fatal(fmt.Sprintf("'%s' not valid", target))
		}
	},
}

func listLanguages() {
	getOutput().Info("Listing languages")
}

func listStars() {
	// Get configuration
	cfg, err := getConfiguration()
	if err != nil {
		getOutput().Fatal(err.Error())
	}

	// Get the database
	db, err := model.InitDB(cfg.DatabasePath, options.verbose)
	if err != nil {
		getOutput().Fatal(err.Error())
	}

	stars, err := model.FindStarsWithLanguageAndTag(db, options.language, options.tag)
	if err != nil {
		getOutput().Error(err.Error())
	} else if stars != nil {
		for _, star := range stars {
			getOutput().StarLine(&star)
		}
	}
}

func listTags() {
	getOutput().Info("Listing tags")
}

func listTrending() {
	getOutput().Info("Listing trending")
}

func init() {
	RootCmd.AddCommand(ListCmd)
}
