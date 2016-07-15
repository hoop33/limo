package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

var union = false
var intersection = false
var notTagged = false

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
	Example: fmt.Sprintf("  %s list languages\n  %s list stars -t vim\n  %s list stars -t cli -l go", config.ProgramName, config.ProgramName, config.ProgramName),
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
	fatalOnError(err)

	var stars []model.Star
	if notTagged {
		stars, err = model.FindUntaggedStars(db)
	} else if options.language != "" && options.tag != "" {
		stars, err = model.FindStarsByLanguageAndOrTag(db, options.language, options.tag, union)
	} else if options.language != "" {
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
	ListCmd.Flags().BoolVarP(&union, "union", "u", false, "Show stars matching any arguments")
	ListCmd.Flags().BoolVarP(&notTagged, "notTagged", "n", false, "Show stars without any tags")
	RootCmd.AddCommand(ListCmd)
}
