package main

import (
	"github.com/c-bata/go-prompt"
)

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
		paceCompletor(config, paceData),
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
