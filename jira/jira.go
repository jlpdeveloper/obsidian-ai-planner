package jira

import (
	"context"
	"fmt"
	"obsidian-ai-planner/configuration"
	"regexp"

	"github.com/andygrunwald/go-jira/v2/cloud"
)

type JiraIntegration struct {
	client *cloud.Client
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

	// In v2/cloud, we use BasicAuthTransport or similar, but the example shows a simpler approach
	// or assumes the client is created with a transport.
	// Actually, the example uses cloud.NewClient.

	tp := cloud.BasicAuthTransport{
		Username: cfg.JiraEmail,
		APIToken: cfg.JiraToken,
	}

	client, err := cloud.NewClient(cfg.JiraBaseUrl, tp.Client())
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
	options := &cloud.SearchOptionsV2{
		Fields: []string{"key", "summary"},
	}
	// Using context.Background() as a default
	issues, _, err := j.client.Issue.SearchV2JQL(context.Background(), jql, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search jira issues: %w", err)
	}

	return j.processIssues(issues), nil
}

func (j *JiraIntegration) processIssues(issues []cloud.Issue) []string {
	var tickets []string
	accountRegex := regexp.MustCompile(`\b\d{6,17}\b`)
	for _, issue := range issues {
		summary := accountRegex.ReplaceAllString(issue.Fields.Summary, "[ACCOUNT_NUMBER]")
		tickets = append(tickets, fmt.Sprintf("%s: %s", issue.Key, summary))
	}

	return tickets
}
