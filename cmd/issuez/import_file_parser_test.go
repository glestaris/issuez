package main_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	main "github.com/glestaris/issuez/cmd/issuez"
	"github.com/glestaris/issuez/domain"
)

/******************************************************************************
 * General Parsing
 *****************************************************************************/

func TestMarkdownParserEmpty(t *testing.T) {
	markdown := ""
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 0)
}

func TestMarkdownParserSingleBug(t *testing.T) {
	markdown := `[Bug] Bug title

Test para.

Test list:

- A
- B

Epic: 123
Labels: label-1, label-2
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
	require.Equal(t, "Bug", issues[0].Type.String())
	require.Equal(t, "Bug title", issues[0].Title)
	require.Equal(t, "123", issues[0].Epic.ID)
	require.Equal(
		t, []domain.Label{
			{Label: "label-1"},
			{Label: "label-2"},
		}, issues[0].Labels,
	)
}

func TestMarkdownParserMultipleIssues(t *testing.T) {
	markdown := `[Bug] Bug title

Test para.

Test list:

- A
- B

Epic: 123
Labels: label-1, label-2

---

A story

Hello world.

Epic: 99
Labels: label-3
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 2)
	require.Equal(t, "Bug", issues[0].Type.String())
	require.Equal(t, "Bug title", issues[0].Title)
	require.Equal(t, "123", issues[0].Epic.ID)
	require.Equal(
		t, []domain.Label{
			{Label: "label-1"},
			{Label: "label-2"},
		}, issues[0].Labels,
	)

	require.Equal(t, "User Story", issues[1].Type.String())
	require.Equal(t, "A story", issues[1].Title)
	require.Equal(t, "99", issues[1].Epic.ID)
	require.Equal(
		t, []domain.Label{
			{Label: "label-3"},
		}, issues[1].Labels,
	)
}

/******************************************************************************
 * Parsing story type
 *****************************************************************************/

func TestMarkdownParserChore(t *testing.T) {
	// Chore
	markdown := "[Chore] Chore title"
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(t, "Chore", issues[0].Type.String())
	require.Equal(t, "Chore title", issues[0].Title)

	// Task
	markdown = "[Task] Task title"
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(t, "Chore", issues[0].Type.String())
	require.Equal(t, "Task title", issues[0].Title)
}

func TestMarkdownParserUserStory(t *testing.T) {
	// No tag
	markdown := "Title"
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(t, "User Story", issues[0].Type.String())
	require.Equal(t, "Title", issues[0].Title)

	// Story
	markdown = "[Story] Story title"
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(t, "User Story", issues[0].Type.String())
	require.Equal(t, "Story title", issues[0].Title)

	// Issue
	markdown = "[Issue] Issue title"
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(t, "User Story", issues[0].Type.String())
	require.Equal(t, "Issue title", issues[0].Title)
}

/******************************************************************************
 * Epic
 *****************************************************************************/

func TestMarkdownParserEpic(t *testing.T) {
	// Using E
	markdown := `Title

E: hello-world`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(t, "hello-world", issues[0].Epic.ID)

	// Using Epic
	markdown = `Title

Epic: hello-world`
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(t, "hello-world", issues[0].Epic.ID)

	// Epic ID w/ numbers
	markdown = `Title

Epic: 1234`
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(t, "1234", issues[0].Epic.ID)

	// Epic ID with spaces
	markdown = `Title

Epic:   hello-world      `
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(t, "hello-world", issues[0].Epic.ID)

	// No epic
	markdown = "Title"
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Nil(t, issues[0].Epic)
}

/******************************************************************************
 * Labels
 *****************************************************************************/

func TestMarkdownParserLabels(t *testing.T) {
	// Using L
	markdown := `Title

L: label-1, label-2`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(
		t, []domain.Label{
			{Label: "label-1"},
			{Label: "label-2"},
		}, issues[0].Labels,
	)

	// Using Labels
	markdown = `Title

Labels: label-1, label-2`
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(
		t, []domain.Label{
			{Label: "label-1"},
			{Label: "label-2"},
		}, issues[0].Labels,
	)

	// Labels without spaces
	markdown = `Title

Labels: label-1,label-2`
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(
		t, []domain.Label{
			{Label: "label-1"},
			{Label: "label-2"},
		}, issues[0].Labels,
	)

	// Labels with many spaces
	markdown = `Title

Labels:      label-1,     label-2    `
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Equal(
		t, []domain.Label{
			{Label: "label-1"},
			{Label: "label-2"},
		}, issues[0].Labels,
	)

	// No labels
	markdown = "Title"
	issues, err = main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)
	require.Len(t, issues, 1)
	require.Empty(t, issues[0].Labels)
}

