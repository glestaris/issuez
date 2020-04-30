package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
	// host
	host string
	// auth
	username string
	token    string
}

func NewJiraClient(
	apiHost string, apiUsername string, apiToken string,
	httpClient *http.Client,
) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
		host:       apiHost,
		username:   apiUsername,
		token:      apiToken,
	}
}

/******************************************************************************
 * Import JIRA Issues
 *****************************************************************************/

type IssueType int

const (
	IssueTypeTask IssueType = iota
	IssueTypeStory
	IssueTypeBug
)

type Issue struct {
	ProjectKey  string
	Type        IssueType
	Summary     string
	Description ADFDocument
	EpicKey     string
	Labels      []string
}

func (i Issue) String() string {
	switch i.Type {
	case IssueTypeBug:
		return "Bug"

	case IssueTypeStory:
		return "Story"

	case IssueTypeTask:
		return "Task"

	default:
		return "Unknown"
	}
}

type ImportIssuesResponse []struct {
	NewIssueKey string
	Err         error
}

type issImpReqIssue struct {
	Fields struct {
		Project struct {
			Key string `json:"key"`
		} `json:"project"`
		IssueType struct {
			Name string `json:"name"`
		} `json:"issuetype"`
		Summary     string      `json:"summary"`
		Description ADFDocument `json:"description,omitempty"`
		Parent      struct {
			Key string `json:"key,omitempty"`
		} `json:"parent,omitempty"`
		Labels []string `json:"labels"`
	} `json:"fields"`
}

type issImpReq struct {
	Issues []issImpReqIssue `json:"issueUpdates"`
}

type issImpRespIssue struct {
	Key string `json:"key"`
}

type issImpRespError struct {
	Status           int `json:"status"`
	FailedElementIdx int `json:"failedElementNumber"`
	ElementErrors    struct {
		ErrorMessages []string          `json:"errorMessages"`
		Errors        map[string]string `json:"errors"`
	} `json:"elementErrors"`
}

type issImpResp struct {
	Issues []issImpRespIssue `json:"issues"`
	Errors []issImpRespError `json:"errors"`
}

func (c *Client) ImportIssues(
	issues []*Issue,
) (ImportIssuesResponse, error) {
	if len(issues) == 0 {
		return ImportIssuesResponse{}, nil
	}

	reqBody := issImpReq{
		Issues: make([]issImpReqIssue, len(issues)),
	}
	for i, issue := range issues {
		reqBodyIssue := issImpReqIssue{}
		reqBodyIssue.Fields.Project.Key = issue.ProjectKey
		reqBodyIssue.Fields.IssueType.Name = issue.String()
		reqBodyIssue.Fields.Summary = issue.Summary
		reqBodyIssue.Fields.Description = issue.Description
		if issue.EpicKey != "" {
			reqBodyIssue.Fields.Parent.Key = issue.EpicKey
		}
		reqBodyIssue.Fields.Labels = issue.Labels

		reqBody.Issues[i] = reqBodyIssue
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize request body: %s", err)
	}

	req, resp, err := c.performRequest(
		"POST", "/rest/api/3/issue/bulk", reqBodyBytes,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request: %s", err)
	}
	if resp.StatusCode != 201 && resp.StatusCode != 400 {
		c.logFailedRequest(req, resp)
		return nil, fmt.Errorf(
			"Failed to create issues: %s", resp.Status,
		)
	}

	respBody := issImpResp{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse API response: %s", err)
	}
	retVal := make(ImportIssuesResponse, len(issues))
	for _, respErr := range respBody.Errors {
		respErrBytes, _ := json.Marshal(respErr.ElementErrors)
		log.Printf(
			"Failed to process issue %d: %v", respErr.FailedElementIdx + 1,
			string(respErrBytes),
		)
		retVal[respErr.FailedElementIdx].Err = errors.New(
			"Failed to process issue",
		)
	}
	i := 0
	for _, respIssue := range respBody.Issues {
		for retVal[i].Err != nil {
			i++
		}
		retVal[i].NewIssueKey = respIssue.Key
		i++
	}

	return retVal, nil
}

/******************************************************************************
 * Test JIRA API Connection
 *****************************************************************************/

func (c *Client) Test() error {
	req, resp, err := c.performRequest("GET", "/rest/api/3/project", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		c.logFailedRequest(req, resp)
		return fmt.Errorf(
			"Failed to test the JIRA API connection: %s", resp.Status,
		)
	}

	return nil
}

/******************************************************************************
 * JIRA API Helpers
 *****************************************************************************/

func (c *Client) performRequest(
	method string, path string, body []byte,
) (*http.Request, *http.Response, error) {
	reqURL := fmt.Sprintf(
		"%s/%s",
		strings.TrimRight(c.host, "/"),
		strings.TrimLeft(path, "/"),
	)
	bodyFile := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, reqURL, bodyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create request: %s", err)
	}
	req.SetBasicAuth(c.username, c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to perform request: %s", err)
	}

	return req, resp, nil
}

func (c *Client) logFailedRequest(req *http.Request, resp *http.Response) {
	log.Printf("JIRA API request %s %s failed", req.Method, req.URL.String())

	log.Printf("\tStatus: '%s'", resp.Status)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		respBodyStr := string(respBody)
		log.Printf("\tResponse body: '%s'", respBodyStr)
	} else {
		log.Printf("\tFailed read response body: %s", err)
	}
}
