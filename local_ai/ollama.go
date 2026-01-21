package local_ai

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
)

func NewOllamaModel(ctx context.Context) *ModelInfo {
	ollamaPlugin := &ollama.Ollama{
		ServerAddress: "http://127.0.0.1:11434",
		Timeout:       60, // Optional field, adjust accordingly
	}
	g := genkit.Init(ctx, genkit.WithPlugins(ollamaPlugin))

	model := ollamaPlugin.DefineModel(g,
		ollama.ModelDefinition{
			Name: "gemma3",
			Type: "generate", // "chat" or "generate"
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

	return &ModelInfo{
		Model:  model,
		GenKit: g,
	}
}
