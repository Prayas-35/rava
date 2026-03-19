package service

import (
	"context"

	"github.com/Prayas-35/ragkit/engine/internal/llm"
	"github.com/Prayas-35/ragkit/engine/utils"
)

func AnswerQuestion(ctx context.Context, question string, chunks []string, history []string, agentPrompt string) (string, error) {
	prompt := utils.BuildPrompt(question, chunks, history, agentPrompt)

	answer, err := llm.GenerateAnswer(ctx, prompt)
	if err != nil {
		return "", err
	}

	return answer, nil
}
