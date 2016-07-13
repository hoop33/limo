package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

var listers = map[string]func(args []string){
	"languages": listLanguages,
	"stars":     listStars,
	"tags":      listTags,
	"trending":  listTrending,
}

// ListCmd lists stars, tags, or trending
var ListCmd = &cobra.Command{
	Use:     "list <languages|stars|tags|trending>",
	Aliases: []string{"ls"},
	Short:   "List languages, stars, tags, or trending",
	Long:    "List languages, stars, tags, or trending that match your specified criteria.",
	Example: fmt.Sprintf("  %s list languages\n  %s list stars -t vim", config.ProgramName, config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			getOutput().Fatal("You must specify languages, stars, tags, or trending")
		}

		if fn, ok := listers[args[0]]; ok {
			fn(args[1:])
		} else {
			getOutput().Fatal(fmt.Sprintf("'%s' not valid", args[0]))
		}
	},
}

func listLanguages(args []string) {
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

func listStars(args []string) {
	output := getOutput()

	db, err := getDatabase()
	fatalOnError(err)

	var stars []model.Star

	if options.language != "" {
		stars, err = model.FindStarsByLanguage(db, options.language)
	} else if options.tag != "" {
		tag, err := model.FindTagByName(db, options.tag)
		fatalOnError(err)

		if tag == nil {
			output.Fatal(fmt.Sprintf("Tag '%s' not found", options.tag))
		}

		err = tag.LoadStars(db)
		fatalOnError(err)

		stars = tag.Stars
	} else {
		stars, err = model.FindStars(db)
	}

	fatalOnError(err)

	if stars != nil {
		for _, star := range stars {
			output.StarLine(&star)
		}
	}
}

func listTags(args []string) {
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

func listTrending(args []string) {
	getOutput().Info("Listing trending")
}

func init() {
	RootCmd.AddCommand(ListCmd)
}
