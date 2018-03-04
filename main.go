package main

import (
	"github.com/c-bata/go-prompt"
)

//TODO FUTURE [maybe]: Add 'State' that is persisted to a file. This would cache issues
//and other data to preserve state and reduce load time.
//TODO FUTURE [will-do]: Add support for Pomodoro timer to clock. Config for 3 levels of time boxing.
//Alerts to user (Mac OS notifications).
//TODO FUTURE [will-do]: Support multiple projects.

func main() {
	config := NewConfig()
	paceData := NewPaceData()

	err := config.InitializeConfig()
	if err != nil {
		panic(err)
	}

	jiraClient, err := configureNewJiraClient(config.JiraInstanceUrl, config.JiraUsername, config.JiraPassword)
	if err != nil {
		panic(err)
	}

	err = paceData.loadIssueQueues(jiraClient, config)
	if err != nil {
		panic(err)
	}

	p := prompt.New(
		paceExecutor(jiraClient, config, paceData),
		paceCompletor(paceData),
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("pace-prompt"),

		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.DarkBlue),
		prompt.OptionSuggestionTextColor(prompt.White),

		prompt.OptionSelectedDescriptionBGColor(prompt.DarkBlue),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkBlue),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
	)
	p.Run()
}
