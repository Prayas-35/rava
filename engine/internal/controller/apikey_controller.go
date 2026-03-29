package controller

import (
	"context"

	"github.com/Prayas-35/ragkit/engine/config"
	"github.com/Prayas-35/ragkit/engine/internal/service"
	"github.com/gofiber/fiber/v2"
)

type APIKeyController struct{}

func NewAPIKeyController() *APIKeyController {
	return &APIKeyController{}
}

func getKeySecret() string {
	cfg := config.LoadConfig()
	if cfg.APIKeySecret != "" {
		return cfg.APIKeySecret
	}
	return cfg.JWTSecret
}

// CreateAPIKey issues a new API key for a given project ID.
func (akc *APIKeyController) CreateAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authenticated user")
		}

		projectID := c.Params("project_id")
		if projectID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "project_id is required")
		}

		owns, err := service.UserOwnsProject(context.Background(), projectID, userID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if !owns {
			return fiber.NewError(fiber.StatusForbidden, "project not found or unauthorized")
		}

		rawKey, apiKey, err := service.CreateAPIKey(context.Background(), projectID, getKeySecret())
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

func (akc *APIKeyController) ListAPIKeys() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authenticated user")
		}

		projectID := c.Params("project_id")
		if projectID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "project_id is required")
		}

		owns, err := service.UserOwnsProject(context.Background(), projectID, userID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if !owns {
			return fiber.NewError(fiber.StatusForbidden, "project not found or unauthorized")
		}

		keys, err := service.ListAPIKeysByProject(context.Background(), projectID, getKeySecret())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(keys)
	}
}
