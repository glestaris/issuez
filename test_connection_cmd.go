package main

import (
	"fmt"
	"os"

	"github.com/glestaris/issuez/domain"
	"github.com/glestaris/issuez/tracker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testConnectionCmd)
}

var testConnectionCmd = &cobra.Command{
	Use:   "test-connection",
	Short: "Tests JIRA connection",
	Run: func(cmd *cobra.Command, args []string) {
		trackerService, err := tracker.NewTrackerService(domain.Tracker{
			Type: "jira",
			Config: map[string]string{
				"apiHost":     jiraAPIHost,
				"apiUsername": jiraAPIUsername,
				"apiToken":    jiraAPIToken,
			},
		})
		if err != nil {
			fmt.Printf(
				"Failed to initialise tracker service for JIRA: %s\n", err,
			)
			os.Exit(1)
		}

		if err := trackerService.TestConnection(); err != nil {
			fmt.Printf("Failed to connect to JIRA: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("OK")
	},
}
