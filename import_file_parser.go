package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"
	"github.com/glestaris/issuez/domain"
)

func ParseImportFile(markdownFile io.Reader) ([]*domain.Issue, error) {
	data, err := ioutil.ReadAll(markdownFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read markdown file: %s", err)
	}

	md := blackfriday.New(blackfriday.WithExtensions(
		blackfriday.FencedCode | blackfriday.Strikethrough,
	))
	node := md.Parse(data)

	// parsing did not produce a doc, no issues
	if node == nil {
		return []*domain.Issue{}, nil
	}

	// make document
	doc, err := newDocument(node)
	if err != nil {
		log.Printf("Failed to create document: %s", err)
		return nil, errors.New("Failed to parse markdown file")
	}

	// find sections
	sections, err := doc.sections()
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to extract issues from markdown file: %s", err,
		)
	}
	if len(sections) == 0 {
		return []*domain.Issue{}, nil
	}

	// create issues
	issues := make([]*domain.Issue, len(sections))
	for i, section := range sections {
		issue, err := section.makeIssue()
		if err != nil {
			return nil, fmt.Errorf(
				"Failed parsing issue %d in markdown file: %s", i+1, err,
			)
		}
		issues[i] = issue
	}

	return issues, nil
}

type document struct {
	root *blackfriday.Node
}

func newDocument(doc *blackfriday.Node) (*document, error) {
	if doc.Type != blackfriday.Document {
		return nil, errors.New("Document not found")
	}
	return &document{root: doc}, nil
}

func isNodeEmpty(node *blackfriday.Node) bool {
	// nil node is empty
	if node == nil {
		return true
	}

	// doesn't have children: empty if literal is empty
	if node.FirstChild == nil {
		return len(node.Literal) == 0
	}

	// has children: empty if all children are empty
	foundNonEmpty := false
	node.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			return blackfriday.GoToNext
		}

		if n == node {
			return blackfriday.GoToNext
		}
		if !isNodeEmpty(n) {
			foundNonEmpty = true
			return blackfriday.Terminate
		}
		return blackfriday.SkipChildren
	})
	return !foundNonEmpty
}

func (d *document) sections() ([]*section, error) {
	boundaries := []*blackfriday.Node{}
	currNode := d.root.FirstChild
	for currNode != nil {
		// remove HRs
		if currNode.Type == blackfriday.HorizontalRule {
			// mark prev as boundary
			// IFF
			//  1) prev != nil: not the first node
			//  2) next != nil: not the last node
			if currNode.Prev != nil && currNode.Next != nil {
				b := currNode.Prev
				// boundary does not already exist in list
				if len(boundaries) == 0 ||
					boundaries[len(boundaries)-1] != b {
					boundaries = append(boundaries, b)
				}
			}

			nextNode := currNode.Next
			currNode.Unlink()
			currNode = nextNode
			continue
		}

		// remove empty nodes
		if isNodeEmpty(currNode) {
			nextNode := currNode.Next
			currNode.Unlink()
			currNode = nextNode
			continue
		}

		currNode = currNode.Next
	}

	if d.root.FirstChild == nil {
		return []*section{}, nil
	}

	var sections []*section
	// no boundaries? - single issue
	if len(boundaries) == 0 {
		sections = []*section{
			{
				firstNode: d.root.FirstChild,
				lastNode:  d.root.LastChild,
			},
		}
	} else {
		sections = make([]*section, len(boundaries)+1)
		// first section
		sections[0] = &section{
			firstNode: d.root.FirstChild,
			lastNode:  boundaries[0],
		}
		// sections in between
		for sectionIdx := 1; sectionIdx < len(boundaries); sectionIdx++ {
			sections[sectionIdx] = &section{
				firstNode: boundaries[sectionIdx-1].Next,
				lastNode:  boundaries[sectionIdx],
			}
		}
		// last section
		sections[len(boundaries)] = &section{
			firstNode: boundaries[len(boundaries)-1].Next,
			lastNode:  d.root.LastChild,
		}
	}

	return sections, nil
}

type section struct {
	firstNode *blackfriday.Node
	lastNode  *blackfriday.Node
}

func (s *section) makeIssue() (*domain.Issue, error) {
	// parse header
	issueType, title, err := s.parseHeader()
	if err != nil {
		return nil, err
	}

	// parse footer
	epicID, labels := s.parseFooter()

	// parse description
	var incLastNode bool
	if epicID == "" && labels == nil {
		// last node was not used as footer
		incLastNode = true
	}
	description, err := s.parseDescription(incLastNode)
	if err != nil {
		return nil, err
	}

	// make issue
	issue := &domain.Issue{}

	// issue title
	issue.Title = title

	// issue type
	if issueType == "" || issueType == "Story" || issueType == "Issue" {
		issue.Type = domain.IssueTypeStory
	} else if issueType == "Bug" {
		issue.Type = domain.IssueTypeBug
	} else if issueType == "Chore" || issueType == "Task" {
		issue.Type = domain.IssueTypeChore
	} else {
		return nil, fmt.Errorf("Unknown issue type %s", issueType)
	}

	// issue description
	issue.Description = description

	// issue epic
	if epicID != "" {
		issue.Epic = &domain.Epic{ID: epicID}
	}

	// issue labels
	if labels != nil && len(labels) != 0 {
		for _, label := range labels {
			issue.Labels = append(issue.Labels, domain.Label{
				Label: label,
			})
		}
	}

	return issue, nil
}

