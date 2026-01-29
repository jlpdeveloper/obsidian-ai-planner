package jira

import (
	"fmt"
	"obsidian-ai-planner/configuration"

	"github.com/andygrunwald/go-jira"
)

type JiraIntegration struct {
	client *jira.Client
}

func New() (*JiraIntegration, error) {
	cfg := &configuration.Config{}
	err := cfg.LoadFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.JiraBaseUrl == "" || cfg.JiraEmail == "" || cfg.JiraToken == "" {
		return nil, fmt.Errorf("jira configuration is incomplete")
	}

	tp := jira.BasicAuthTransport{
		Username: cfg.JiraEmail,
		Password: cfg.JiraToken,
	}

	client, err := jira.NewClient(tp.Client(), cfg.JiraBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create jira client: %w", err)
	}

	return &JiraIntegration{
		client: client,
	}, nil
}

func (j *JiraIntegration) GetAssignedTickets() ([]string, error) {
	// JQL to get tickets assigned to the current user that are not closed
	jql := "assignee = currentUser() AND statusCategory != Done"

	issues, _, err := j.client.Issue.Search(jql, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search jira issues: %w", err)
	}

	var tickets []string
	for _, issue := range issues {
		tickets = append(tickets, fmt.Sprintf("%s: %s", issue.Key, issue.Fields.Summary))
	}

	return tickets, nil
}
