package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
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
	Use:     "list <languages|stars|tags|trending>",
	Aliases: []string{"ls"},
	Short:   "List languages, stars, trending, or tags",
	Long:    "List languages, stars, trending, or tags that match your specified criteria.",
	Example: fmt.Sprintf("  %s list languages\n  %s list stars -t vim", config.ProgramName, config.ProgramName),
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
	fatalOnError(err)

	languages, err := model.FindLanguages(db)
	fatalOnError(err)

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
		stars, err = model.FindStarsByLanguage(db, options.language)
	} else if options.tag != "" {
		tag, _, err := model.FindOrCreateTagByName(db, options.tag)
		if err != nil {
			output.Fatal(err.Error())
		}
		err = tag.LoadStars(db)
		if err != nil {
			output.Fatal(err.Error())
		}
		stars, err = tag.Stars, nil
	} else {
		stars, err = model.FindStars(db)
	}

	if err != nil {
		output.Fatal(err.Error())
	}
	if stars != nil {
		for _, star := range stars {
			output.StarLine(&star)
		}
	}
}

func listTags() {
	output := getOutput()

	db, err := getDatabase()
	if err != nil {
		output.Fatal(err.Error())
	}

	tags, err := model.FindTags(db)
	if err != nil {
		output.Error(err.Error())
	} else {
		for _, tag := range tags {
			output.Info(tag.Name)
		}
	}
}

func listTrending() {
	getOutput().Info("Listing trending")
}

func init() {
	RootCmd.AddCommand(ListCmd)
}
