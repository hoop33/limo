package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/hoop33/limo/service"
	"github.com/spf13/cobra"
)

var del = false

// PruneCmd prunes local stars that are no longer starred on a remote service
var PruneCmd = &cobra.Command{
	Use:     "prune",
	Short:   "Prune unstarred repositories",
	Long:    "Prune from your local database any repositories you no longer have starred on [--service] (default: github).",
	Example: fmt.Sprintf("  %s prune", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		db, err := getDatabase()
		fatalOnError(err)

		svc, err := getService("")
		fatalOnError(err)

		serviceName := service.Name(svc)
		dbSvc, _, err := model.FindOrCreateServiceByName(db, serviceName)
		fatalOnError(err)

		prunable, err := model.FindPrunableStars(db, dbSvc)
		fatalOnError(err)

		for _, star := range prunable {
			output.StarLine(&star)
			if del {
				fatalOnError(star.Delete(db))
			}
		}
	},
}

func init() {
	PruneCmd.Flags().BoolVarP(&del, "delete", "d", false, "Actually delete from your local database (default: display what would be deleted)")
	RootCmd.AddCommand(PruneCmd)
}
