package routes

import (
	"github.com/Prayas-35/ragkit/engine/internal/controller"
	"github.com/Prayas-35/ragkit/engine/internal/vector"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Register(app *fiber.App, ch *amqp.Channel, store *vector.Store) {
	documentController := controller.NewDocumentController(ch)
	userController := controller.NewUserController()
	projectController := controller.NewProjectController()
	apiKeyController := controller.NewAPIKeyController()
	generateController := controller.NewGenerateController(store)

	api := app.Group("/api")
	api.Post("/users", userController.CreateUser())
	api.Post("/projects", projectController.CreateProject())
	api.Post("/projects/:project_id/keys", apiKeyController.CreateAPIKey())
	api.Put("/projects/:project_id/agent-prompt", projectController.UpdateAgentPrompt())
	api.Post("/ingest", documentController.IngestDocument())
	api.Post("/query", generateController.Query())
}
