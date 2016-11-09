package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/hoop33/limo/service"
	"github.com/spf13/cobra"
)

var union = false
var notTagged = false

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
	Example: fmt.Sprintf("  %s list languages\n  %s list stars -t vim\n  %s list stars -t cli -l go", config.ProgramName, config.ProgramName, config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			getOutput().Fatal("You must specify languages, stars, tags, or trending")
		}

		if fn, ok := listers[args[0]]; ok {
			fn(args)
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

	match := ""
	if len(args) > 1 {
		match = args[1]
	}

	var stars []model.Star
	if notTagged {
		stars, err = model.FindUntaggedStars(db, match)
	} else if options.language != "" && options.tag != "" {
		stars, err = model.FindStarsByLanguageAndOrTag(db, match, options.language, options.tag, union)
	} else if options.language != "" {
		stars, err = model.FindStarsByLanguage(db, match, options.language)
	} else if options.tag != "" {
		tag, err := model.FindTagByName(db, options.tag)
		fatalOnError(err)

		if tag == nil {
			output.Fatal(fmt.Sprintf("Tag '%s' not found", options.tag))
		}

		err = tag.LoadStars(db, match)
		fatalOnError(err)

		stars = tag.Stars
	} else {
		stars, err = model.FindStars(db, match)
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

	tags, err := model.FindTagsWithStarCount(db)
	if err != nil {
		output.Error(err.Error())
	} else {
		for _, tag := range tags {
			output.Tag(&tag)
		}
	}
}

func listTrending(args []string) {
	// Get configuration
	cfg, err := getConfiguration()
	fatalOnError(err)

	// Get the specified service
	svc, err := getService()
	fatalOnError(err)

	// Create a channel to receive trending, since service can page
	trendingChan := make(chan *model.StarResult, 20)

	// Get trending for the specified service
	go svc.GetTrending(trendingChan, cfg.GetService(service.Name(svc)).Token, options.language, options.verbose)

	output := getOutput()

	for starResult := range trendingChan {
		if starResult.Error != nil {
			output.Error(starResult.Error.Error())
		} else {
			output.StarLine(starResult.Star)
		}
	}
}

func init() {
	ListCmd.Flags().BoolVarP(&union, "union", "u", false, "Show stars matching any arguments")
	ListCmd.Flags().BoolVarP(&notTagged, "notTagged", "n", false, "Show stars without any tags")
	RootCmd.AddCommand(ListCmd)
}
