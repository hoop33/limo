package cmd

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lucmski/limo/config"
	"github.com/lucmski/limo/model"
	"github.com/spf13/cobra"
)

var counters = map[string]func(ctx context.Context){
	"languages": countLanguages,
	"stars":     countStars,
	"tags":      countTags,
}

// CountCmd counts languages, stars, or tags
var CountCmd = &cobra.Command{
	Use:     "count <languages|stars|tags>",
	Aliases: []string{"c"},
	Short:   "Count languages, stars, or tags",
	Long:    "Count languages, stars, or tags that match your specified criteria",
	Example: fmt.Sprintf("  %s count languages\n  %s count stars -t vim", config.ProgramName, config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var which string
		if len(args) == 0 {
			which = "stars"
		} else {
			which = args[0]
		}

		if fn, ok := counters[which]; ok {
			fn(ctx)
		} else {
			getOutput().Fatal(fmt.Sprintf("'%s' not valid", which))
		}
	},
}

func countLanguages(_ context.Context) {
	counter(model.CountLanguages)
}

func countStars(_ context.Context) {
	counter(model.CountStarsByLanguageAndTag, options.language, options.tag)
}

func countTags(_ context.Context) {
	counter(model.CountTags)
}

func counter(fn func(db *gorm.DB, args ...string) (int, error), args ...string) {
	output := getOutput()

	db, err := getDatabase()
	if err != nil {
		output.Fatal(err.Error())
	}

	count, err := fn(db, args...)
	if err != nil {
		output.Fatal(err.Error())
	}
	output.Info(fmt.Sprintf("%d", count))
}

func init() {
	RootCmd.AddCommand(CountCmd)
}
