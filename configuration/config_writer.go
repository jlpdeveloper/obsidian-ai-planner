package configuration

import (
	"encoding/json"
	"os"
)

var configLocation string = "./config.json"

type Config struct {
	CalendarUrl string `json:"calendar_url"`
	JiraEmail   string `json:"jira_email"`
	JiraToken   string `json:"jira_token"`
}

func (c *Config) Write() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configLocation, data, 0644)
}

func (c *Config) LoadFromFile() error {
	data, err := os.ReadFile(configLocation)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}
