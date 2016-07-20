package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

// SearchCmd does a full-text search
var SearchCmd = &cobra.Command{
	Use:     "search <search string>",
	Aliases: []string{"find", "query", "q"},
	Short:   "Search stars",
	Long:    "Perform a full-text search on your stars",
	Example: fmt.Sprintf("  %s search robust", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) == 0 {
			output.Fatal("You must specify a search string")
		}

		index, err := getIndex()
		fatalOnError(err)

		query := bleve.NewMatchQuery(strings.Join(args, " "))
		request := bleve.NewSearchRequest(query)
		results, err := index.Search(request)
		fatalOnError(err)

		db, err := getDatabase()
		fatalOnError(err)

		for _, hit := range results.Hits {
			ID, err := strconv.Atoi(hit.ID)
			if err != nil {
				output.Error(err.Error())
			} else {
				star, err := model.FindStarByID(db, uint(ID))
				if err != nil {
					output.Error(err.Error())
				} else {
					output.Inline(fmt.Sprintf("(%f) ", hit.Score))
					output.StarLine(star)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(SearchCmd)
}
