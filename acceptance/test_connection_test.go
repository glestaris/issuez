package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionTest(t *testing.T) {
	tc := newTestConfig(t)

	cmd := newTestableCommand(tc)
	out, err := cmd.RunSubcommand("test-connection")
	assert.NoError(t, err)
	assert.Equal(t, "OK\n", string(out))
}

func TestConnectionTestError(t *testing.T) {
	tc := newTestConfig(t)

	cmd := newTestableCommand(testConfig{
		issuezExePath:   tc.issuezExePath,
		jiraAPIHost:     tc.jiraAPIHost,
		jiraAPIUsername: "foo@bar.com",
		jiraAPIToken:    "abc123",
	})
	out, err := cmd.RunSubcommand("test-connection")
	assert.Error(t, err)
	assert.Regexp(
		t, regexp.MustCompile(".*Failed to connect to JIRA.*"), string(out),
	)
}
