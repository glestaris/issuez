package domain

func NewDocument() Document {
	return Document{}
}

type Document struct {
	Nodes []DocumentNode
}

type DocumentNodeType int

const (
	// Paragraph
	DocumentNodeTypeParagraph DocumentNodeType = iota
	// Headings
	DocumentNodeTypeHeading
	// Lists
	DocumentNodeTypeList
	// Code Block
	DocumentNodeTypeCodeBlock
)

func (t DocumentNodeType) String() string {
	switch t {
	case DocumentNodeTypeHeading:
		return "Heading"
	case DocumentNodeTypeList:
		return "List"
	case DocumentNodeTypeParagraph:
		return "Paragraph"
	case DocumentNodeTypeCodeBlock:
		return "Code Block"
	default:
		return "Unknown"
	}
}

type DocumentNode struct {
	Type DocumentNodeType

	// Paragraph
	*ParagraphData

	// Heading
	*HeadingData

	// Lists
	*ListData

	// Code block
	*CodeBlockData
}

/******************************************************************************
 * Headings
 *****************************************************************************/

type HeadingLevel int

const (
	HeadingLevel1 HeadingLevel = iota
	HeadingLevel2
	HeadingLevel3
	HeadingLevel4
	HeadingLevel5
)

type HeadingData struct {
	Level HeadingLevel
	Text  string
}

func (d *Document) AddHeading(headingLevel HeadingLevel, text string) {
	node := DocumentNode{
		Type: DocumentNodeTypeHeading,
		HeadingData: &HeadingData{
			Level: headingLevel,
			Text:  text,
		},
	}
	d.Nodes = append(d.Nodes, node)
}

/******************************************************************************
 * Code Block
 *****************************************************************************/

type CodeBlockData struct {
	Language string
	Code     string
}

func (d *Document) AddCodeBlock(language string, code string) {
	node := DocumentNode{
		Type: DocumentNodeTypeCodeBlock,
		CodeBlockData: &CodeBlockData{
			Language: language,
			Code:     code,
		},
	}
	d.Nodes = append(d.Nodes, node)
}

/******************************************************************************
 * Text Container
 *****************************************************************************/

type TextContainer struct {
	Elements []TextElement
}

type TextMode struct {
	Bold          bool
	Italics       bool
	Strikethrough bool
	Code          bool
}

type TextElement struct {
	Text string
	Mode TextMode

	// Links
	LinkURL string
}

func (t *TextContainer) AddText(text string, mode TextMode) {
	te := TextElement{
		Text: text,
		Mode: mode,
	}
	t.Elements = append(t.Elements, te)
}

func (t *TextContainer) AddLink(text string, linkURL string, mode TextMode) {
	te := TextElement{
		Text:    text,
		Mode:    mode,
		LinkURL: linkURL,
	}
	t.Elements = append(t.Elements, te)
}

/******************************************************************************
 * Paragraphs
 *****************************************************************************/

type ParagraphData struct {
	Content TextContainer
}

func (d *Document) AddParagraph() *TextContainer {
	data := &ParagraphData{}
	node := DocumentNode{
		Type:          DocumentNodeTypeParagraph,
		ParagraphData: data,
	}
	d.Nodes = append(d.Nodes, node)
	return &data.Content
}

/******************************************************************************
 * Lists
 *****************************************************************************/

type ListData struct {
	IsOrdered bool
	Items     []TextContainer
}

func (d *Document) AddOrderedList() *ListData {
	data := &ListData{
		IsOrdered: true,
	}
	node := DocumentNode{
		Type:     DocumentNodeTypeList,
		ListData: data,
	}
	d.Nodes = append(d.Nodes, node)
	return data
}

func (d *Document) AddUnorderedList() *ListData {
	data := &ListData{
		IsOrdered: false,
	}
	node := DocumentNode{
		Type:     DocumentNodeTypeList,
		ListData: data,
	}
	d.Nodes = append(d.Nodes, node)
	return data
}

func (l *ListData) AddItem() *TextContainer {
	l.Items = append(l.Items, TextContainer{})
	return &l.Items[len(l.Items)-1]
}
