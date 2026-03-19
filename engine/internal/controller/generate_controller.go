package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/Prayas-35/ragkit/engine/internal/service"
	"github.com/Prayas-35/ragkit/engine/internal/vector"
	"github.com/gofiber/fiber/v2"
)

type GenerateController struct {
	Store *vector.Store
}

type QueryRequest struct {
	Question string   `json:"question"`
	TopK     int      `json:"top_k"`
	History  []string `json:"history"`
}

func NewGenerateController(store *vector.Store) *GenerateController {
	return &GenerateController{Store: store}
}

func (ctrl *GenerateController) Query() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if strings.TrimSpace(apiKey) == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "X-API-Key header is required")
		}

		projectID, err := service.ResolveProjectIDByAPIKey(context.Background(), apiKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or revoked API key")
		}

		var req QueryRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		question := strings.TrimSpace(req.Question)
		if question == "" {
			return fiber.NewError(fiber.StatusBadRequest, "question is required")
		}

		limit := req.TopK
		if limit <= 0 {
			limit = 5
		}
		if limit > 20 {
			limit = 20
		}

		results, err := ctrl.Store.VectorSearch(
			c.Context(),
			projectID,
			question,
			limit,
		)

		if err != nil {
			fmt.Println("vector search error:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "failed to run vector search")
		}

		var chunks []string
		for _, r := range results {
			chunks = append(chunks, r.Content)
		}

		agentPrompt, err := service.GetProjectAgentPrompt(c.Context(), projectID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get project agent prompt")
		}

		answer, err := service.AnswerQuestion(c.Context(), question, chunks, req.History, agentPrompt)
		if err != nil {
			return fiber.NewError(fiber.StatusBadGateway, "failed to generate answer")
		}

		return c.JSON(fiber.Map{
			"answer": answer,
		})
	}
}
