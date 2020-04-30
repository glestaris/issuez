package jirapm

import (
	"log"

	"github.com/glestaris/issuez/domain"
	"github.com/glestaris/issuez/jira"
)

// APP layer
//  Test using integration tests

type jiraTrackerService struct {
	jiraClient *jira.Client
	projectKey string
}

func newJiraTrackerService(
	apiHost string, apiUsername string, apiToken string, projectKey string,
) TrackerService {
	jiraClient := jira.NewJiraClient(apiHost, apiUsername, apiToken, nil)
	return &jiraTrackerService{
		jiraClient: jiraClient,
		projectKey: projectKey,
	}
}

func (j *jiraTrackerService) ImportIssues(domainIssues []*domain.Issue) error {
	jiraIssues := make([]*jira.Issue, len(domainIssues))
	for i, domainIssue := range domainIssues {
		jiraIssue := &jira.Issue{}

		// map project
		jiraIssue.ProjectKey = j.projectKey

		// map issue type
		jiraIssue.Type = jira.IssueTypeStory
		switch domainIssue.Type {
		case domain.IssueTypeBug:
			jiraIssue.Type = jira.IssueTypeBug
		case domain.IssueTypeChore:
			jiraIssue.Type = jira.IssueTypeTask
		}

		// map title
		jiraIssue.Summary = domainIssue.Title

		// map description
		jiraDescriptionDoc, err := mapDocument(domainIssue.Description)
		if err != nil {
			log.Printf("Failed to map description for issue '%s': %s",
				domainIssue.Title,
				err)
			continue
		}
		jiraIssue.Description = jiraDescriptionDoc

		// map epic
		if domainIssue.Epic != nil {
			jiraIssue.EpicKey = domainIssue.Epic.ID
		}

		// map labels
		for _, domainLabel := range domainIssue.Labels {
			jiraIssue.Labels = append(jiraIssue.Labels, domainLabel.Label)
		}

		jiraIssues[i] = jiraIssue
	}

	resp, err := j.jiraClient.ImportIssues(jiraIssues)
	if err != nil {
		return err
	}

	for i, entry := range resp {
		if entry.Err != nil {
			log.Printf("Failed to import issue '%s': %s",
				domainIssues[i].Title,
				err)
			continue
		}
		domainIssues[i].ID = entry.NewIssueKey
	}

	return nil
}

func (j *jiraTrackerService) TestConnection() error {
	return j.jiraClient.Test()
}
