package jira

import (
	"obsidian-ai-planner/configuration"
	"testing"
)

func TestJiraIntegrationBasic(t *testing.T) {
	// Just test that New() fails gracefully if config is missing or incomplete
	cfg := &configuration.Config{}
	err := cfg.LoadFromFile()

	if err != nil {
		ji, err := New()
		if err == nil {
			t.Error("Expected error when config is missing, got nil")
		}
		if ji != nil {
			t.Error("Expected nil JiraIntegration when config is missing")
		}
	} else {
		// If config exists, check if it's incomplete
		if cfg.JiraBaseUrl == "" || cfg.JiraEmail == "" || cfg.JiraToken == "" {
			ji, err := New()
			if err == nil {
				t.Error("Expected error when config is incomplete, got nil")
			}
			if ji != nil {
				t.Error("Expected nil JiraIntegration when config is incomplete")
			}
		}
	}
}
