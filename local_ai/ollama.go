package local_ai

import (
	"context"
	"obsidian-ai-planner/calendar"
	"obsidian-ai-planner/jira"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
)

func NewOllamaModel(ctx context.Context) (*ModelInfo, error) {
	ollamaPlugin := &ollama.Ollama{
		ServerAddress: "http://127.0.0.1:11434",
		Timeout:       60, // Optional field, adjust accordingly
	}
	g := genkit.Init(ctx, genkit.WithPlugins(ollamaPlugin))

	model := ollamaPlugin.DefineModel(g,
		ollama.ModelDefinition{
			Name: "llama3.1",
			Type: "chat", // "chat" or "generate"
		},
		&ai.ModelOptions{
			Supports: &ai.ModelSupports{
				Multiturn:  true,
				SystemRole: true,
				Tools:      false,
				Media:      false,
			},
		},
	)

	cal, _ := calendar.New(ctx)
	ji, _ := jira.New()

	return &ModelInfo{
		Model:    model,
		GenKit:   g,
		Calendar: cal,
		Jira:     ji,
	}, nil
}
