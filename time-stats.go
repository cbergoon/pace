package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/cbergoon/gocolor"
)

type TimeStats struct {
	TotalTime  int
	TimeBlocks []WorklogItem
}

type WorklogItem struct {
	Start             time.Time
	End               time.Time
	TimeSpent         string
	TimeSpentSeconds  int
	Author            string
	IssueKey          string
	IssueSummary      string
	AddsTimeToKey     string
	AddsTimeToSummary string
}

func printTimeStats(timeStats *TimeStats) {
	totalTime := time.Duration(time.Second * time.Duration(timeStats.TotalTime))
	totalTimeString := gocolor.ColorClear(gocolor.COLOR_RED, totalTime.String())
	totalEntriesString := gocolor.ColorClear(gocolor.COLOR_RED, strconv.Itoa(len(timeStats.TimeBlocks)))

	fmt.Printf("Total Time: %s Entries: %s\n", totalTimeString, totalEntriesString)
	for _, block := range timeStats.TimeBlocks {
		timeSpentString := gocolor.ColorClear(gocolor.COLOR_RED, stringPadRight(block.TimeSpent, " ", 10))
		issueKeyString := stringPadRight(block.IssueKey+":", " ", 16)
		authorString := stringPadRight(block.Author, " ", 28)
		linkString := ""
		if block.AddsTimeToKey != "" {
			linkString = fmt.Sprint(" ->  ", stringPadRight(block.AddsTimeToKey+":", " ", 16), block.AddsTimeToSummary)
		}
		fmt.Printf("%s %s %s to %s \t %s %s %s \n", authorString, timeSpentString, block.Start.Format(DISPLAY_DATE_FORMAT), block.End.Format(DISPLAY_DATE_FORMAT), issueKeyString, block.IssueSummary, linkString)
	}
}

func generateTimeStats(issues []jira.Issue, after func() time.Time, user string) *TimeStats {
	stats := &TimeStats{}
	for _, issue := range issues {
		if len(issue.Fields.IssueLinks) > 0 {
			b, _ := json.MarshalIndent(issue, " ", "\t")
			fmt.Println(string(b))
		}
		for _, worklog := range issue.Fields.Worklog.Worklogs {
			if worklog.Author.Name == user {
				if time.Time(worklog.Started).After(after()) {
					addsTimeToKey := ""
					addsTimeToSummary := ""
					for i := 0; i < len(issue.Fields.IssueLinks); i++ {
						if issue.Fields.IssueLinks[i].Type.Name == "Adds Time" && issue.Fields.IssueLinks[i].OutwardIssue != nil {
							addsTimeToKey = issue.Fields.IssueLinks[i].OutwardIssue.Key
							addsTimeToSummary = issue.Fields.IssueLinks[i].OutwardIssue.Fields.Summary
						}
					}
					stats.TotalTime += worklog.TimeSpentSeconds
					stats.TimeBlocks = append(stats.TimeBlocks, WorklogItem{
						TimeSpentSeconds:  worklog.TimeSpentSeconds,
						TimeSpent:         worklog.TimeSpent,
						Start:             time.Time(worklog.Started),
						End:               time.Time(worklog.Started).Add(time.Second * time.Duration(worklog.TimeSpentSeconds)),
						Author:            worklog.Author.Name,
						IssueKey:          issue.Key,
						IssueSummary:      issue.Fields.Summary,
						AddsTimeToKey:     addsTimeToKey,
						AddsTimeToSummary: addsTimeToSummary,
					})
				}
			}
		}
	}
	sort.Slice(stats.TimeBlocks, func(i, j int) bool {
		return stats.TimeBlocks[i].Start.Before(stats.TimeBlocks[j].Start)
	})
	return stats
}
