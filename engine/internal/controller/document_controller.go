package controller

import (
	"context"

	"github.com/Prayas-35/ragkit/engine/internal/service"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DocumentController struct {
	ch *amqp.Channel
}

func NewDocumentController(ch *amqp.Channel) *DocumentController {
	return &DocumentController{ch: ch}
}

// IngestDocument handles HTTP requests to create a document and queue it for ingestion.
func (dc *DocumentController) IngestDocument() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "X-API-Key header is required")
		}

		projectID, err := service.ResolveProjectIDByAPIKey(context.Background(), apiKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or revoked API key")
		}

		var req service.DocumentIngestRequest

		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if req.Name == "" || req.Content == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name and content are required")
		}

		docID, err := service.QueueDocumentIngestion(dc.ch, projectID, req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"status":      "queued",
			"document_id": docID,
		})
	}
}
