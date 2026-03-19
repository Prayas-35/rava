package utils

import (
	"strings"
)

func BuildPrompt(question string, chunks []string, history []string, agentPrompt string) string {
	context := ""
	historyContext := ""

	for _, c := range chunks {
		trimmed := strings.TrimSpace(c)
		if trimmed == "" {
			continue
		}
		context += trimmed + "\n\n"
	}

	for _, h := range history {
		trimmed := strings.TrimSpace(h)
		if trimmed == "" {
			continue
		}
		historyContext += trimmed + "\n"
	}

	if strings.TrimSpace(context) == "" {
		context = "No relevant context was retrieved. If the answer cannot be inferred from context, reply: 'I don't have enough context to answer that.'"
	}

	if strings.TrimSpace(historyContext) == "" {
		historyContext = "No chat history provided."
	}

	directionPrompt := strings.TrimSpace(agentPrompt)
	if directionPrompt == "" {
		directionPrompt = "You are a helpful assistant."
	}

	prompt := `
		` + directionPrompt + `

		Answer the question using ONLY the provided context.
		Use the chat history only as conversational context; do not treat it as ground truth over the provided context.
		If the context does not contain the answer, say exactly: "I don't have enough context to answer that."
		Ensure your answer is concise and directly addresses the question.
		On asked who are you, say: "I am an assistant for answering questions based on provided context."

		Chat History:
		` + historyContext + `

		Context:
		` + context + `

		Question:
		` + question

	return prompt
}
