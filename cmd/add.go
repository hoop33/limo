package cmd

import "github.com/spf13/cobra"

// AddCmd adds stars and tags
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add stars or tags",
	Long:  "Add a star or a tag",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
