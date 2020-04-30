package jira_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/glestaris/issuez/jira"
)

func TestADFDocumentHeading(t *testing.T) {
	doc := jira.NewADFDocument()
	require.NotNil(t, doc)

	doc.AddHeading(jira.ADFHeadingLevel1, "A paragraph")

	docJSON, err := json.Marshal(doc)
	require.NoError(t, err)
	require.JSONEq(t, `
{
  "version": 1,
  "type": "doc",
  "content": [
    {
      "type": "heading",
      "attrs": {
        "level": 1
      },
      "content": [
        {
          "type": "text",
          "text": "A paragraph"
        }
      ]
    }
  ]
}
    `, string(docJSON))
}

func TestADFDocumentParagraph(t *testing.T) {
	doc := jira.NewADFDocument()
	require.NotNil(t, doc)

	p := doc.AddParagraph()
	p.AddText("Hello world of ", jira.ADFTextMode{})
	p.AddText("bold text", jira.ADFTextMode{Strong: true})
	p.AddText(" which can be also ", jira.ADFTextMode{})
	p.AddText("italics", jira.ADFTextMode{Em: true})
	p.AddText(", ", jira.ADFTextMode{})
	p.AddText("stikethrough", jira.ADFTextMode{Strike: true})
	p.AddText(" or ", jira.ADFTextMode{})
	p.AddText("code", jira.ADFTextMode{Code: true})
	p.AddText(".", jira.ADFTextMode{})
	p.AddText(" And of course, it can be ", jira.ADFTextMode{})
	p.AddText("multiple things at once", jira.ADFTextMode{
		Strong: true,
		Em:     true,
	})
	p.AddText(" as well.", jira.ADFTextMode{})

	docJSON, err := json.Marshal(doc)
	require.NoError(t, err)
	require.JSONEq(t, `
{
  "version": 1,
  "type": "doc",
  "content": [
    {
      "type": "paragraph",
      "content": [
        {
          "type": "text",
          "text": "Hello world of "
        },
        {
          "type": "text",
          "text": "bold text",
          "marks": [{ "type": "strong" }]
        },
        {
          "type": "text",
          "text": " which can be also "
        },
        {
          "type": "text",
          "text": "italics",
          "marks": [{ "type": "em" }]
        },
        {
          "type": "text",
          "text": ", "
        },
        {
          "type": "text",
          "text": "stikethrough",
          "marks": [{ "type": "strike" }]
        },
        {
          "type": "text",
          "text": " or "
        },
        {
          "type": "text",
          "text": "code",
          "marks": [{ "type": "code" }]
        },
        {
          "type": "text",
          "text": "."
        },
        {
          "type": "text",
          "text": " And of course, it can be "
        },
        {
          "type": "text",
          "text": "multiple things at once",
          "marks": [{ "type": "strong" }, { "type": "em" }]
        },
        {
          "type": "text",
          "text": " as well."
        }
      ]
    }
  ]
}
    `, string(docJSON))
}

func TestADFDocumentOrderedList(t *testing.T) {
	doc := jira.NewADFDocument()
	require.NotNil(t, doc)

	ol := doc.AddOrderedList()
	ol.AddItem().AddText("New nodes", jira.ADFTextMode{})
	olItem := ol.AddItem()
	olItem.AddText("And ", jira.ADFTextMode{})
	olItem.AddLink("links", "https://google.com", jira.ADFTextMode{})

	docJSON, err := json.Marshal(doc)
	require.NoError(t, err)
	require.JSONEq(t, `
{
  "version": 1,
  "type": "doc",
  "content": [
    {
      "type": "orderedList",
      "content": [
        {
          "type": "listItem",
          "content": [
            {
              "type": "paragraph",
              "content": [
                {
                  "type": "text",
                  "text": "New nodes"
                }
              ]
            }
          ]
        },
        {
          "type": "listItem",
          "content": [
            {
              "type": "paragraph",
              "content": [
                {
                  "type": "text",
                  "text": "And "
                },
                {
                  "type": "text",
                  "text": "links",
                  "marks": [
                    {
                      "type": "link",
                      "attrs": { "href": "https://google.com" }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
    `, string(docJSON))
}

