package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/json-iterator/go"
	"github.com/k0kubun/pp"
	"github.com/lucmski/limo/config"
	"github.com/lucmski/limo/model"
	"github.com/lucmski/limo/service"
	"github.com/spf13/cobra"
	"gopkg.in/olivere/elastic.v6"
	// "github.com/satori/go.uuid"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// WebuiCmd updates your stars from a remote service
var WebuiCmd = &cobra.Command{
	Use:     "webui",
	Short:   "webui stars from a service",
	Long:    "webui your local database with your stars from the service specified by [--service] (default: github).",
	Example: fmt.Sprintf("  %s webui", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Get configuration
		cfg, err := getConfiguration()
		fatalOnError(err)

		// Get the database
		db, err := getDatabase()
		fatalOnError(err)

		// Get the search index
		// index, err := getIndex()
		// fatalOnError(err)

		// curl -s -XGET 'http://127.0.0.1:9200/_nodes/http?pretty=1
		// Create a client
		// elastic.SetURL(elasticHost)
		client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"), elastic.SetSniff(false))
		fatalOnError(err)
		pp.Println("elastic new client connected")

		// Create an index
		_, _ = client.CreateIndex("limo").Do(context.Background())
		//fatalOnError(err)
		pp.Println("elastic CreateIndex ok")

		// Get the specified service
		svc, err := getService("")
		fatalOnError(err)

		// Get the database record for the specified service
		serviceName := service.Name(svc)
		dbSvc, _, err := model.FindOrCreateServiceByName(db, serviceName)
		fatalOnError(err)

		startTime := time.Now()

		// Create a channel to receive stars, since service can page
		starChan := make(chan *model.StarResult, 100)

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

					// add to es index
					// find a gist
					// pp.Println(*starResult.Star)
					// tweet := Tweet{User: "olivere", Message: "Take Five"}
					docinfo := fmt.Sprintf("%d-%s", starResult.Star.ServiceID, starResult.Star.RemoteID)
					// pp.Println("docinfo:", docinfo)
					star := *starResult.Star
					b, err := json.Marshal(&star)
					fatalOnError(err)
					_, err = client.Index().
						Index("limo").
						Type("doc").
						Id(docinfo).
						BodyJson(string(b)).
						Refresh("wait_for").
						Do(context.Background())
					fatalOnError(err)
					output.Tick("", "Updating")
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
	RootCmd.AddCommand(WebuiCmd)
}
