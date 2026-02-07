package planner

import (
	"context"
	"fmt"
	"obsidian-ai-planner/calendar"
	"obsidian-ai-planner/jira"
	"obsidian-ai-planner/obsidian"
	"time"
)

type PlannerContext struct {
	WeeklyGoals  string           `json:"weeklyGoals"`
	Calendar     []calendar.Event `json:"calendar"`
	JiraTickets  []string         `json:"jiraTickets"`
	CurrentTasks []string         `json:"currentTasks"`
}

type PlannerService struct {
	Calendar      *calendar.GoogleCalendarIntegration
	Jira          *jira.JiraIntegration
	ObsidianVault string
}

func (s *PlannerService) FetchContext(ctx context.Context) (*PlannerContext, error) {
	now := time.Now()
	weeklyGoals := ""
	currentTasks := []string{}

	if s.ObsidianVault != "" {
		// Fetch Weekly Goals
		weeklyPath := obsidian.GetWeeklyNotePath(s.ObsidianVault, now)
		goals, err := obsidian.ReadSection(weeklyPath, "Weekly Goals")
		if err == nil {
			weeklyGoals = goals
		} else {
			fmt.Printf("Warning: failed to fetch weekly goals: %v\n", err)
		}

		// Fetch Daily Tasks
		dailyPath := obsidian.GetDailyNotePath(s.ObsidianVault, now)
		tasks, err := obsidian.ReadSection(dailyPath, "Tasks")
		if err == nil {
			currentTasks = []string{tasks}
		} else {
			fmt.Printf("Warning: failed to fetch daily tasks: %v\n", err)
		}
	}

	jiraTickets := []string{"Jira-123: Update db (Mock)", "Jira-456: Fix bug on backend (Mock)"}
	if s.Jira != nil {
		tickets, err := s.Jira.GetAssignedTickets(ctx)
		if err == nil {
			jiraTickets = tickets
		} else {
			fmt.Printf("Warning: failed to fetch Jira tickets: %v\n", err)
		}
	}

	var calendarEvents []calendar.Event
	if s.Calendar != nil {
		// Today's date helper
		now := time.Now()
		year, month, day := now.Date()
		today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

		events, err := s.Calendar.GetCalendarEvents(ctx, today)
		if err != nil {
			return nil, err
		}
		calendarEvents = events
	}

	return &PlannerContext{
		WeeklyGoals:  weeklyGoals,
		Calendar:     calendarEvents,
		JiraTickets:  jiraTickets,
		CurrentTasks: currentTasks,
	}, nil
}
