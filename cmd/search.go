package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/spf13/cobra"
)

// SearchCmd does a full-text search
var SearchCmd = &cobra.Command{
	Use:     "search <search string>",
	Aliases: []string{"find", "q"},
	Short:   "Search stars",
	Long:    "Perform a full-text search on your stars",
	Example: fmt.Sprintf("  %s search robust", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(SearchCmd)
}
