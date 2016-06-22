package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/service"
	"github.com/spf13/cobra"
)

// LoginCmd lets you log in
var LoginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Log in to a service",
	Long:    "Log in to the service specified by [--service] (default: github).",
	Example: fmt.Sprintf("  %s login", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the specified service and log in
		svc, err := getService()
		if err != nil {
			getOutput().Fatal(err.Error())
		}

		token, err := svc.Login()
		if err != nil {
			getOutput().Fatal(err.Error())
		}

		// Update configuration with token
		config, err := getConfiguration()
		if err != nil {
			getOutput().Fatal(err.Error())
		}

		config.GetService(service.Name(svc)).Token = token
		err = config.WriteConfig()
		if err != nil {
			getOutput().Fatal(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(LoginCmd)
}
