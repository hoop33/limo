package cmd

import "github.com/spf13/cobra"

// TagCmd tags a star
var TagCmd = &cobra.Command{
	Use: "tag",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(TagCmd)
}
