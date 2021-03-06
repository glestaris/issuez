package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	jiraAPIHost     string
	jiraAPIUsername string
	jiraAPIToken    string
)

var rootCmd = &cobra.Command{
	Use:   "issuez",
    Short: "Importing Tickets to Jira from Markdown.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please use one of the subcommands.")
	},
	Version: "1.0.0-beta.01",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&jiraAPIHost, "api", "a", "", "JIRA host to use",
	)
	rootCmd.MarkPersistentFlagRequired("api")
	rootCmd.PersistentFlags().StringVarP(
		&jiraAPIUsername, "username", "u", "",
		"Username to use to connect to JIRA",
	)
	rootCmd.MarkPersistentFlagRequired("username")
	rootCmd.PersistentFlags().StringVarP(
		&jiraAPIToken, "token", "t", "",
		"API token to use to connect to JIRA",
	)
	rootCmd.MarkPersistentFlagRequired("token")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
