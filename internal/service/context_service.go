package service

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"ai-rules-link/internal/utils"
)

// ContextService implements domain.ContextInitializer and domain.PromptGenerator.
type ContextService struct {
	RulesFS fs.FS
}

// NewContextService creates a new ContextService with the given embedded FS.
func NewContextService(rulesFS fs.FS) *ContextService {
	return &ContextService{RulesFS: rulesFS}
}

// GeneratePrompt combines the base and technology-specific prompts.
func (s *ContextService) GeneratePrompt(ctx context.Context, technology string) ([]byte, error) {
	basePrompt, err := fs.ReadFile(s.RulesFS, "rules/baserules.mdc")
	if err != nil {
		return nil, fmt.Errorf("read base prompt: %w", err)
	}

	techPromptPath := fmt.Sprintf("rules/prompt.%s.mdc", technology)
	techPrompt, err := fs.ReadFile(s.RulesFS, techPromptPath)
	if err != nil {
		return nil, fmt.Errorf("read tech prompt: %w", err)
	}

	return append(basePrompt, techPrompt...), nil
}

// GeneratePromptFlexible generates the prompt based on the provided options.
// If baseOnly is true, returns only the base prompt.
// If langOnly is true, returns only the language prompt.
// If both are false, returns base + language (default).
func (s *ContextService) GeneratePromptFlexible(ctx context.Context, language string, baseOnly, langOnly bool) ([]byte, error) {
	if baseOnly && langOnly {
		return nil, fmt.Errorf("cannot set both baseOnly and langOnly to true")
	}

	if baseOnly {
		basePrompt, err := fs.ReadFile(s.RulesFS, "rules/baserules.mdc")
		if err != nil {
			return nil, fmt.Errorf("read base prompt: %w", err)
		}
		return basePrompt, nil
	}

	techPromptPath := fmt.Sprintf("rules/prompt.%s.mdc", language)
	techPrompt, err := fs.ReadFile(s.RulesFS, techPromptPath)
	if err != nil {
		return nil, fmt.Errorf("read tech prompt: %w", err)
	}

	if langOnly {
		return techPrompt, nil
	}

	basePrompt, err := fs.ReadFile(s.RulesFS, "rules/baserules.mdc")
	if err != nil {
		return nil, fmt.Errorf("read base prompt: %w", err)
	}

	return append(basePrompt, techPrompt...), nil
}

// Initialize sets up the context for the given tool and technology.
func (s *ContextService) Initialize(ctx context.Context, tool, technology string) error {
	prompt, err := s.GeneratePrompt(ctx, technology)
	if err != nil {
		return err
	}

	projectContextDir := ".context"
	if err := os.MkdirAll(projectContextDir, 0755); err != nil {
		return fmt.Errorf("create context dir: %w", err)
	}

	if err := utils.EnsureGitignore(); err != nil {
		return fmt.Errorf("ensure .gitignore: %w", err)
	}

	projectPromptFile := filepath.Join(projectContextDir, "prompt.mdc")
	if err := utils.WriteBytes(projectPromptFile, prompt); err != nil {
		return fmt.Errorf("write prompt file: %w", err)
	}

	var symlinkPath string
	switch tool {
	case "gemini":
		symlinkPath = filepath.Join(".gemini", "context.mdc")
	case "cursor":
		symlinkPath = filepath.Join(".vscode", "cursor_prompt.mdc")
	default:
		return nil // No symlink for unsupported tools
	}

	if err := os.MkdirAll(filepath.Dir(symlinkPath), 0755); err != nil {
		return fmt.Errorf("create symlink dir: %w", err)
	}

	relPromptPath, err := filepath.Rel(filepath.Dir(symlinkPath), projectPromptFile)
	if err != nil {
		return fmt.Errorf("rel prompt path: %w", err)
	}

	if _, err := os.Lstat(symlinkPath); err == nil {
		if err := os.Remove(symlinkPath); err != nil {
			return fmt.Errorf("remove existing symlink: %w", err)
		}
	}

	if err := os.Symlink(relPromptPath, symlinkPath); err != nil {
		return fmt.Errorf("create symlink: %w", err)
	}

	return nil
}

// InitializeFlexible sets up the context for the given language, with options for baseOnly or langOnly.
func (s *ContextService) InitializeFlexible(ctx context.Context, language string, baseOnly, langOnly bool) error {
	prompt, err := s.GeneratePromptFlexible(ctx, language, baseOnly, langOnly)
	if err != nil {
		return err
	}

	projectContextDir := ".context"
	if err := os.MkdirAll(projectContextDir, 0755); err != nil {
		return fmt.Errorf("create context dir: %w", err)
	}

	if err := utils.EnsureGitignore(); err != nil {
		return fmt.Errorf("ensure .gitignore: %w", err)
	}

	projectPromptFile := filepath.Join(projectContextDir, "prompt.mdc")
	if err := utils.WriteBytes(projectPromptFile, prompt); err != nil {
		return fmt.Errorf("write prompt file: %w", err)
	}

	return nil
}
