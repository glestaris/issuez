package domain

type IssueType int

const (
	IssueTypeChore IssueType = iota
	IssueTypeStory
	IssueTypeBug
)

func (it IssueType) String() string {
	switch it {
	case IssueTypeChore:
		return "Chore"
	case IssueTypeStory:
		return "User Story"
	case IssueTypeBug:
		return "Bug"
	default:
		return "Unknown Type"
	}
}

type Tracker struct {
	Type   string
	Config map[string]string
}

type Epic struct {
	ID string
}

type Label struct {
	Label string
}

type Issue struct {
	ID          string
	Type        IssueType
	Title       string
	Description *Document
	Epic        *Epic
	Labels      []Label
}
