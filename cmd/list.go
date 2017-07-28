package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/hoop33/entrevista"
	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/hoop33/limo/service"
	"github.com/spf13/cobra"
)

var any = false
var browse = false
var notTagged = false
var page = 1
var count = 1
var user = ""

var listers = map[string]func(ctx context.Context, args []string){
	"events":    listEvents,
	"languages": listLanguages,
	"stars":     listStars,
	"tags":      listTags,
	"trending":  listTrending,
}

// ListCmd lists stars, tags, or trending
var ListCmd = &cobra.Command{
	Use:     "list <events|languages|stars|tags|trending>",
	Aliases: []string{"ls"},
	Short:   "List events, languages, stars, tags, or trending",
	Long:    "List events, languages, stars, tags, or trending that match your specified criteria.",
	Example: fmt.Sprintf("  %s list events\n  %s list languages\n  %s list stars -t vim\n  %s list stars -t cli -l go", config.ProgramName, config.ProgramName, config.ProgramName, config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var which string
		if len(args) == 0 {
			which = "events"
		} else {
			which = args[0]
		}

		if fn, ok := listers[which]; ok {
			fn(ctx, args)
		} else {
			getOutput().Fatal(fmt.Sprintf("'%s' not valid", which))
		}
	},
}

func listEvents(ctx context.Context, args []string) {
	cfg, err := getConfiguration()
	fatalOnError(err)

	svc, err := getService()
	fatalOnError(err)

	if user == "" {
		user = cfg.GetService(service.Name(svc)).User
		if user == "" {
			var err error
			user, err = getUser()
			fatalOnError(err)
			cfg.GetService(service.Name(svc)).User = user
			fatalOnError(cfg.WriteConfig())
		}
	}

	eventChan := make(chan *model.EventResult, 20)

	go svc.GetEvents(ctx, eventChan, cfg.GetService(service.Name(svc)).Token, user, page, count)

	output := getOutput()

	for eventResult := range eventChan {
		if eventResult.Error != nil {
			output.Error(eventResult.Error.Error())
		} else {
			output.Event(eventResult.Event)
			if browse {
				err := eventResult.Event.OpenInBrowser()
				if err != nil {
					output.Error(err.Error())
				}
			}
		}
	}
}

func listLanguages(ctx context.Context, args []string) {
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

func listStars(ctx context.Context, args []string) {
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
		stars, err = model.FindStarsByLanguageAndOrTag(db, match, options.language, options.tag, any)
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
			if browse {
				err := star.OpenInBrowser(false)
				if err != nil {
					output.Error(err.Error())
				}
			}
		}
	}
}

func listTags(ctx context.Context, args []string) {
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

func listTrending(ctx context.Context, args []string) {
	// Get configuration
	cfg, err := getConfiguration()
	fatalOnError(err)

	// Get the specified service
	svc, err := getService()
	fatalOnError(err)

	// Create a channel to receive trending, since service can page
	trendingChan := make(chan *model.StarResult, 20)

	// Get trending for the specified service
	go svc.GetTrending(ctx, trendingChan, cfg.GetService(service.Name(svc)).Token, options.language, options.verbose)

	output := getOutput()

	for starResult := range trendingChan {
		if starResult.Error != nil {
			output.Error(starResult.Error.Error())
		} else {
			output.StarLine(starResult.Star)
			if browse {
				err := starResult.Star.OpenInBrowser(false)
				if err != nil {
					output.Error(err.Error())
				}
			}
		}
	}
}

func getUser() (string, error) {
	interview := entrevista.NewInterview()
	interview.ShowOutput = func(message string) {
		fmt.Print(color.GreenString(message))
	}
	interview.ShowError = func(message string) {
		color.Red(message)
	}
	interview.Questions = []entrevista.Question{
		{
			Key:      "user",
			Text:     "Enter your user ID",
			Required: true,
			Hidden:   false,
		},
	}

	answers, err := interview.Run()
	if err != nil {
		return "", err
	}
	return answers["user"].(string), nil
}

func init() {
	ListCmd.Flags().BoolVarP(&any, "any", "a", false, "Show stars matching any arguments")
	ListCmd.Flags().BoolVarP(&browse, "browse", "b", false, "Open listed items in your default browser")
	ListCmd.Flags().BoolVarP(&notTagged, "notTagged", "n", false, "Show stars without any tags")
	ListCmd.Flags().IntVarP(&page, "page", "p", 1, "First event page to list")
	ListCmd.Flags().IntVarP(&count, "count", "c", 1, "Count of event pages to list")
	ListCmd.Flags().StringVarP(&user, "user", "u", "", "User for event list")
	RootCmd.AddCommand(ListCmd)
}