/******************************************************************************
 * Horizontal rules in the beginning and/or the end of the file
 *****************************************************************************/

func TestMarkdownParserSingleBugStartingWHR(t *testing.T) {
	markdown := `---

[Bug] Bug title

Test para.

Test list:

- A
- B

Epic: 123
Labels: label-1, label-2
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
}

func TestMarkdownParserSingleBugEndingWHR(t *testing.T) {
	markdown := `[Bug] Bug title

Test para.

Test list:

- A
- B

Epic: 123
Labels: label-1, label-2

---
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
}

func TestMarkdownParserSingleBugSurroundedByHR(t *testing.T) {
	markdown := `---

[Bug] Bug title

Test para.

Test list:

- A
- B

Epic: 123
Labels: label-1, label-2

---
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
}

func TestMarkdownParserMultipleIssuesWithMultipleHR(t *testing.T) {
	markdown := `[Bug] Bug title

Test para.

Test list:

- A
- B

Epic: 123
Labels: label-1, label-2

---

---
---

A story

Hello world.

Epic: 99
Labels: label-3
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 2)
}

/******************************************************************************
 * Parse description
 *****************************************************************************/

func TestMarkdownParserDescription(t *testing.T) {
	markdown := `[Bug] Bug title

Test para.

Test list:

- A
- B

Epic: 123
Labels: label-1, label-2
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
	require.Equal(
		t, &domain.Document{
			Nodes: []domain.DocumentNode{
				{
					Type: domain.DocumentNodeTypeParagraph,
					ParagraphData: &domain.ParagraphData{
						Content: domain.TextContainer{
							Elements: []domain.TextElement{
								{Text: "Test para."},
							},
						},
					},
				},
				{
					Type: domain.DocumentNodeTypeParagraph,
					ParagraphData: &domain.ParagraphData{
						Content: domain.TextContainer{
							Elements: []domain.TextElement{
								{Text: "Test list:"},
							},
						},
					},
				},
				{
					Type: domain.DocumentNodeTypeList,
					ListData: &domain.ListData{
						IsOrdered: false,
						Items: []domain.TextContainer{
							{
								Elements: []domain.TextElement{
									{Text: "A"},
								},
							},
							{
								Elements: []domain.TextElement{
									{Text: "B"},
								},
							},
						},
					},
				},
			},
		}, issues[0].Description,
	)
}

func TestMarkdownParserDescriptionWithMarks(t *testing.T) {
	markdown := `[Bug] Bug title

Test para **bold** and _italics_ and ~~strikethrough~~.
Sometimes, ~~**combinations**~~ of both. And some [links](https://google.com).
And **[links with marks](https://google.com)**.

- **Bold** list item

Epic: 123
Labels: label-1, label-2
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
	require.Equal(
		t, &domain.Document{
			Nodes: []domain.DocumentNode{
				{
					Type: domain.DocumentNodeTypeParagraph,
					ParagraphData: &domain.ParagraphData{
						Content: domain.TextContainer{
							Elements: []domain.TextElement{
								{Text: "Test para "},
								{
									Text: "bold",
									Mode: domain.TextMode{Bold: true},
								},
								{Text: " and "},
								{
									Text: "italics",
									Mode: domain.TextMode{Italics: true},
								},
								{Text: " and "},
								{
									Text: "strikethrough",
									Mode: domain.TextMode{Strikethrough: true},
								},
								{Text: ".\nSometimes, "},
								{
									Text: "combinations",
									Mode: domain.TextMode{
										Strikethrough: true,
										Bold:          true,
									},
								},
								{Text: " of both. And some "},
								{
									Text:    "links",
									LinkURL: "https://google.com",
								},
								{Text: ".\nAnd "},
								{
									Text:    "links with marks",
									Mode:    domain.TextMode{Bold: true},
									LinkURL: "https://google.com",
								},
								{Text: "."},
							},
						},
					},
				},
				{
					Type: domain.DocumentNodeTypeList,
					ListData: &domain.ListData{
						IsOrdered: false,
						Items: []domain.TextContainer{
							{
								Elements: []domain.TextElement{
									{
										Text: "Bold",
										Mode: domain.TextMode{Bold: true},
									},
									{Text: " list item"},
								},
							},
						},
					},
				},
			},
		}, issues[0].Description,
	)
}

func TestMarkdownParserDescriptionWithCode(t *testing.T) {
	markdown := "[Bug] Bug title\n\nTest para `with code`.\n\nEpic: 123"
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
	require.Equal(
		t, &domain.Document{
			Nodes: []domain.DocumentNode{
				{
					Type: domain.DocumentNodeTypeParagraph,
					ParagraphData: &domain.ParagraphData{
						Content: domain.TextContainer{
							Elements: []domain.TextElement{
								{Text: "Test para "},
								{
									Text: "with code",
									Mode: domain.TextMode{Code: true},
								},
								{Text: "."},
							},
						},
					},
				},
			},
		}, issues[0].Description,
	)
}

func TestMarkdownParserDescriptionWithCodeBlock(t *testing.T) {
	markdown := "[Bug] Bug title\n\nTest para with code block:\n\n" +
		"```python\nx = 12\n```\n\nEpic: 123"
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
	require.Equal(
		t, &domain.Document{
			Nodes: []domain.DocumentNode{
				{
					Type: domain.DocumentNodeTypeParagraph,
					ParagraphData: &domain.ParagraphData{
						Content: domain.TextContainer{
							Elements: []domain.TextElement{
								{Text: "Test para with code block:"},
							},
						},
					},
				},
				{
					Type: domain.DocumentNodeTypeCodeBlock,
					CodeBlockData: &domain.CodeBlockData{
						Code:     "x = 12\n",
						Language: "python",
					},
				},
			},
		}, issues[0].Description,
	)
}

func TestMarkdownParserDescriptionHeadings(t *testing.T) {
	markdown := `[Bug] Bug title

# Title

Hello world.

## Subtitle

Test test.

Epic: 123
Labels: label-1, label-2
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
	require.Equal(
		t, &domain.Document{
			Nodes: []domain.DocumentNode{
				{
					Type: domain.DocumentNodeTypeHeading,
					HeadingData: &domain.HeadingData{
						Level: domain.HeadingLevel1,
						Text:  "Title",
					},
				},
				{
					Type: domain.DocumentNodeTypeParagraph,
					ParagraphData: &domain.ParagraphData{
						Content: domain.TextContainer{
							Elements: []domain.TextElement{
								{Text: "Hello world."},
							},
						},
					},
				},
				{
					Type: domain.DocumentNodeTypeHeading,
					HeadingData: &domain.HeadingData{
						Level: domain.HeadingLevel2,
						Text:  "Subtitle",
					},
				},
				{
					Type: domain.DocumentNodeTypeParagraph,
					ParagraphData: &domain.ParagraphData{
						Content: domain.TextContainer{
							Elements: []domain.TextElement{
								{Text: "Test test."},
							},
						},
					},
				},
			},
		}, issues[0].Description,
	)
}

func TestMarkdownParserDescriptionOrderedLists(t *testing.T) {
	markdown := `[Bug] Bug title

An ordered list:

1. List item 1 
1. List item 2 
1. List item 3

Epic: 123
Labels: label-1, label-2
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
	require.Equal(
		t, &domain.Document{
			Nodes: []domain.DocumentNode{
				{
					Type: domain.DocumentNodeTypeParagraph,
					ParagraphData: &domain.ParagraphData{
						Content: domain.TextContainer{
							Elements: []domain.TextElement{
								{Text: "An ordered list:"},
							},
						},
					},
				},
				{
					Type: domain.DocumentNodeTypeList,
					ListData: &domain.ListData{
						IsOrdered: true,
						Items: []domain.TextContainer{
							{
								Elements: []domain.TextElement{
									{Text: "List item 1"},
								},
							},
							{
								Elements: []domain.TextElement{
									{Text: "List item 2"},
								},
							},
							{
								Elements: []domain.TextElement{
									{Text: "List item 3"},
								},
							},
						},
					},
				},
			},
		}, issues[0].Description,
	)
}

func TestMarkdownParserDescriptionWithoutFooter(t *testing.T) {
	markdown := `[Bug] Bug title

Hello world.
`
	issues, err := main.ParseImportFile(
		strings.NewReader(markdown),
	)
	require.NoError(t, err)

	require.Len(t, issues, 1)
	require.Equal(
		t, &domain.Document{
			Nodes: []domain.DocumentNode{
				{
					Type: domain.DocumentNodeTypeParagraph,
					ParagraphData: &domain.ParagraphData{
						Content: domain.TextContainer{
							Elements: []domain.TextElement{
								{Text: "Hello world."},
							},
						},
					},
				},
			},
		}, issues[0].Description,
	)
}
