package main

import (
	"bytes"

	"github.com/andygrunwald/go-jira"
)

//Worklog request represents the request object used to make request to JIRA API.
//This is a data type is a minimal subset of the Worklog type included in the go-jira
//library. The included data type JSON representation cannot be used to create a
//worklog entry with only the fields below.
type WorklogRequest struct {
	TimeSpentSeconds string `json:"timeSpentSeconds"`
	Started          string `json:"started"`
	Comment          string `json:"comment"`
}

//Calls JIRA worklog endpoint to create a worklog entry for an issue.
//Implemented using generic interface provided by go-jira due to this endpoint not being
//implemented at the time of implementation.
func addWorklog(jiraClient *jira.Client, issueKey string, worklog WorklogRequest) error {
	req, _ := jiraClient.NewRequest("POST", "/rest/api/2/issue/"+issueKey+"/worklog", worklog)
	resp, err := jiraClient.Do(req, nil)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return err
}
