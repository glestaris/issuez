package tracker

import (
	"fmt"

	"github.com/glestaris/issuez/domain"
	"github.com/glestaris/issuez/jira"
)

func mapTextMode(tceMode domain.TextMode) (jtMode jira.ADFTextMode) {
	jtMode.Strong = tceMode.Bold
	jtMode.Code = tceMode.Code
	jtMode.Em = tceMode.Italics
	jtMode.Strike = tceMode.Strikethrough
	return
}

func mapTextContainer(jt jira.ADFNodeText, tc domain.TextContainer) {
	for _, tce := range tc.Elements {
		jtMode := mapTextMode(tce.Mode)
		if tce.LinkURL != "" {
			jt.AddLink(tce.Text, tce.LinkURL, jtMode)
		} else {
			jt.AddText(tce.Text, jtMode)
		}
	}
}

func addParagraph(jd jira.ADFDocument, node domain.DocumentNode) {
	p := jd.AddParagraph()
	mapTextContainer(p, node.Content)
}

func mapHeadingLevel(hl domain.HeadingLevel) (jhl jira.ADFHeadingLevel) {
	switch hl {
	case domain.HeadingLevel1:
		jhl = jira.ADFHeadingLevel1
	case domain.HeadingLevel2:
		jhl = jira.ADFHeadingLevel2
	case domain.HeadingLevel3:
		jhl = jira.ADFHeadingLevel3
	case domain.HeadingLevel4:
		jhl = jira.ADFHeadingLevel4
	case domain.HeadingLevel5:
		jhl = jira.ADFHeadingLevel5
	}
	return
}

func addHeading(jd jira.ADFDocument, node domain.DocumentNode) {
	jd.AddHeading(
		mapHeadingLevel(node.HeadingData.Level), node.HeadingData.Text,
	)
}

func addList(jd jira.ADFDocument, node domain.DocumentNode) {
	var l jira.ADFNodeList
	if node.ListData.IsOrdered {
		l = jd.AddOrderedList()
	} else {
		l = jd.AddOrderedList()
	}
	for _, tc := range node.ListData.Items {
		jt := l.AddItem()
		mapTextContainer(jt, tc)
	}
}

func addCodeBlock(jd jira.ADFDocument, node domain.DocumentNode) {
	jd.AddCodeBlock(node.CodeBlockData.Language, node.CodeBlockData.Code)
}

func mapDocument(domainDoc *domain.Document) (jira.ADFDocument, error) {
	if domainDoc == nil {
		return nil, nil
	}

	jd := jira.NewADFDocument()
	for _, node := range domainDoc.Nodes {
		switch node.Type {
		case domain.DocumentNodeTypeParagraph:
			addParagraph(jd, node)
		case domain.DocumentNodeTypeList:
			addList(jd, node)
		case domain.DocumentNodeTypeHeading:
			addHeading(jd, node)
		case domain.DocumentNodeTypeCodeBlock:
			addCodeBlock(jd, node)
		default:
			return nil, fmt.Errorf("Cannot map document node type %s to Jira",
				node.Type)
		}
	}

	return jd, nil
}
