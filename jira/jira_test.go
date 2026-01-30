package jira

import (
	"obsidian-ai-planner/configuration"
	"testing"

	"github.com/andygrunwald/go-jira/v2/cloud"
)

func TestProcessIssues(t *testing.T) {
	ji := &JiraIntegration{}

	issues := []cloud.Issue{
		{
			Key: "PROJ-1",
			Fields: &cloud.IssueFields{
				Summary: "Normal ticket",
			},
		},
		{
			Key: "PROJ-2",
			Fields: &cloud.IssueFields{
				Summary: "Ticket with 123456789 account",
			},
		},
		{
			Key: "PROJ-3",
			Fields: &cloud.IssueFields{
				Summary: "Ticket with short 12345 digits",
			},
		},
		{
			Key: "PROJ-4",
			Fields: &cloud.IssueFields{
				Summary: "Ticket with long 123456789012345678 digits",
			},
		},
		{
			Key: "PROJ-5",
			Fields: &cloud.IssueFields{
				Summary: "Ticket with 6 digits 123456",
			},
		},
		{
			Key: "PROJ-6",
			Fields: &cloud.IssueFields{
				Summary: "Ticket with 17 digits 12345678901234567",
			},
		},
	}

	result := ji.processIssues(issues)

	expected := []string{
		"PROJ-1: Normal ticket",
		"PROJ-2: Ticket with [ACCOUNT_NUMBER] account",
		"PROJ-3: Ticket with short 12345 digits",
		"PROJ-4: Ticket with long 123456789012345678 digits",
		"PROJ-5: Ticket with 6 digits [ACCOUNT_NUMBER]",
		"PROJ-6: Ticket with 17 digits [ACCOUNT_NUMBER]",
	}

	if len(result) != len(expected) {
		t.Errorf("Expected %d tickets, got %d", len(expected), len(result))
	}

	for i, ticket := range result {
		if i < len(expected) && ticket != expected[i] {
			t.Errorf("Expected ticket %d to be %q, got %q", i, expected[i], ticket)
		}
	}
}

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
