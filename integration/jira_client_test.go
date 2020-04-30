package integration_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/glestaris/issuez/jira"
)

func TestConnectionTest(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)

	err := jiraClient.Test()
	require.NoError(t, err)
}

func TestConnectionTestInvalidAuth(t *testing.T) {
	tc := newTestConfig(t)
	tc.jiraAPIToken = "fail"
	jiraClient := newJiraClient(tc)

	err := jiraClient.Test()
	require.Error(t, err)
}

func TestImportZeroIssues(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)

	resp, err := jiraClient.ImportIssues([]*jira.Issue{})
	require.NoError(t, err)
	require.Len(t, resp, 0)
}

func TestImportOneIssue(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)
	gjc := newGoJiraClient(t, tc)

	descriptionDoc := jira.NewADFDocument()
	descriptionDoc.AddParagraph().
		AddText("This is a paragraph", jira.ADFTextMode{})

	resp, err := jiraClient.ImportIssues([]*jira.Issue{
		{
			Type:        jira.IssueTypeStory,
			Summary:     "Hello world",
			Description: descriptionDoc,
			EpicKey:     tc.jiraEpicKey,
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-1", "label-2"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)

	issue, _, err := gjc.Issue.Get(resp[0].NewIssueKey, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(resp[0].NewIssueKey)
	require.Equal(t, issue.Fields.Project.Key, tc.jiraProjectKey)
	require.Equal(t, issue.Fields.Type.Name, "Story")
	require.Equal(t, issue.Fields.Summary, "Hello world")
	require.Equal(t, issue.Fields.Parent.Key, tc.jiraEpicKey)
	require.Equal(t, issue.Fields.Labels, []string{"label-1", "label-2"})
}

func TestImportTask(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)
	gjc := newGoJiraClient(t, tc)

	descriptionDoc := jira.NewADFDocument()
	descriptionDoc.AddParagraph().
		AddText("This is a task paragraph", jira.ADFTextMode{})

	resp, err := jiraClient.ImportIssues([]*jira.Issue{
		{
			Type:        jira.IssueTypeTask,
			Summary:     "Hello world of tasks",
			Description: descriptionDoc,
			EpicKey:     tc.jiraEpicKey,
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-2"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)

	issue, _, err := gjc.Issue.Get(resp[0].NewIssueKey, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(resp[0].NewIssueKey)
	require.Equal(t, issue.Fields.Project.Key, tc.jiraProjectKey)
	require.Equal(t, issue.Fields.Type.Name, "Task")
	require.Equal(t, issue.Fields.Summary, "Hello world of tasks")
	require.Equal(t, issue.Fields.Parent.Key, tc.jiraEpicKey)
	require.Equal(t, issue.Fields.Labels, []string{"label-2"})
}

func TestImportBug(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)
	gjc := newGoJiraClient(t, tc)

	descriptionDoc := jira.NewADFDocument()
	descriptionDoc.AddParagraph().
		AddText("This is a bug paragraph", jira.ADFTextMode{})

	resp, err := jiraClient.ImportIssues([]*jira.Issue{
		{
			Type:        jira.IssueTypeBug,
			Summary:     "Hello world of bugs",
			Description: descriptionDoc,
			EpicKey:     tc.jiraEpicKey,
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-3"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)

	issue, _, err := gjc.Issue.Get(resp[0].NewIssueKey, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(resp[0].NewIssueKey)
	require.Equal(t, issue.Fields.Project.Key, tc.jiraProjectKey)
	require.Equal(t, issue.Fields.Type.Name, "Bug")
	require.Equal(t, issue.Fields.Summary, "Hello world of bugs")
	require.Equal(t, issue.Fields.Parent.Key, tc.jiraEpicKey)
	require.Equal(t, issue.Fields.Labels, []string{"label-3"})
}

func TestImportManyIssues(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)
	gjc := newGoJiraClient(t, tc)

	descriptionADoc := jira.NewADFDocument()
	descriptionADoc.AddParagraph().
		AddText("This is a paragraph", jira.ADFTextMode{})

	descriptionBDoc := jira.NewADFDocument()
	descriptionBDoc.AddParagraph().
		AddText("This is another paragraph", jira.ADFTextMode{})

	resp, err := jiraClient.ImportIssues([]*jira.Issue{
		{
			Type:        jira.IssueTypeStory,
			Summary:     "Hello world A",
			Description: descriptionADoc,
			EpicKey:     tc.jiraEpicKey,
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-1"},
		},
		{
			Type:        jira.IssueTypeBug,
			Summary:     "Hello world B",
			Description: descriptionBDoc,
			EpicKey:     tc.jiraEpicKey,
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-2"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 2)

	issue, _, err := gjc.Issue.Get(resp[0].NewIssueKey, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(resp[0].NewIssueKey)
	require.Equal(t, issue.Fields.Type.Name, "Story")
	require.Equal(t, issue.Fields.Summary, "Hello world A")
	require.Equal(t, issue.Fields.Labels, []string{"label-1"})
	issue, _, err = gjc.Issue.Get(resp[1].NewIssueKey, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(resp[1].NewIssueKey)
	require.Equal(t, issue.Fields.Type.Name, "Bug")
	require.Equal(t, issue.Fields.Summary, "Hello world B")
	require.Equal(t, issue.Fields.Labels, []string{"label-2"})
}

func TestImportNoEpic(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)
	gjc := newGoJiraClient(t, tc)

	descriptionDoc := jira.NewADFDocument()
	descriptionDoc.AddParagraph().
		AddText("This is a paragraph", jira.ADFTextMode{})

	resp, err := jiraClient.ImportIssues([]*jira.Issue{
		{
			Type:        jira.IssueTypeStory,
			Summary:     "Hello world",
			Description: descriptionDoc,
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-1", "label-2"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)

	issue, _, err := gjc.Issue.Get(resp[0].NewIssueKey, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(resp[0].NewIssueKey)
	require.Nil(t, issue.Fields.Parent)
}

func TestImportNoLabels(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)
	gjc := newGoJiraClient(t, tc)

	descriptionDoc := jira.NewADFDocument()
	descriptionDoc.AddParagraph().
		AddText("This is a paragraph", jira.ADFTextMode{})

	resp, err := jiraClient.ImportIssues([]*jira.Issue{
		{
			Type:        jira.IssueTypeStory,
			Summary:     "Hello world",
			Description: descriptionDoc,
			EpicKey:     tc.jiraEpicKey,
			ProjectKey:  tc.jiraProjectKey,
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 1)

	issue, _, err := gjc.Issue.Get(resp[0].NewIssueKey, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(resp[0].NewIssueKey)
	require.Empty(t, issue.Fields.Labels)
}

func TestImportManyIssuesWAuthError(t *testing.T) {
	tc := newTestConfig(t)
	tc.jiraAPIToken = "fail"
	jiraClient := newJiraClient(tc)

	descriptionADoc := jira.NewADFDocument()
	descriptionADoc.AddParagraph().
		AddText("This is a paragraph", jira.ADFTextMode{})

	descriptionBDoc := jira.NewADFDocument()
	descriptionBDoc.AddParagraph().
		AddText("This is another paragraph", jira.ADFTextMode{})

	resp, err := jiraClient.ImportIssues([]*jira.Issue{
		{
			Type:        jira.IssueTypeStory,
			Summary:     "Hello world A",
			Description: descriptionADoc,
			EpicKey:     tc.jiraEpicKey,
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-1"},
		},
		{
			Type:        jira.IssueTypeBug,
			Summary:     "Hello world B",
			Description: descriptionBDoc,
			EpicKey:     tc.jiraEpicKey,
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-2"},
		},
	})
	require.Error(t, err)
	require.EqualError(t, err, "Failed to create issues: 403 Forbidden")
	require.Len(t, resp, 0)
}

func TestImportManyIssuesWError(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)
	gjc := newGoJiraClient(t, tc)

	descriptionADoc := jira.NewADFDocument()
	descriptionADoc.AddParagraph().
		AddText("This is a paragraph", jira.ADFTextMode{})

	descriptionBDoc := jira.NewADFDocument()
	descriptionBDoc.AddParagraph().
		AddText("This is another paragraph", jira.ADFTextMode{})

	resp, err := jiraClient.ImportIssues([]*jira.Issue{
		{
			Type:        jira.IssueTypeStory,
			Summary:     "Hello world A",
			Description: descriptionADoc,
			EpicKey:     "NON-EXISTENT-EPIC",
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-1"},
		},
		{
			Type:        jira.IssueTypeBug,
			Summary:     "Hello world B",
			Description: descriptionBDoc,
			EpicKey:     tc.jiraEpicKey,
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-2"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 2)

	require.Error(t, resp[0].Err)
	require.NoError(t, resp[1].Err)
	issue, _, err := gjc.Issue.Get(resp[1].NewIssueKey, nil)
	require.NoError(t, err)
	defer gjc.Issue.Delete(resp[1].NewIssueKey)
	require.Equal(t, issue.Fields.Type.Name, "Bug")
	require.Equal(t, issue.Fields.Summary, "Hello world B")
	require.Equal(t, issue.Fields.Labels, []string{"label-2"})
}

func TestImportManyIssuesWErrorAll(t *testing.T) {
	tc := newTestConfig(t)
	jiraClient := newJiraClient(tc)

	descriptionADoc := jira.NewADFDocument()
	descriptionADoc.AddParagraph().
		AddText("This is a paragraph", jira.ADFTextMode{})

	descriptionBDoc := jira.NewADFDocument()
	descriptionBDoc.AddParagraph().
		AddText("This is another paragraph", jira.ADFTextMode{})

	resp, err := jiraClient.ImportIssues([]*jira.Issue{
		{
			Type:        jira.IssueTypeStory,
			Summary:     "Hello world A",
			Description: descriptionADoc,
			EpicKey:     "NON-EXISTENT-EPIC",
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-1"},
		},
		{
			Type:        jira.IssueTypeBug,
			Summary:     "Hello world B",
			Description: descriptionBDoc,
			EpicKey:     "NON-EXISTENT-EPIC",
			ProjectKey:  tc.jiraProjectKey,
			Labels:      []string{"label-2"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp, 2)

	require.Error(t, resp[0].Err)
	require.Error(t, resp[1].Err)
}
