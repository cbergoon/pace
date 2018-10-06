<h1 align="center">Pace - A CLI Work Log Manager for Jira</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/pace"><img src="https://goreportcard.com/badge/github.com/cbergoon/pace?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/pace"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

[![asciicast](https://asciinema.org/a/204973.png)](https://asciinema.org/a/204973?t=8)

Pace provides an easy way to manage and track worklog entries in Jira. 

#### Install

```
$ go get github.com/cbergoon/pace
```

Run Pace with ```$ pace``` to generate the default configuration file in ```~/.config/pace/config.toml```

Update the configuration file with details from your Jira instance and Jira Projects. 

#### Usage

Run Pace with ```$ pace``` 

##### Time

Time allows you to view worklogs for a given time period. 

##### Log

Log allows you to create worklogs for a given issue. 

The format of the command is: 
``` 
>>> log ISSUE-100 1h30m 20181005 0900
``` 


