package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/lucmski/limo/config"
	"github.com/lucmski/limo/model"
	"github.com/lucmski/limo/service"
	"github.com/spf13/cobra"
)

// UpdateCmd updates your stars from a remote service
var UpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update stars from a service",
	Long:    "Update your local database with your stars from the service specified by [--service] (default: github).",
	Example: fmt.Sprintf("  %s update", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Get configuration
		cfg, err := getConfiguration()
		fatalOnError(err)

		// Get the database
		db, err := getDatabase()
		fatalOnError(err)

		// Get the search index
		index, err := getIndex()
		fatalOnError(err)

		// Get the specified service
		svc, err := getService("")
		fatalOnError(err)

		// Get the database record for the specified service
		serviceName := service.Name(svc)
		dbSvc, _, err := model.FindOrCreateServiceByName(db, serviceName)
		fatalOnError(err)

		startTime := time.Now()

		// Create a channel to receive stars, since service can page
		starChan := make(chan *model.StarResult, 20)

		// Get the stars for the authenticated user
		go svc.GetStars(ctx, starChan, cfg.GetService(serviceName).Token, "")

		output := getOutput()

		totalCreated, totalUpdated, totalErrors := 0, 0, 0

		for starResult := range starChan {
			if starResult.Error != nil {
				totalErrors++
				output.Error(starResult.Error.Error())
			} else {
				created, err := model.CreateOrUpdateStar(db, starResult.Star, dbSvc)
				if err != nil {
					totalErrors++
					output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
				} else {
					if created {
						totalCreated++
					} else {
						totalUpdated++
					}
					err = starResult.Star.Index(index, db)
					if err != nil {
						totalErrors++
						output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
					}
					output.Tick("", "Updating "+*starResult.Star.FullName)
				}
			}
		}

		if totalCreated > 0 || totalUpdated > 0 {
			dbSvc.LastSuccess = startTime
			fatalOnError(db.Save(dbSvc).Error)
		}

		output.Info(fmt.Sprintf("\nCreated: %d; Updated: %d; Errors: %d", totalCreated, totalUpdated, totalErrors))
	},
}

func init() {
	RootCmd.AddCommand(UpdateCmd)
}
