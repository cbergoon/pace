<h1 align="center">Pace - A CLI Work Log Manager for Jira</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/pace"><img src="https://goreportcard.com/badge/github.com/cbergoon/pace?1=2" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/pace"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

[![asciicast](https://asciinema.org/a/204973.png)](https://asciinema.org/a/204973?t=8)

Pace provides an easy way to manage and track worklog entries in Jira. 

### Install

```
$ go get github.com/cbergoon/pace
```

Run Pace with ```$ pace``` to generate the default configuration file in ```~/.config/pace/config.toml```

Update the configuration file with details from your Jira instance and Jira Projects. 

### Usage

Run Pace with ```$ pace```. Your Jira instance will be queried for open issues in the projects specified in the configuration. Worklogs for each issue will be retrieved and an ordered time log will be constructed. 

#### Time

Time allows you to view worklogs for a given time period. 

#### Log

Log allows you to create worklogs for a given issue. 

The general format of the command is: 
``` 
>>> log ISSUE-100 1h30m 20181005 0900
``` 

Issues will be searched by the issue key and a list of options to select from will be displayed below. The duration parameter should be in the form of \<H\>h\<M\>m (e.g. 1h30m for 1 hour and 30 minutes). The date and time should follow the format suggested by the prompt; "YYYYMMDD" and "HHMM" (e.g. 20181005 for 5 October 2018). If the ```FillOptionEnabled``` flag is true in the cofiguration file and alternative option will be available for date and time: ```--fill```. This option will set the starting time of the worklog entry to be created to the computed end time of the last work log. 

#### Refresh

The refresh option will refresh all of the issues and worklog information retrieved at startup. If the ```FillOptionEnabled``` setting is set to true the data will be refreshed after each worklog entry is created. 


