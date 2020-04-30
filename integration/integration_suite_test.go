package integration_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/glestaris/issuez/tracker"
	"github.com/glestaris/issuez/domain"
	"github.com/glestaris/issuez/jira"
	gojira "gopkg.in/andygrunwald/go-jira.v1"
)

type testConfig struct {
	jiraAPIHost     string
	jiraAPIUsername string
	jiraAPIToken    string
	jiraProjectKey  string
	jiraEpicKey     string
}

func newTestConfigMsg(varName string) string {
	return fmt.Sprintf(
		"'%s' environment variable is required to run the jira_client unit test",
		varName,
	)
}

func newTestConfig(t *testing.T) testConfig {
	tc := testConfig{
		jiraAPIHost:     os.Getenv("TEST_JIRA_API_HOST"),
		jiraAPIUsername: os.Getenv("TEST_JIRA_API_USERNAME"),
		jiraAPIToken:    os.Getenv("TEST_JIRA_API_TOKEN"),
		jiraProjectKey:  os.Getenv("TEST_JIRA_PROJECT_KEY"),
		jiraEpicKey:     os.Getenv("TEST_JIRA_EPIC_KEY"),
	}
	require.NotEmpty(
		t, tc.jiraAPIHost, newTestConfigMsg("TEST_JIRA_API_HOST"),
	)
	require.NotEmpty(
		t, tc.jiraAPIUsername, newTestConfigMsg("TEST_JIRA_API_USERNAME"),
	)
	require.NotEmpty(
		t, tc.jiraAPIToken, newTestConfigMsg("TEST_JIRA_API_TOKEN"),
	)
	require.NotEmpty(
		t, tc.jiraProjectKey, newTestConfigMsg("TEST_JIRA_PROJECT_KEY"),
	)
	require.NotEmpty(
		t, tc.jiraEpicKey, newTestConfigMsg("TEST_JIRA_EPIC_KEY"),
	)

	return tc
}

func newJiraClient(tc testConfig) *jira.Client {
	client := jira.NewJiraClient(
		tc.jiraAPIHost, tc.jiraAPIUsername, tc.jiraAPIToken, nil,
	)
	return client
}

func newGoJiraClient(t *testing.T, tc testConfig) *gojira.Client {
	gjc, err := gojira.NewClient(nil, tc.jiraAPIHost)
	require.NoError(t, err)
	gjc.Authentication = &gojira.AuthenticationService{}
	gjc.Authentication.SetBasicAuth(tc.jiraAPIUsername, tc.jiraAPIToken)
	return gjc
}

func newTrackerService(t *testing.T, tc testConfig) tracker.TrackerService {
	trackerService, err := tracker.NewTrackerService(domain.Tracker{
		Type: "jira",
		Config: map[string]string{
			"apiHost":     tc.jiraAPIHost,
			"apiUsername": tc.jiraAPIUsername,
			"apiToken":    tc.jiraAPIToken,
			"projectKey":  tc.jiraProjectKey,
		},
	})
	require.NoError(t, err)

	return trackerService
}
