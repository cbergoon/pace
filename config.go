package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

//TODO (cbergoon): Add ability to round to nearest x duration.

const defaultConfigFileLocation = "/.config/pace/"
const defaultConfigFileName = "config.toml"

const defaultConfigFileContent = `# Pace Configuration File

# Jira Instance Information

# Jira URL of Atlassian instance to operate on. Must begin with "https://".
JiraInstanceUrl = "https://mycompany.atlassian.net"

# Jira username of to operate as. User must have ability to read and update issues. 
JiraUsername = "john.doe"

# Password of above user account. 
JiraPassword = "password"

##############################################################################################################

# Appearance/Behavior

# Prompt prefix display text. 
Prompt = ">>> "

##############################################################################################################

# Queries

# Jira username to be used to in JQL queries to retrieve issues. Only time for specified user will be 
# included in time log. 
QueryUsername = "john.doe"
SuggestionProjects = ["PR1", "PR2"]

##############################################################################################################

# Fill Option Behavior

# FillOptionEnabled allows the log to be filled without entering date/time with each entry. Depending on 
# network speeds, Jira instance speeds, number of projects, and project size this may cause performance 
# issues and can be disabled with the option below.
FillOptionEnabled = false

# LogFillStartTime option sets the time that should be considered the start of the work day. The fill option
# will use this time when entering creating the first entry of a day. Format: #### (e.g. 0800 for 8:00 AM, 
# 1630 for 4:30 PM)
LogFillStartTime = "0800"`

type Config struct {
	JiraInstanceUrl string
	JiraUsername    string
	JiraPassword    string

	Prompt string

	QueryUsername      string
	SuggestionProjects []string

	ClockStartTime    time.Time
	LogFillStartTime  string //format 1625 or 0425; HHMM
	FillOptionEnabled bool
}

type ProjectShortcut struct {
	Name     string
	Shortcut string
}

func NewConfig() *Config {
	return &Config{}
}

func (config *Config) InitializeConfig() error {
	//TODO (cbergoon): use file path joiner to remove os specific file path separator; add function to get dir; add function to get file
	if _, err := os.Stat(os.Getenv("HOME") + defaultConfigFileLocation + defaultConfigFileName); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(os.Getenv("HOME")+defaultConfigFileLocation, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return errors.New("failed to create pace config directory")
			}
			file, err := os.Create(os.Getenv("HOME") + defaultConfigFileLocation + defaultConfigFileName)
			if err != nil {
				return errors.New("failed to create config file")
			}
			defer file.Close()

			_, err = file.WriteString(defaultConfigFileContent)
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

	cbuf, err := ioutil.ReadFile(os.Getenv("HOME") + defaultConfigFileLocation + defaultConfigFileName)
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
	//TODO (cbergoon): Implement validation.
	return true, nil
}

func (config *Config) PersistConfig() error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return err
	}
	file, err := os.Create(os.Getenv("HOME") + defaultConfigFileLocation + defaultConfigFileName)
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
