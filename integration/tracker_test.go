package integration_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/glestaris/issuez/domain"
)

func TestTrackerCreateSingleIssue(t *testing.T) {
	tc := newTestConfig(t)
	trackerService := newTrackerService(t, tc)
	gjc := newGoJiraClient(t, tc)

	issues := []*domain.Issue{
		{
			Type:  domain.IssueTypeBug,
			Title: "Hello world 1",
			Labels: []domain.Label{
				{Label: "label-1"},
			},
		},
	}
	err := trackerService.ImportIssues(issues)
	require.NoError(t, err)

	issue, _, err := gjc.Issue.Get(issues[0].ID, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(issues[0].ID)
	require.Equal(t, "Bug", issue.Fields.Type.Name)
	require.Equal(t, "Hello world 1", issue.Fields.Summary)
	require.Equal(t, []string{"label-1"}, issue.Fields.Labels)
}

func TestTrackerCreateManyIssues(t *testing.T) {
	tc := newTestConfig(t)
	trackerService := newTrackerService(t, tc)
	gjc := newGoJiraClient(t, tc)

	issues := []*domain.Issue{
		{
			Type:  domain.IssueTypeBug,
			Title: "Hello world 1",
			Labels: []domain.Label{
				{Label: "label-1"},
			},
		},
		{
			Type:  domain.IssueTypeStory,
			Title: "Hello world 2",
			Labels: []domain.Label{
				{Label: "label-2"},
				{Label: "label-3"},
			},
			Epic: &domain.Epic{ID: tc.jiraEpicKey},
		},
	}
	err := trackerService.ImportIssues(issues)
	require.NoError(t, err)

	issue, _, err := gjc.Issue.Get(issues[0].ID, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(issues[0].ID)
	require.Equal(t, "Bug", issue.Fields.Type.Name)
	require.Equal(t, "Hello world 1", issue.Fields.Summary)
	require.Equal(t, []string{"label-1"}, issue.Fields.Labels)
	issue, _, err = gjc.Issue.Get(issues[1].ID, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(issues[1].ID)
	require.Equal(t, "Story", issue.Fields.Type.Name)
	require.Equal(t, "Hello world 2", issue.Fields.Summary)
	require.Equal(t, []string{"label-2", "label-3"}, issue.Fields.Labels)
	require.Equal(t, tc.jiraEpicKey, issue.Fields.Parent.Key)
}

func TestTrackerMapsDescriptionDocument(t *testing.T) {
	tc := newTestConfig(t)
	trackerService := newTrackerService(t, tc)
	gjc := newGoJiraClient(t, tc)

	issues := []*domain.Issue{
		{
			Type:  domain.IssueTypeBug,
			Title: "Hello world 1",
			Description: &domain.Document{
				Nodes: []domain.DocumentNode{
					{
						Type: domain.DocumentNodeTypeHeading,
						HeadingData: &domain.HeadingData{
							Level: domain.HeadingLevel3,
							Text:  "Paragraph coming up",
						},
					},
					{
						Type: domain.DocumentNodeTypeParagraph,
						ParagraphData: &domain.ParagraphData{
							Content: domain.TextContainer{
								Elements: []domain.TextElement{
									{Text: "Test paragraph. "},
									{
										Text: "Bold sentence.",
										Mode: domain.TextMode{Bold: true},
									},
								},
							},
						},
					},
				},
			},
			Labels: []domain.Label{
				{Label: "label-1"},
			},
		},
	}
	err := trackerService.ImportIssues(issues)
	require.NoError(t, err)

	issue, _, err := gjc.Issue.Get(issues[0].ID, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(issues[0].ID)
	require.Equal(t, `h3. Paragraph coming up

Test paragraph. *Bold sentence.*`,
		issue.Fields.Description)
}
