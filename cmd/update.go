package cmd

import (
	"github.com/hoop33/limo/model"
	"github.com/hoop33/limo/service"
	"github.com/spf13/cobra"
)

// UpdateCmd lets you log in
var UpdateCmd = &cobra.Command{
	Use:   "update [--service service]",
	Short: "Update [service]",
	Long:  "Update your database from the service specified by [--service] (default: github)",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the specified service
		svc, err := getService()
		if err != nil {
			getOutput().Fatal(err.Error())
		}

		// Get configuration
		config, err := getConfiguration()
		if err != nil {
			getOutput().Fatal(err.Error())
		}

		// Create a channel to receive stars, since service can page
		starChan := make(chan *model.StarResult, 1)

		// Get the stars for the authenticated user
		go svc.GetStars(starChan, config.GetService(service.Name(svc)).Token, "")

		for starResult := range starChan {
			if starResult.Error != nil {
				getOutput().Error(starResult.Error.Error())
			} else {
				getOutput().Info(starResult.Star.String())
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(UpdateCmd)
}
