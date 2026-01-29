package local_ai

import (
	"context"
	"fmt"
	"obsidian-ai-planner/calendar"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type ModelInfo struct {
	GenKit   *genkit.Genkit
	Model    ai.Model
	Calendar *calendar.GoogleCalendarIntegration
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PlannerInput struct {
	UserPrompt string    `json:"userPrompt"`
	History    []Message `json:"history"`
}

type InternalPlannerContext struct {
	WeeklyGoals  string           `json:"weeklyGoals"`
	Calendar     []calendar.Event `json:"calendar"`
	JiraTickets  []string         `json:"jiraTickets"`
	CurrentTasks []string         `json:"currentTasks"`
}

func (m *ModelInfo) fetchContext(ctx context.Context) (*InternalPlannerContext, error) {
	// TODO: Pull from Obsidian
	weeklyGoals := "Plan for project unicorn, Review roadmap, Improve test coverage 10%"
	// TODO: Pull from Daily Note
	currentTasks := []string{}
	// TODO: Pull from Jira
	jiraTickets := []string{"Jira-123: Update db", "Jira-456: Fix bug on backend"}

	var calendarEvents []calendar.Event
	if m.Calendar != nil {
		// Today's date helper
		now := time.Now()
		year, month, day := now.Date()
		today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

		events, err := m.Calendar.GetCalendarEvents(today)
		if err != nil {
			return nil, err
		}
		calendarEvents = events
	}

	return &InternalPlannerContext{
		WeeklyGoals:  weeklyGoals,
		Calendar:     calendarEvents,
		JiraTickets:  jiraTickets,
		CurrentTasks: currentTasks,
	}, nil
}

func (m *ModelInfo) Chat(ctx context.Context, input PlannerInput) (string, error) {
	pContext, err := m.fetchContext(ctx)
	if err != nil {
		return "", err
	}

	systemPrompt := fmt.Sprintf(`
	You are an opinionated personal planning analyst assisting a software engineer with their day.

You have no access to outside tools, memory, or systems beyond the data provided in this message.
You must reason only from the supplied goals, calendar events, Jira tickets, and tasks.

Your primary responsibility is to assess overcommitment.

Overcommitment means:
Planned work + context switching + cognitive overhead exceeds realistic daily capacity.

You are allowed to:
- Estimate task effort when no duration is provided
- Make reasonable assumptions about a standard workday
- Apply an overhead factor for meetings and context switching (state your assumptions)
- Be uncertain but still decisive

You are encouraged to:
- Call out when the plan does not mathematically fit in the day
- Point out hidden overload, fragmentation, or unrealistic sequencing
- Push back on priorities when trade-offs are required

You should:
- Look for alignment between weekly goals and Jira tickets
- Treat calendar events as hard constraints
- Treat tasks and tickets as flexible unless stated otherwise

You should NOT:
- Attempt to optimize or rewrite the full day
- Generate a complete daily note
- Store or assume long-term user behavior

Tone and format:
- Be concise, conversational, and direct
- Prefer clear assertions over vague suggestions
- If the day appears overcommitted, say so plainly

Inputs:
Current Weekly Goals: %s
Calendar Events: %v
Jira Tickets: %v
Current Tasks: %v

Respond by discussing the plan, highlighting risks or mismatches, or answering the user's question.

`, pContext.WeeklyGoals, pContext.Calendar, pContext.JiraTickets, pContext.CurrentTasks)

	var messages []*ai.Message
	messages = append(messages, ai.NewSystemMessage(ai.NewTextPart(systemPrompt)))

	for _, msg := range input.History {
		if msg.Role == "user" {
			messages = append(messages, ai.NewUserMessage(ai.NewTextPart(msg.Content)))
		} else if msg.Role == "model" || msg.Role == "bot" || msg.Role == "assistant" {
			messages = append(messages, ai.NewModelMessage(ai.NewTextPart(msg.Content)))
		}
	}

	messages = append(messages, ai.NewUserMessage(ai.NewTextPart(input.UserPrompt)))

	resp, err := genkit.Generate(ctx, m.GenKit,
		ai.WithModel(m.Model),
		ai.WithMessages(messages...),
	)
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}

func (m *ModelInfo) GeneratePlan(ctx context.Context, input PlannerInput) (string, error) {
	pContext, err := m.fetchContext(ctx)
	if err != nil {
		return "", err
	}

	systemPrompt := fmt.Sprintf(`
You are a personal AI planner. Your goal is to help a software engineer plan their day by generating structured updates for their daily note.
Current Weekly Goals: %s
Calendar Events: %v
Jira Tickets: %v
Current Tasks: %v

Please generate the content for the 'Goals', 'Meetings', and 'Bonus Items' sections. 
Be specific and professional. Use Markdown format.
`, pContext.WeeklyGoals, pContext.Calendar, pContext.JiraTickets, pContext.CurrentTasks)

	var messages []*ai.Message
	messages = append(messages, ai.NewSystemMessage(ai.NewTextPart(systemPrompt)))

	for _, msg := range input.History {
		if msg.Role == "user" {
			messages = append(messages, ai.NewUserMessage(ai.NewTextPart(msg.Content)))
		} else if msg.Role == "model" || msg.Role == "bot" || msg.Role == "assistant" {
			messages = append(messages, ai.NewModelMessage(ai.NewTextPart(msg.Content)))
		}
	}

	messages = append(messages, ai.NewUserMessage(ai.NewTextPart(input.UserPrompt)))

	resp, err := genkit.Generate(ctx, m.GenKit,
		ai.WithModel(m.Model),
		ai.WithMessages(messages...),
	)
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}

func (m *ModelInfo) Condense(ctx context.Context, history []Message) (string, error) {
	systemPrompt := "You are a helpful assistant. Summarize the following conversation history concisely, preserving all key decisions, tasks, and context. This summary will be used as the starting point for a new conversation session."

	var messages []*ai.Message
	messages = append(messages, ai.NewSystemMessage(ai.NewTextPart(systemPrompt)))

	for _, msg := range history {
		if msg.Role == "user" {
			messages = append(messages, ai.NewUserMessage(ai.NewTextPart(msg.Content)))
		} else if msg.Role == "model" || msg.Role == "bot" || msg.Role == "assistant" {
			messages = append(messages, ai.NewModelMessage(ai.NewTextPart(msg.Content)))
		}
	}

	resp, err := genkit.Generate(ctx, m.GenKit,
		ai.WithModel(m.Model),
		ai.WithMessages(messages...),
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
