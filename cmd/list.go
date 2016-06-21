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
	output := getOutput()

	db, err := getDatabase()
	if err != nil {
		output.Fatal(err.Error())
	}

	languages, err := model.FindLanguages(db)
	if err != nil {
		output.Fatal(err.Error())
	}

	for _, language := range languages {
		if language != "" {
			output.Info(language)
		}
	}
}

func listStars() {
	output := getOutput()

	db, err := getDatabase()
	if err != nil {
		output.Fatal(err.Error())
	}

	var stars []model.Star

	if options.language != "" {
		stars, err = model.FindStarsWithLanguage(db, options.language)
	} else {
		stars, err = model.FindStars(db)
	}

	if err != nil {
		output.Error(err.Error())
	} else if stars != nil {
		for _, star := range stars {
			output.StarLine(&star)
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
