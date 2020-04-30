package tracker

import (
	"fmt"

	"github.com/glestaris/issuez/domain"
)

// APP layer
//  Test using integration tests

type TrackerService interface {
	ImportIssues(issues []*domain.Issue) error
	TestConnection() error
}

func NewTrackerService(tracker domain.Tracker) (TrackerService, error) {
	if tracker.Type == "jira" {
		return newJiraTrackerService(
			tracker.Config["apiHost"],
			tracker.Config["apiUsername"],
			tracker.Config["apiToken"],
			tracker.Config["projectKey"],
		), nil
	}

	return nil, fmt.Errorf("Unknown tracker type '%s'", tracker.Type)
}
