package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	jirapm "github.com/glestaris/issuez"
	"github.com/glestaris/issuez/domain"
)

var jiraProjectKey string

func init() {
	importCmd.PersistentFlags().StringVarP(
		&jiraProjectKey, "project-key", "p", "", "JIRA project key",
	)
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import <Path to Markdown file>",
	Short: "Imports markdown file as JIRA issues",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		markdownFilePath := args[0]
		markdownFile, err := os.Open(markdownFilePath)
		if err != nil {
			fmt.Printf(
				"Failed to open markdown file '%s': %s\n", markdownFilePath,
				err,
			)
			os.Exit(1)
		}
		defer markdownFile.Close()

		issues, err := ParseImportFile(markdownFile)
		if err != nil {
			fmt.Printf("Failed to parse markdown file: %s\n", err)
			os.Exit(1)
		}
		if len(issues) == 0 {
			fmt.Println("No issues were found")
			os.Exit(0)
		} else {
			fmt.Printf("Found %d issues in the markdown file\n", len(issues))
		}

		trackerService, err := jirapm.NewTrackerService(domain.Tracker{
			Type: "jira",
			Config: map[string]string{
				"apiHost":     jiraAPIHost,
				"apiUsername": jiraAPIUsername,
				"apiToken":    jiraAPIToken,
				"projectKey":  jiraProjectKey,
			},
		})
		if err != nil {
			fmt.Printf("Failed to initalise tracker service: %s\n", err)
			os.Exit(1)
		}

		err = trackerService.ImportIssues(issues)
		if err != nil {
			fmt.Printf("Failed to import issues: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Imported issues:\n")
		for _, issue := range issues {
			// - Task (TEST-124): Subject
			if issue.ID == "" {
				fmt.Printf("- %s (FAILED): %s\n", issue.Type, issue.Title)
			} else {
				fmt.Printf("- %s (%s): %s\n", issue.Type, issue.ID, issue.Title)

			}
		}
	},
}
