package acceptance_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	tc := newTestConfig(t)
	gjc := newGoJiraClient(t, tc)

	// define cases
	cases := []struct {
		fixturePath string
		returnErr   bool
		name        string
		callback    func(*testing.T, []byte)
	}{
		{
			name:        "Task",
			fixturePath: "./fixtures/task.md",
			returnErr:   false,
			callback: func(t *testing.T, out []byte) {
				strOut := string(out)
				require.Regexp(
					t, regexp.MustCompile(
						"Found 1 issues in the markdown file",
					),
					strOut,
				)
				matches := regexp.MustCompile(
					`- Chore \(([A-Z0-9\-]+)\)`,
				).FindAllStringSubmatch(strOut, 1000)
				require.Len(t, matches, 1)
				require.NotEqual(t, matches[0][1], "FAILED")

				issue, _, err := gjc.Issue.Get(matches[0][1], nil)
				require.NoError(t, err)
				defer gjc.Issue.Delete(matches[0][1])
				require.Equal(t, tc.jiraProjectKey, issue.Fields.Project.Key)
				require.Equal(t, "Task", issue.Fields.Type.Name)
				require.Equal(t, "A summary", issue.Fields.Summary)
				require.Equal(t, "TEST-1", issue.Fields.Parent.Key)
				require.Equal(t, []string{"label-1", "label-2"},
					issue.Fields.Labels)
			},
		},
		{
			name:        "Multiple Issues",
			fixturePath: "./fixtures/multiple.md",
			returnErr:   false,
			callback: func(t *testing.T, out []byte) {
				strOut := string(out)
				require.Regexp(
					t, regexp.MustCompile(
						"Found 2 issues in the markdown file",
					),
					strOut,
				)

				usMatches := regexp.MustCompile(
					`- User Story \(([A-Z0-9\-]+)\)`,
				).FindAllStringSubmatch(strOut, 1000)
				require.Len(t, usMatches, 1)
				usID := usMatches[0][1]
				require.NotEqual(t, usID, "FAILED")

				issue, _, err := gjc.Issue.Get(usID, nil)
				require.NoError(t, err)
				defer gjc.Issue.Delete(usID)
				require.Equal(t, "A summary", issue.Fields.Summary)
				require.Equal(t, `Code example:

{code:python}test = 12
{code}`,
					issue.Fields.Description)

				chMatches := regexp.MustCompile(
					`- Chore \(([A-Z0-9\-]+)\)`,
				).FindAllStringSubmatch(strOut, 1000)
				require.Len(t, chMatches, 1)
				chID := chMatches[0][1]
				require.NotEqual(t, chID, "FAILED")

				issue, _, err = gjc.Issue.Get(chID, nil)
				require.NoError(t, err)
				defer gjc.Issue.Delete(chID)
				require.Equal(t, "A chore story", issue.Fields.Summary)
				require.Equal(t, "Do it.", issue.Fields.Description)
			},
		},
	}

	// run through cases
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			cmd := newTestableCommand(tc)
			out, err := cmd.RunSubcommand("import", "-p", tc.jiraProjectKey,
				test.fixturePath)
			if test.returnErr {
				require.Error(t, err)
				return
			}
			if test.callback == nil {
				return
			}
			test.callback(t, out)
		})
	}
}
