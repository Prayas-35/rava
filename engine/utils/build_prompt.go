package utils

import (
	"strings"
)

func BuildPrompt(question string, chunks []string) string {
	context := ""

	for _, c := range chunks {
		trimmed := strings.TrimSpace(c)
		if trimmed == "" {
			continue
		}
		context += trimmed + "\n\n"
	}

	if strings.TrimSpace(context) == "" {
		context = "No relevant context was retrieved. If the answer cannot be inferred from context, reply: 'I don't have enough context to answer that.'"
	}

	prompt := `
		You are a helpful assistant.

		Answer the question using ONLY the provided context.
		If the context does not contain the answer, say exactly: "I don't have enough context to answer that."
		Ensure your answer is concise and directly addresses the question.
		On asked who are you, say: "I am an assistant for answering questions based on provided context."

		Context:
		` + context + `

		Question:
		` + question

	return prompt
}
