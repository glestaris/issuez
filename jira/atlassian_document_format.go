package jira

/******************************************************************************
 * Interfaces
 *****************************************************************************/

type ADFTextMode struct {
	Strong bool
	Em     bool
	Code   bool
	Strike bool
}

type ADFNodeText interface {
	AddText(text string, mode ADFTextMode)
	AddLink(text string, url string, mode ADFTextMode)
}

type ADFNodeList interface {
	AddItem() ADFNodeText
}

type ADFHeadingLevel int

const (
	ADFHeadingLevel1 ADFHeadingLevel = iota
	ADFHeadingLevel2
	ADFHeadingLevel3
	ADFHeadingLevel4
	ADFHeadingLevel5
	ADFHeadingLevel6
)

type ADFDocument interface {
	AddParagraph() ADFNodeText
	AddBulletList() ADFNodeList
	AddOrderedList() ADFNodeList
	AddCodeBlock(language string, code string)
	AddHeading(level ADFHeadingLevel, text string)
}

func NewADFDocument() ADFDocument {
	return newADFDocument()
}

/******************************************************************************
 * Concrete
 *****************************************************************************/

type adfDocument struct {
	Version int        `json:"version"`
	Type    string     `json:"type"`
	Content []*adfNode `json:"content"`
}

type adfNode struct {
	Type    string                 `json:"type"`
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
	Content []*adfNode             `json:"content,omitempty"`
	Text    string                 `json:"text,omitempty"`
	Marks   []*adfNode             `json:"marks,omitempty"`
}

func newADFDocument() ADFDocument {
	return &adfDocument{
		Version: 1,
		Type:    "doc",
		Content: []*adfNode{},
	}
}

type textNodeContainer struct {
	*adfNode
}

func (c *textNodeContainer) AddText(text string, mode ADFTextMode) {
	marks := []*adfNode{}
	if mode.Strong {
		marks = append(marks, &adfNode{Type: "strong"})
	}
	if mode.Em {
		marks = append(marks, &adfNode{Type: "em"})
	}
	if mode.Code {
		marks = append(marks, &adfNode{Type: "code"})
	}
	if mode.Strike {
		marks = append(marks, &adfNode{Type: "strike"})
	}

	node := &adfNode{
		Type:  "text",
		Text:  text,
		Marks: marks,
	}
	c.adfNode.Content = append(c.adfNode.Content, node)
}

func (c *textNodeContainer) AddLink(
	text string, linkURL string, mode ADFTextMode,
) {
	marks := []*adfNode{}
	if mode.Strong {
		marks = append(marks, &adfNode{Type: "strong"})
	}
	if mode.Em {
		marks = append(marks, &adfNode{Type: "em"})
	}
	if mode.Code {
		marks = append(marks, &adfNode{Type: "code"})
	}
	if mode.Strike {
		marks = append(marks, &adfNode{Type: "strike"})
	}
	marks = append(marks, &adfNode{
		Type: "link",
		Attrs: map[string]interface{}{
			"href": linkURL,
		},
	})

	node := &adfNode{
		Type:  "text",
		Text:  text,
		Marks: marks,
	}
	c.adfNode.Content = append(c.adfNode.Content, node)
}

func (d *adfDocument) AddParagraph() ADFNodeText {
	p := &adfNode{
		Type:    "paragraph",
		Content: []*adfNode{},
	}
	d.Content = append(d.Content, p)
	return &textNodeContainer{p}
}

func (d *adfDocument) AddHeading(headingLevel ADFHeadingLevel, text string) {
	headingLevelInt := 6
	switch headingLevel {
	case ADFHeadingLevel1:
		headingLevelInt = 1
	case ADFHeadingLevel2:
		headingLevelInt = 2
	case ADFHeadingLevel3:
		headingLevelInt = 3
	case ADFHeadingLevel4:
		headingLevelInt = 4
	case ADFHeadingLevel5:
		headingLevelInt = 5
	}

	h := &adfNode{
		Type: "heading",
		Attrs: map[string]interface{}{
			"level": headingLevelInt,
		},
		Content: []*adfNode{
			{
				Type: "text",
				Text: text,
			},
		},
	}
	d.Content = append(d.Content, h)
}

type listNodeContainer struct {
	*adfNode
}

func (l *listNodeContainer) AddItem() ADFNodeText {
	p := &adfNode{Type: "paragraph"}
	li := &adfNode{
		Type:    "listItem",
		Content: []*adfNode{p},
	}
	l.adfNode.Content = append(l.adfNode.Content, li)
	return &textNodeContainer{p}
}

func (d *adfDocument) AddOrderedList() ADFNodeList {
	l := &adfNode{
		Type:    "orderedList",
		Content: []*adfNode{},
	}
	d.Content = append(d.Content, l)
	return &listNodeContainer{l}
}

func (d *adfDocument) AddBulletList() ADFNodeList {
	l := &adfNode{
		Type:    "bulletList",
		Content: []*adfNode{},
	}
	d.Content = append(d.Content, l)
	return &listNodeContainer{l}
}

func (d *adfDocument) AddCodeBlock(language string, code string) {
	cb := &adfNode{
		Type: "codeBlock",
		Attrs: map[string]interface{}{
			"language": language,
		},
		Content: []*adfNode{
			{
				Type: "text",
				Text: code,
			},
		},
	}
	d.Content = append(d.Content, cb)
}