func TestADFDocumentMany(t *testing.T) {
	doc := jira.NewADFDocument()
	require.NotNil(t, doc)

	doc.AddHeading(jira.ADFHeadingLevel1, "A paragraph")

	p := doc.AddParagraph()
	p.AddText("Hello world of ", jira.ADFTextMode{})
	p.AddText("bold text", jira.ADFTextMode{Strong: true})
	p.AddText(" which can be also ", jira.ADFTextMode{})
	p.AddText("italics", jira.ADFTextMode{Em: true})
	p.AddText(", ", jira.ADFTextMode{})
	p.AddText("stikethrough", jira.ADFTextMode{Strike: true})
	p.AddText(" or ", jira.ADFTextMode{})
	p.AddText("code", jira.ADFTextMode{Code: true})
	p.AddText(".", jira.ADFTextMode{})
	p.AddText(" And of course, it can be ", jira.ADFTextMode{})
	p.AddText("multiple things at once", jira.ADFTextMode{
		Strong: true,
		Em:     true,
	})
	p.AddText(" as well.", jira.ADFTextMode{})

	doc.AddHeading(jira.ADFHeadingLevel2, "An ordered list...")

	ol := doc.AddOrderedList()
	ol.AddItem().AddText("New nodes", jira.ADFTextMode{})
	olItem := ol.AddItem()
	olItem.AddText("And ", jira.ADFTextMode{})
	olItem.AddLink("links", "https://google.com", jira.ADFTextMode{})

	doc.AddHeading(jira.ADFHeadingLevel2, "An un-ordered (bullet) list...")

	uol := doc.AddBulletList()
	uol.AddItem().AddText("We like bullets", jira.ADFTextMode{})
	uol.AddItem().AddText("Too", jira.ADFTextMode{})

	doc.AddHeading(jira.ADFHeadingLevel3, "And some code...")

	doc.AddCodeBlock("python", "x = 12\n")

	docJSON, err := json.Marshal(doc)
	require.NoError(t, err)
	require.JSONEq(t, `
{
  "version": 1,
  "type": "doc",
  "content": [
    {
      "type": "heading",
      "attrs": {
        "level": 1
      },
      "content": [
        {
          "type": "text",
          "text": "A paragraph"
        }
      ]
    },
    {
      "type": "paragraph",
      "content": [
        {
          "type": "text",
          "text": "Hello world of "
        },
        {
          "type": "text",
          "text": "bold text",
          "marks": [{ "type": "strong" }]
        },
        {
          "type": "text",
          "text": " which can be also "
        },
        {
          "type": "text",
          "text": "italics",
          "marks": [{ "type": "em" }]
        },
        {
          "type": "text",
          "text": ", "
        },
        {
          "type": "text",
          "text": "stikethrough",
          "marks": [{ "type": "strike" }]
        },
        {
          "type": "text",
          "text": " or "
        },
        {
          "type": "text",
          "text": "code",
          "marks": [{ "type": "code" }]
        },
        {
          "type": "text",
          "text": "."
        },
        {
          "type": "text",
          "text": " And of course, it can be "
        },
        {
          "type": "text",
          "text": "multiple things at once",
          "marks": [{ "type": "strong" }, { "type": "em" }]
        },
        {
          "type": "text",
          "text": " as well."
        }
      ]
    },
    {
      "type": "heading",
      "attrs": {
        "level": 2
      },
      "content": [
        {
          "type": "text",
          "text": "An ordered list..."
        }
      ]
    },
    {
      "type": "orderedList",
      "content": [
        {
          "type": "listItem",
          "content": [
            {
              "type": "paragraph",
              "content": [
                {
                  "type": "text",
                  "text": "New nodes"
                }
              ]
            }
          ]
        },
        {
          "type": "listItem",
          "content": [
            {
              "type": "paragraph",
              "content": [
                {
                  "type": "text",
                  "text": "And "
                },
                {
                  "type": "text",
                  "text": "links",
                  "marks": [
                    {
                      "type": "link",
                      "attrs": { "href": "https://google.com" }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "type": "heading",
      "attrs": {
        "level": 2
      },
      "content": [
        {
          "type": "text",
          "text": "An un-ordered (bullet) list..."
        }
      ]
    },
    {
      "type": "bulletList",
      "content": [
        {
          "type": "listItem",
          "content": [
            {
              "type": "paragraph",
              "content": [
                {
                  "type": "text",
                  "text": "We like bullets"
                }
              ]
            }
          ]
        },
        {
          "type": "listItem",
          "content": [
            {
              "type": "paragraph",
              "content": [
                {
                  "type": "text",
                  "text": "Too"
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "type": "heading",
      "attrs": {
        "level": 3
      },
      "content": [
        {
          "type": "text",
          "text": "And some code..."
        }
      ]
    },
    {
      "type": "codeBlock",
      "attrs": {
        "language": "python"
      },
      "content": [
        {
          "type": "text",
          "text": "x = 12\n"
        }
      ]
    }
  ]
}
    `, string(docJSON))
}
