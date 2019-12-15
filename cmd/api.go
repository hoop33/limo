package cmd

import (
	// "context"
	"fmt"
	"net/http"

	// "github.com/jinzhu/gorm"
	// "github.com/lucmski/limo/config"
	"github.com/lucmski/limo/model"
	"github.com/qor/admin"
	"github.com/spf13/cobra"
)

/*
var counters = map[string]func(ctx context.Context){
	"languages": countLanguages,
	"stars":     countStars,
	"tags":      countTags,
}
*/

// CountCmd counts languages, stars, or tags
var ApiCmd = &cobra.Command{
	Use:     "api",
	Aliases: []string{"a"},
	Short:   "API for languages, stars, or tags",
	Long:    "API for languages, stars, or tags",
	// Example: fmt.Sprintf("  %s count languages\n  %s count stars -t vim", config.ProgramName, config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		// ctx := context.Background()

		dbs, err := getDatabase()
		if err != nil {
			panic(err)
		}

		dbs.AutoMigrate(&model.Star{}, &model.Tag{})

		// Initalize
		Admin := admin.New(&admin.AdminConfig{DB: dbs})

		// Allow to use Admin to manage User, Product
		Admin.AddResource(&model.Star{})
		Admin.AddResource(&model.Tag{})

		// initalize an HTTP request multiplexer
		mux := http.NewServeMux()

		// Mount admin interface to mux
		Admin.MountTo("/admin", mux)

		fmt.Println("Listening on: 9000")
		http.ListenAndServe(":9000", mux)

	},
}

func init() {
	RootCmd.AddCommand(ApiCmd)
}
