package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
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
		// Get configuration
		cfg, err := getConfiguration()
		if err != nil {
			getOutput().Fatal(err.Error())
		}

		// Get the database
		db, err := config.InitDB(cfg.DatabasePath)
		if err != nil {
			getOutput().Fatal(err.Error())
		}

		// Get the specified service
		svc, err := getService()
		if err != nil {
			getOutput().Fatal(err.Error())
		}

		// Get the database record for the specified service
		serviceName := service.Name(svc)
		dbSvc, err := model.GetOrCreateService(db, serviceName)
		if err != nil {
			getOutput().Fatal(err.Error())
		}

		// Create a channel to receive stars, since service can page
		starChan := make(chan *model.StarResult, 20)

		// Get the stars for the authenticated user
		go svc.GetStars(starChan, cfg.GetService(serviceName).Token, "")

		totalCreated, totalUpdated, totalErrors := 0, 0, 0

		for starResult := range starChan {
			if starResult.Error != nil {
				totalErrors++
				getOutput().Error(starResult.Error.Error())
			} else {
				created, err := model.CreateOrUpdateStar(db, starResult.Star, dbSvc)
				if err != nil {
					totalErrors++
					getOutput().Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
				} else {
					if created {
						totalCreated++
					} else {
						totalUpdated++
					}
					getOutput().Tick()
				}
			}
		}
		getOutput().Info(fmt.Sprintf("\nCreated: %d; Updated: %d; Errors: %d", totalCreated, totalUpdated, totalErrors))
	},
}

func init() {
	RootCmd.AddCommand(UpdateCmd)
}
