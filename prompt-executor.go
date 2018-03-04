package main

import (
	"strings"
	"fmt"
	"time"
	"strconv"
	"github.com/andygrunwald/go-jira"
)

//TODO FUTURE [will-do]: Project shortcuts ie 'WEBPORTAL-' -> 'WP'.
//TODO FUTURE [will-do]: Default/favorite project.
//TODO FUTURE [maybe]: On times ability to omit the leading '0'/'20' i.e. '0330' == '330', '2018' == '18'.

func paceExecutor(jiraClient *jira.Client, config *Config, paceData *PaceData) func(in string) {
	return func(in string) {
		fields := strings.Fields(in)
		if strings.Compare(fields[0], "time") == 0 {
			if strings.Compare(fields[1], "today") == 0 {
				printTimeStats(generateTimeStats(paceData.DayIssues, beginningOfDay, config.QueryUsername))
			} else if strings.Compare(fields[1], "week") == 0 {
				printTimeStats(generateTimeStats(paceData.WeekIssues, beginningOfWeek, config.QueryUsername))
			} else if strings.Compare(fields[1], "month") == 0 {
				printTimeStats(generateTimeStats(paceData.MonthIssues, beginningOfMonth, config.QueryUsername))
			}
		} else if strings.Compare(fields[0], "clock") == 0 {
			if config.ClockStartTime.Unix() < 0 {
				fmt.Println("The clock must be started before it is stopped.")
				return
			}
			if len(fields) == 1 {
				fmt.Println("Clock started " + time.Now().Sub(config.ClockStartTime).Truncate(time.Minute).String() + " ago at " + config.ClockStartTime.Format(DISPLAY_DATE_FORMAT))
				return
			}
			if strings.Compare(fields[1], "start") == 0 {
				config.ClockStartTime = time.Now()
				err := config.ConfigUpdated()
				if err != nil {
					fmt.Println("Something went wrong, ensure the clock is started.")
				}
			} else if strings.Compare(fields[1], "stop") == 0 {
				if len(fields) < 3 {
					fmt.Println("An issue must be selected in order to log time.")
					return
				}

				comment := "Logged with Care by Pace.\n"
				if len(fields) >= 4 {
					comment += strings.Join(fields[3:], " ")
				}

				elapsedTime := time.Now().Sub(config.ClockStartTime).Truncate(time.Minute)
				err := addWorklog(jiraClient, fields[2], WorklogRequest{
					TimeSpentSeconds: fmt.Sprint(elapsedTime.Seconds()),
					Started:          config.ClockStartTime.Format("2006-01-02T15:04:05.000-0700"),
					Comment:          comment,
				})
				if err != nil {
					fmt.Println("Something went wrong, ensure the worklog was created.")
				}

				fmt.Println("Logged " + elapsedTime.String() + " to issue " + fields[2])

				config.ClockStartTime = time.Now()
				err = config.ConfigUpdated()
				if err != nil {
					fmt.Println("Something went wrong, ensure the clock is started.")
				}
			}
		} else if strings.Compare(fields[0], "log") == 0 {
			if len(fields) < 3 {
				fmt.Println("An issue and duration must be specified to log time.")
				return
			}

			var hour int64 = 0
			var minute int64 = 0
			if strings.Contains(fields[2], "h") {
				tmp := strings.Replace(fields[2], "m", "", 1)
				parts := strings.Split(tmp, "h")
				hour, _ = strconv.ParseInt(parts[0], 10, 32)
				minute, _ = strconv.ParseInt(parts[1], 10, 32)
			} else {
				tmp := strings.Replace(fields[2], "m", "", 1)
				minute, _ = strconv.ParseInt(tmp, 10, 32)
			}
			d := time.Duration((time.Minute * time.Duration(minute)) + (time.Duration(time.Hour * time.Duration(hour))))

			comment := "Logged with Care by Pace.\n"
			worklogRequest := WorklogRequest{
				TimeSpentSeconds: fmt.Sprint(d.Seconds()),
				Comment:          comment,
			}

			if len(fields) == 3 {
				worklogRequest.Started = time.Now().Add(-1 * d).Truncate(time.Second).Format("2006-01-02T15:04:05.000-0700")
				err := addWorklog(jiraClient, fields[1], worklogRequest)
				if err != nil {
					fmt.Println("Something went wrong, ensure the worklog was created.")
					return
				}
				return
			}

			var dateFill time.Time
			if len(fields) >= 4 {
				if fields[3] == "--fill" {
					//log from default or previous
					dateFill = paceData.findLastWorklogEndTimeForUserForDate(time.Now(), config.QueryUsername)
					if dateFill.Unix() < 0 {
						if len(config.LogFillStartTime) <= 0 {
							fmt.Println("Configuration does not contain default log fill start time.")
							return
						}
						dateFill = time.Now()
						timeOnlyDateFill, err := time.ParseInLocation("1504", fmt.Sprint(config.LogFillStartTime), time.Local)
						if err != nil {
							fmt.Println("Invalid date/time format.")
							return
						}
						dateFill = time.Date(dateFill.Year(), dateFill.Month(), dateFill.Day(), timeOnlyDateFill.Hour(), timeOnlyDateFill.Minute(), 0, 0, time.Local)
					}
				}else{
					var err error
					dateFill, err = time.ParseInLocation("20060102", fmt.Sprint(fields[3]), time.Local)
					if err != nil {
						fmt.Println("Invalid date/time format.")
						return
					}
				}
			}
			timeFill := dateFill
			if fields[3] != "--fill" {
				if len(fields) >= 5 {
					if fields[4] == "--fill" {
						//log from specified date w/ previous time
						timeFill = paceData.findLastWorklogEndTimeForUserForDate(dateFill, config.QueryUsername)
						if timeFill.Unix() < 0 {
							if len(config.LogFillStartTime) <= 0 {
								fmt.Println("Configuration does not contain default log fill start time.")
								return
							}
							timeFill = dateFill
							timeOnlyTimeFill, err := time.ParseInLocation("1504", fmt.Sprint(config.LogFillStartTime), time.Local)
							if err != nil {
								fmt.Println("Invalid date/time format.")
								return
							}
							timeFill = time.Date(dateFill.Year(), dateFill.Month(), dateFill.Day(), timeOnlyTimeFill.Hour(), timeOnlyTimeFill.Minute(), 0, 0, time.Local)
						}
					} else {
						var err error
						timeFill, err = time.ParseInLocation("20060102 1504", fmt.Sprint(fields[3]+" "+fields[4]), time.Local)
						if err != nil {
							fmt.Println("Invalid date/time format.")
							return
						}
					}
				}
			}

			started := timeFill
			err := addWorklog(jiraClient, fields[1], WorklogRequest{
				Started:          started.Truncate(time.Second).Format("2006-01-02T15:04:05.000-0700"),
				TimeSpentSeconds: fmt.Sprint(d.Seconds()),
				Comment:          comment,
			})
			if err != nil {
				fmt.Println("Something went wrong, ensure the worklog was created.")
				return
			}

			paceData.loadIssueQueues(jiraClient, config)

			return
		} else if strings.Compare(fields[0], "issues") == 0 {
			fmt.Println("Not Implemented")
		} else if strings.Compare(fields[0], "refresh") == 0 {
			paceData.loadIssueQueues(jiraClient, config)
		} else if strings.Compare(fields[0], "settings") == 0 {
			fmt.Println("Not Implemented")
		}
	}
}
