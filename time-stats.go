package main

import (
	"time"
	"github.com/cbergoon/gocolor"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"sort"
	"strconv"
)

type TimeStats struct {
	TotalTime  int
	TimeBlocks []WorkLogItem
}

type WorkLogItem struct {
	Start            time.Time
	End              time.Time
	TimeSpent        string
	TimeSpentSeconds int
}

func printTimeStats(timeStats *TimeStats) {
	totalTime := time.Duration(time.Second * time.Duration(timeStats.TotalTime))
	totalTimeString := gocolor.ColorClear(gocolor.COLOR_RED, totalTime.String())
	totalEntriesString := gocolor.ColorClear(gocolor.COLOR_RED, strconv.Itoa(len(timeStats.TimeBlocks)))

	fmt.Printf("Total Time: %s Entries: %s\n", totalTimeString, totalEntriesString)
	for _, block := range timeStats.TimeBlocks {
		timeSpentString := gocolor.ColorClear(gocolor.COLOR_RED, stringPadRight(block.TimeSpent, " ", 10))
		fmt.Printf("%s %s to %s\n", timeSpentString, block.Start.Format(DISPLAY_DATE_FORMAT), block.End.Format(DISPLAY_DATE_FORMAT))
	}
}

func generateTimeStats(issues []jira.Issue, after func() time.Time, user string) *TimeStats {
	stats := &TimeStats{}
	for _, issue := range issues {
		for _, worklog := range issue.Fields.Worklog.Worklogs {
			if worklog.Author.Name == user {
				if time.Time(worklog.Started).After(after()) {
					stats.TotalTime += worklog.TimeSpentSeconds
					stats.TimeBlocks = append(stats.TimeBlocks, WorkLogItem{
						TimeSpentSeconds: worklog.TimeSpentSeconds,
						TimeSpent:        worklog.TimeSpent,
						Start:            time.Time(worklog.Started),
						End:              time.Time(worklog.Started).Add(time.Second * time.Duration(worklog.TimeSpentSeconds)),
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
