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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PlannerInput struct {
	WeeklyGoals  string           `json:"weeklyGoals"`
	Calendar     []calendar.Event `json:"calendar"`
	JiraTickets  []string         `json:"jiraTickets"`
	CurrentTasks []string         `json:"currentTasks"`
	UserPrompt   string           `json:"userPrompt"`
	History      []Message        `json:"history"`
}

func (m *ModelInfo) Chat(ctx context.Context, input PlannerInput) (string, error) {
	systemPrompt := fmt.Sprintf(`
You are a personal AI planner assistant. You are having a conversation with a software engineer about their day. 
You have no access to outside tools and no awareness of exterior systems except with the data provided. 
You should examine the weekly goals and infer if any Jira tickets align with them. 
Current Weekly Goals: %s
Calendar Events: %v
Jira Tickets: %v
Current Tasks: %v

Please respond concisely and conversationally. Do not generate the full daily note yet. Just discuss the plan or answer questions.
`, input.WeeklyGoals, input.Calendar, input.JiraTickets, input.CurrentTasks)

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
	systemPrompt := fmt.Sprintf(`
You are a personal AI planner. Your goal is to help a software engineer plan their day by generating structured updates for their daily note.
Current Weekly Goals: %s
Calendar Events: %v
Jira Tickets: %v
Current Tasks: %v

Please generate the content for the 'Goals', 'Meetings', and 'Bonus Items' sections. 
Be specific and professional. Use Markdown format.
`, input.WeeklyGoals, input.Calendar, input.JiraTickets, input.CurrentTasks)

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
