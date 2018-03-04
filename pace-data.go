package main

import (
	"github.com/c-bata/go-prompt"
	"github.com/andygrunwald/go-jira"
	"strconv"
	"time"
)

type PaceData struct {
	IssueSuggestions []prompt.Suggest
	TimeSuggestions  []prompt.Suggest

	UnresolvedIssues []jira.Issue

	DayIssues   []jira.Issue
	WeekIssues  []jira.Issue
	MonthIssues []jira.Issue
}

func NewPaceData() *PaceData {
	return &PaceData{
		TimeSuggestions: buildTimeSuggestionTable(),
	}
}

func (paceData *PaceData) loadIssueQueues(jiraClient *jira.Client, config *Config) error {
	paceData.IssueSuggestions = paceData.IssueSuggestions[:0]
	paceData.UnresolvedIssues = paceData.UnresolvedIssues[:0]
	paceData.DayIssues = paceData.DayIssues[:0]
	paceData.WeekIssues = paceData.WeekIssues[:0]
	paceData.MonthIssues = paceData.MonthIssues[:0]

	err := jiraClient.Issue.SearchPages(config.SuggestionIssueQuery, &jira.SearchOptions{Fields: []string{"worklog", "summary"}}, func(i jira.Issue) error {
		paceData.UnresolvedIssues = append(paceData.UnresolvedIssues, i)
		return nil
	})
	if err != nil {
		return err
	}
	err = checkAndRetrievePaginatedWorklogs(jiraClient, paceData.UnresolvedIssues)
	if err != nil {
		return err
	}

	err = jiraClient.Issue.SearchPages(config.DayWorklogIssueQuery, &jira.SearchOptions{Fields: []string{"worklog", "summary"}}, func(i jira.Issue) error {
		paceData.DayIssues = append(paceData.DayIssues, i)
		return nil
	})
	if err != nil {
		return err
	}
	err = checkAndRetrievePaginatedWorklogs(jiraClient, paceData.DayIssues)
	if err != nil {
		return err
	}

	err = jiraClient.Issue.SearchPages(config.WeekWorklogIssueQuery, &jira.SearchOptions{Fields: []string{"worklog", "summary"}}, func(i jira.Issue) error {
		paceData.WeekIssues = append(paceData.WeekIssues, i)
		return nil
	})
	if err != nil {
		return err
	}
	err = checkAndRetrievePaginatedWorklogs(jiraClient, paceData.WeekIssues)
	if err != nil {
		return err
	}

	err = jiraClient.Issue.SearchPages(config.MonthWorklogIssueQuery, &jira.SearchOptions{Fields: []string{"worklog", "summary"}}, func(i jira.Issue) error {
		paceData.MonthIssues = append(paceData.MonthIssues, i)
		return nil
	})
	if err != nil {
		return err
	}
	err = checkAndRetrievePaginatedWorklogs(jiraClient, paceData.MonthIssues)
	if err != nil {
		return err
	}

	paceData.IssueSuggestions = buildIssueSuggestionTable(paceData.UnresolvedIssues)

	return nil
}

func (paceData *PaceData) findLastWorklogEndTimeForUserForDate(date time.Time, user string) time.Time {
	month := findLastWorklogEndTimeForUserForDate(paceData.MonthIssues, date, user)
	week := findLastWorklogEndTimeForUserForDate(paceData.WeekIssues, date, user)
	day := findLastWorklogEndTimeForUserForDate(paceData.DayIssues, date, user)
	return maximumTime(month, week, day)
}

func findLastWorklogEndTimeForUserForDate(issues []jira.Issue, date time.Time, user string) time.Time {
	var maxForDate time.Time
	truncatedDate := date.Truncate(time.Hour * 24)
	for _, issue := range issues {
		for _, worklog := range issue.Fields.Worklog.Worklogs {
			if worklog.Author.Name == user {
				started := time.Time(worklog.Started)
				ended := started.Add(time.Duration(worklog.TimeSpentSeconds) * time.Second)
				if ended.After(truncatedDate) && ended.Before(truncatedDate.Add(time.Duration(24)*time.Hour)) {
					if ended.After(maxForDate) {
						maxForDate = ended
					}
				}
			}
		}
	}
	return maxForDate
}

func checkAndRetrievePaginatedWorklogs(jiraClient *jira.Client, issues []jira.Issue) error {
	for _, i := range issues {
		if len(i.Fields.Worklog.Worklogs) < i.Fields.Worklog.Total {
			worklog, _, err := jiraClient.Issue.GetWorklogs(i.Key)
			if err != nil {
				return err
			}
			i.Fields.Worklog = worklog
		}
	}
	return nil
}

func buildIssueSuggestionTable(issues []jira.Issue) []prompt.Suggest {
	var suggestions []prompt.Suggest
	for _, i := range issues {
		suggestions = append(suggestions, prompt.Suggest{Text: i.Key, Description: i.Fields.Summary})
	}
	return suggestions
}

func buildTimeSuggestionTable() []prompt.Suggest {
	var timeTable []prompt.Suggest
	for i := 0; i < 24; i++ {
		for j := 0; j < 60; j++ {
			if i == 0 {
				timeTable = append(timeTable, prompt.Suggest{Text: strconv.Itoa(j) + "m", Description: strconv.Itoa(j) + " minutes"})
			} else {
				timeTable = append(timeTable, prompt.Suggest{Text: strconv.Itoa(i) + "h" + strconv.Itoa(j) + "m", Description: strconv.Itoa(i) + " hours " + strconv.Itoa(j) + " minutes"})
			}
		}
	}
	return timeTable
}
