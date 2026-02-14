package agent

import (
	"context"

	"github.com/tmc/langchaingo/llms"
)

type BaseAgent struct {
	LLM llms.Model
}

func NewBaseAgent(llm llms.Model) *BaseAgent {
	return &BaseAgent{LLM: llm}
}

func (a *BaseAgent) Generate(ctx context.Context, prompt string) (string, error) {
	// 使用 GenerateFromSinglePrompt 简化调用
	return llms.GenerateFromSinglePrompt(ctx, a.LLM, prompt)
}
