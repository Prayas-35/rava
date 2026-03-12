package controller

import (
	"context"

	"github.com/Prayas-35/ragkit/engine/internal/service"
	"github.com/gofiber/fiber/v2"
)

type APIKeyController struct{}

func NewAPIKeyController() *APIKeyController {
	return &APIKeyController{}
}

// CreateAPIKey issues a new API key for a given project ID.
func (akc *APIKeyController) CreateAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		projectID := c.Params("project_id")
		if projectID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "project_id is required")
		}

		rawKey, apiKey, err := service.CreateAPIKey(context.Background(), projectID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"api_key":    rawKey,
			"id":         apiKey.ID,
			"project_id": apiKey.ProjectID,
		})
	}
}
