package main

import (
	"errors"

	"github.com/andygrunwald/go-jira"
)

func configureNewJiraClient(url, username, password string) (*jira.Client, error) {
	var err error
	jiraClient, err := jira.NewClient(nil, url)
	if err != nil {
		return nil, err
	}
	ok, err := jiraClient.Authentication.AcquireSessionCookie(username, password)
	if !ok {
		return nil, errors.New("failed to create jira client; authentication failed")
	}
	if err != nil {
		return nil, err
	}
	return jiraClient, nil
}
