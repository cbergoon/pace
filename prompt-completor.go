package main

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

func paceCompletor(paceData *PaceData) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		fields := strings.Fields(d.Text)
		if len(fields) > 0 {
			if strings.Compare("time", fields[0]) == 0 {
				if len(fields) > 1 {
					if strings.Compare("today", fields[1]) == 0 {
						return nil
					} else if strings.Compare("week", fields[1]) == 0 {
						return nil
					} else if strings.Compare("month", fields[1]) == 0 {
						return nil
					}
				}
				s := []prompt.Suggest{
					{Text: "today", Description: "time logged today"},
					{Text: "week", Description: "time logged this week"},
					{Text: "month", Description: "time logged this month"},
				}
				return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
			} else if strings.Compare("clock", fields[0]) == 0 {
				if len(fields) > 1 {
					if strings.Compare("start", fields[1]) == 0 {
						return nil
					} else if strings.Compare("stop", fields[1]) == 0 {
						if strings.Compare(string(d.TextBeforeCursor()[len(d.TextBeforeCursor())-1]), " ") == 0 {
							if len(fields) > 2 {
								return nil
							}
						}
						return prompt.FilterContains(paceData.IssueSuggestions, d.GetWordBeforeCursor(), true)
					} else if strings.Compare("pause", fields[1]) == 0 {
						return nil
					}
				}
				s := []prompt.Suggest{
					{Text: "start", Description: "start the clock"},
					{Text: "stop", Description: "stop the clock"},
				}
				return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
			} else if strings.Compare("log", fields[0]) == 0 {
				if (len(fields) == 2 && d.GetCharRelativeToCursor(0) == 32) || (len(fields) == 3 && d.GetCharRelativeToCursor(0) != 32) {
					return prompt.FilterContains(paceData.TimeSuggestions, d.GetWordBeforeCursor(), true)
				}
				if (len(fields) == 3 && d.GetCharRelativeToCursor(0) == 32) || (len(fields) == 4 && d.GetCharRelativeToCursor(0) != 32) {
					return []prompt.Suggest{
						{Text: "--fill", Description: "fill current day from previous or default time"},
						{Text: "YYYYMMDD", Description: "date"},
					}
				}
				if (len(fields) >= 4 && d.GetCharRelativeToCursor(0) == 32) || (len(fields) == 5 && d.GetCharRelativeToCursor(0) != 32) {
					//Terminal Catch
					if len(fields) >= 5 && d.GetCharRelativeToCursor(0) == 32 {
						return nil
					}
					if fields[3] == "--fill"{
						return nil
					}
					return []prompt.Suggest{
						{Text: "--fill", Description: "fill date from previous or default time"},
						{Text: "HHMM", Description: "time"},
					}
				}
				return prompt.FilterContains(paceData.IssueSuggestions, d.GetWordBeforeCursor(), true)
			} else if strings.Compare("issues", fields[0]) == 0 {
				s := []prompt.Suggest{
					{Text: "Not Implemented", Description: "Not Implemented"},
				}
				return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
			} else if strings.Compare("refresh", fields[0]) == 0 {
				return nil
			} else if strings.Compare("settings", fields[0]) == 0 {
				s := []prompt.Suggest{
					{Text: "Not Implemented", Description: "Not Implemented"},
				}
				return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
			}
		}
		s := []prompt.Suggest{
			{Text: "time", Description: "time in period"},
			{Text: "clock", Description: "clock to track time"},
			{Text: "log", Description: "manually log time"},
			{Text: "issues", Description: "search issues"},
			{Text: "refresh", Description: "refresh issues"},
			{Text: "settings", Description: "manage pace settings"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
}
