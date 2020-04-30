package acceptance_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	gojira "gopkg.in/andygrunwald/go-jira.v1"
)

type testConfig struct {
	issuezExePath   string
	jiraAPIHost     string
	jiraAPIUsername string
	jiraAPIToken    string
	jiraProjectKey  string
}

func newTestConfigMsg(varName string) string {
	return fmt.Sprintf(
		"'%s' environment variable is required to run the jira_client unit test",
		varName,
	)
}

func newTestConfig(t *testing.T) testConfig {
	tc := testConfig{
		issuezExePath:   os.Getenv("TEST_ISSUEZ_EXE_PATH"),
		jiraAPIHost:     os.Getenv("TEST_JIRA_API_HOST"),
		jiraAPIUsername: os.Getenv("TEST_JIRA_API_USERNAME"),
		jiraAPIToken:    os.Getenv("TEST_JIRA_API_TOKEN"),
		jiraProjectKey:  os.Getenv("TEST_JIRA_PROJECT_KEY"),
	}
	require.NotEmpty(
		t, tc.issuezExePath, newTestConfigMsg("TEST_ISSUEZ_EXE_PATH"),
	)
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

	return tc
}

type testableCommand struct {
	config testConfig
}

func newTestableCommand(tc testConfig) *testableCommand {
	return &testableCommand{config: tc}
}

func (c *testableCommand) RunSubcommand(
	subcommand string, arg ...string,
) ([]byte, error) {
	cmdArg := []string{
		"--api", c.config.jiraAPIHost,
		"--username", c.config.jiraAPIUsername,
		"--token", c.config.jiraAPIToken,
		subcommand,
	}
	cmdArg = append(cmdArg, arg...)
	cmd := exec.Command(c.config.issuezExePath, cmdArg...)
	return cmd.CombinedOutput()
}

func newGoJiraClient(t *testing.T, tc testConfig) *gojira.Client {
	gjc, err := gojira.NewClient(nil, tc.jiraAPIHost)
	require.NoError(t, err)
	gjc.Authentication = &gojira.AuthenticationService{}
	gjc.Authentication.SetBasicAuth(tc.jiraAPIUsername, tc.jiraAPIToken)
	return gjc
}
