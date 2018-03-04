package main

import (
	"io/ioutil"
	"github.com/pkg/errors"
	"os"
	"github.com/BurntSushi/toml"
	"fmt"
	"time"
	"bytes"
)

//TODO FUTURE [will-do]: Add functionality for default config values.
//TODO FUTURE [maybe]: Add ability to round to nearest x duration.

const DEFAULT_CONFIG_FILE_LOCATION = "/.config/pace/"
const DEFAULT_CONFIG_FILE_NAME = "config.toml"

const DEFAULT_CONFIG_FILE_CONTENT = `# Pace Configuration File

# Jira Instance Information
JiraInstanceUrl = "https://mycompany.atlassian.net"
JiraUsername = "john.doe"
JiraPassword = "password"

# Appearance/Behavior
Prompt = ">>> "

# Prebuilt Query Defaults
ProjectName = "MYPROJECT"
QueryUsername = "john.doe"

# JIRA User Defined JQL
SuggestionIssueQuery = "project = MYPROJECT AND assignee = john.doe AND resolution = Unresolved ORDER BY updated DESC"
DayWorklogIssueQuery = "project = MYPROJECT AND assignee = john.doe AND worklogDate >= startOfDay() ORDER BY created DESC"
WeekWorklogIssueQuery = "project = MYPROJECT AND assignee = john.doe AND worklogDate >= startOfWeek() ORDER BY created DESC"
MonthWorklogIssueQuery = "project = MYPROJECT AND assignee = john.doe AND worklogDate >= startOfMonth() ORDER BY created DESC"

# Fill Option Behavior
LogFillStartTime = "0800"`

type Config struct {
	JiraInstanceUrl string
	JiraUsername    string
	JiraPassword    string

	Prompt string

	ProjectName   string
	QueryUsername string

	SuggestionIssueQuery   string
	DayWorklogIssueQuery   string
	WeekWorklogIssueQuery  string
	MonthWorklogIssueQuery string

	DefaultProject string //todo implement
	ProjectShortcuts []ProjectShortcut //todo implement

	ClockStartTime time.Time

	LogFillStartTime string //format 1625 or 0425; HHMM
}

type ProjectShortcut struct {
	Name string
	Shortcut string
}

func NewConfig() *Config {
	return &Config{}
}

//todo use file path joiner to remove os specific file path separator; add function to get dir; add function to get file

func (config *Config) InitializeConfig() error {
	if _, err := os.Stat(os.Getenv("HOME") + DEFAULT_CONFIG_FILE_LOCATION + DEFAULT_CONFIG_FILE_NAME); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(os.Getenv("HOME")+DEFAULT_CONFIG_FILE_LOCATION, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return errors.New("failed to create pace config directory")
			}
			file, err := os.Create(os.Getenv("HOME") + DEFAULT_CONFIG_FILE_LOCATION + DEFAULT_CONFIG_FILE_NAME)
			if err != nil {
				return errors.New("failed to create config file")
			}
			defer file.Close()

			_, err = file.WriteString(DEFAULT_CONFIG_FILE_CONTENT)
			if err != nil {
				return errors.New("failed to write config file")
			}
			err = file.Sync()
			if err != nil {
				return errors.New("failed to sync config file")
			}
			return nil
		} else {
			return errors.New("failed to check for configuration file")
		}
	}

	cbuf, err := ioutil.ReadFile(os.Getenv("HOME") + DEFAULT_CONFIG_FILE_LOCATION + DEFAULT_CONFIG_FILE_NAME)
	if err != nil {
		return errors.New("failed to read configuration file")
	}

	if _, err := toml.Decode(string(cbuf), config); err != nil {
		return errors.New("failed to parse configuration file")
	}

	valid, err := validateConfig(config)
	if err != nil || !valid {
		return errors.New("invalid configuration file")
	}

	return nil
}

func (config *Config) ApplyDefaults() {
	return
}

func validateConfig(config *Config) (bool, error) {
	//todo implement validation
	return true, nil
}

func (config *Config) PersistConfig() error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return err
	}
	file, err := os.Create(os.Getenv("HOME") + DEFAULT_CONFIG_FILE_LOCATION + DEFAULT_CONFIG_FILE_NAME)
	if err != nil {
		return errors.New("failed to create config file")
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return errors.New("failed to write config file")
	}
	err = file.Sync()
	if err != nil {
		return errors.New("failed to sync config file")
	}
	return nil
}

func (config *Config) ConfigUpdated() error {
	return config.PersistConfig()
}
