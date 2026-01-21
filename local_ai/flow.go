package local_ai

import (
	"context"
	"fmt"
	"obsidian-ai-planner/calendar"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type ModelInfo struct {
	GenKit *genkit.Genkit
	Model  ai.Model
}

type PlannerInput struct {
	WeeklyGoals  string           `json:"weeklyGoals"`
	Calendar     []calendar.Event `json:"calendar"`
	JiraTickets  []string         `json:"jiraTickets"`
	CurrentTasks []string         `json:"currentTasks"`
	UserPrompt   string           `json:"userPrompt"`
}

func (m *ModelInfo) GeneratePlan(ctx context.Context, input PlannerInput) (string, error) {
	prompt := fmt.Sprintf(`
You are a personal AI planner. Your goal is to help a software engineer plan their day.
Current Weekly Goals:
%s

Calendar Events:
%v

Jira Tickets:
%v

Current Tasks in Daily Note:
%v

User says: %s

Please propose a plan for the day, specifically updating the 'Goals', 'Meetings', and 'Bonus Items' sections as needed.
`, input.WeeklyGoals, input.Calendar, input.JiraTickets, input.CurrentTasks, input.UserPrompt)

	resp, err := genkit.Generate(ctx, m.GenKit,
		ai.WithModel(m.Model),
		ai.WithPrompt(prompt),
	)
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}

func DefinePlannerFlow(m *ModelInfo) {
	genkit.DefineFlow(m.GenKit, "plannerFlow", func(ctx context.Context, input PlannerInput) (string, error) {
		return m.GeneratePlan(ctx, input)
	})
}