func (s *section) parseHeader() (string, string, error) {
	f := s.firstNode
	if f.Type != blackfriday.Paragraph ||
		f.FirstChild == nil ||
		f.FirstChild != f.LastChild ||
		f.FirstChild.Type != blackfriday.Text {
		return "", "", errors.New(
			"First line in issue section needs to be of the form" +
				" '[ISSUE TYPE] ISSUE TITLE'",
		)
	}
	firstLine := string(f.FirstChild.Literal)

	re := regexp.MustCompile(`^\s*(?:\[([^\[\]]+)\])?\s*(.+)\s*$`)
	matches := re.FindStringSubmatch(firstLine)
	if len(matches) == 2 {
		return "", strings.TrimSpace(matches[1]), nil
	}
	if len(matches) == 3 {
		return strings.TrimSpace(matches[1]),
			strings.TrimSpace(matches[2]),
			nil
	}

	return "", "", errors.New(
		"First line in issue section needs to be of the form" +
			" '[ISSUE TYPE] ISSUE TITLE'",
	)
}

func (s *section) parseFooter() (string, []string) {
	l := s.lastNode
	if l.Type != blackfriday.Paragraph ||
		l.FirstChild == nil ||
		l.FirstChild != l.LastChild ||
		l.FirstChild.Type != blackfriday.Text {
		return "", nil
	}
	lastParagraph := string(l.FirstChild.Literal)

	var epicID string
	epicRe := regexp.MustCompile(`(?:E|Epic):\s*(.+)`)
	epicReMatches := epicRe.FindStringSubmatch(lastParagraph)
	if len(epicReMatches) == 2 {
		epicID = strings.TrimSpace(epicReMatches[1])
	}

	var labels []string
	labelsRe := regexp.MustCompile(`(?:L|Labels):\s*(.+)`)
	labelsReMatches := labelsRe.FindStringSubmatch(lastParagraph)
	if len(labelsReMatches) == 2 {
		labels = []string{}
		for _, label := range strings.Split(labelsReMatches[1], ",") {
			labels = append(labels, strings.TrimSpace(label))
		}
	}

	return epicID, labels
}

func parseTextContainer(node *blackfriday.Node, tc *domain.TextContainer) {
	textMode := domain.TextMode{}
	linkURL := ""
	node.Walk(func(
		in *blackfriday.Node, entering bool,
	) blackfriday.WalkStatus {
		switch in.Type {
		case blackfriday.Strong:
			textMode.Bold = !textMode.Bold
		case blackfriday.Emph:
			textMode.Italics = !textMode.Italics
		case blackfriday.Del:
			textMode.Strikethrough = !textMode.Strikethrough
		case blackfriday.Link:
			if entering {
				linkURL = string(in.LinkData.Destination)
			} else {
				linkURL = ""
			}

			// Leafs
		case blackfriday.Code:
			text := string(in.Literal)
			if text == "" {
				return blackfriday.GoToNext
			}

			textMode.Code = true
			if linkURL == "" {
				tc.AddText(text, textMode)
			} else {
				tc.AddLink(text, linkURL, textMode)
			}
			textMode.Code = false
		case blackfriday.Text:
			text := string(in.Literal)
			if text == "" {
				return blackfriday.GoToNext
			}

			if linkURL == "" {
				tc.AddText(text, textMode)
			} else {
				tc.AddLink(text, linkURL, textMode)
			}
		}

		return blackfriday.GoToNext
	})
}

func parseText(node *blackfriday.Node) string {
	text := ""
	node.Walk(func(
		in *blackfriday.Node, entering bool,
	) blackfriday.WalkStatus {
		text += string(in.Literal)
		return blackfriday.GoToNext
	})
	return text
}

func headingLevel(nodeHeadingLevel int) domain.HeadingLevel {
	switch nodeHeadingLevel {
	case 1:
		return domain.HeadingLevel1
	case 2:
		return domain.HeadingLevel2
	case 3:
		return domain.HeadingLevel3
	case 4:
		return domain.HeadingLevel4
	case 5:
		return domain.HeadingLevel5
	default:
		return domain.HeadingLevel5
	}
}

func (s *section) parseDescription(incLastNode bool) (*domain.Document, error) {
	if s.firstNode == s.lastNode {
		// no description
		return nil, nil
	}

	startNode := s.firstNode.Next
	stopNode := s.lastNode.Prev
	if incLastNode {
		stopNode = s.lastNode
	}
	domainDoc := &domain.Document{}
	for node := startNode; node != stopNode.Next; node = node.Next {
		switch node.Type {
		case blackfriday.Paragraph:
			tc := domainDoc.AddParagraph()
			parseTextContainer(node, tc)

		case blackfriday.List:
			var list *domain.ListData
			if node.ListData.ListFlags&blackfriday.ListTypeOrdered ==
				blackfriday.ListTypeOrdered {
				list = domainDoc.AddOrderedList()
			} else {
				list = domainDoc.AddUnorderedList()
			}
			for in := node.FirstChild; in != nil; in = in.Next {
				tc := list.AddItem()
				parseTextContainer(in, tc)
			}

		case blackfriday.CodeBlock:
			domainDoc.AddCodeBlock(
				string(node.CodeBlockData.Info), string(node.Literal),
			)

		case blackfriday.Heading:
			text := parseText(node)
			domainDoc.AddHeading(headingLevel(node.HeadingData.Level), text)

		default:
			return nil, fmt.Errorf("Unknown node type: %s", node.Type)
		}
	}
	return domainDoc, nil
}
