package domain

import "context"

// ContextInitializer defines the interface for initializing AI context in a project.
type ContextInitializer interface {
	// Initialize sets up the context for the given tool and technology.
	Initialize(ctx context.Context, tool, technology string) error
}

// PromptGenerator defines the interface for generating prompts.
type PromptGenerator interface {
	// GeneratePrompt returns the combined prompt for the given technology.
	GeneratePrompt(ctx context.Context, technology string) ([]byte, error)
}
