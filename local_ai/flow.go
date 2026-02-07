package local_ai

import (
	"context"
	"obsidian-ai-planner/planner"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type ModelInfo struct {
	GenKit         *genkit.Genkit
	Model          ai.Model
	PlannerService *planner.PlannerService
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PlannerInput struct {
	UserPrompt string    `json:"userPrompt"`
	History    []Message `json:"history"`
}

func (m *ModelInfo) Chat(ctx context.Context, input PlannerInput) (string, error) {
	pContext, err := m.PlannerService.FetchContext(ctx)
	if err != nil {
		return "", err
	}

	systemPrompt := planner.BuildSystemPrompt(pContext)

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
	pContext, err := m.PlannerService.FetchContext(ctx)
	if err != nil {
		return "", err
	}

	systemPrompt := planner.BuildGeneratePlanPrompt(pContext)

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
	systemPrompt := planner.BuildCondensePrompt()

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
