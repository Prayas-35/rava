package routes

import (
	"github.com/Prayas-35/ragkit/engine/internal/controller"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Register(app *fiber.App, ch *amqp.Channel) {
	documentController := controller.NewDocumentController(ch)
	userController := controller.NewUserController()
	projectController := controller.NewProjectController()
	apiKeyController := controller.NewAPIKeyController()

	api := app.Group("/api")
	api.Post("/users", userController.CreateUser())
	api.Post("/projects", projectController.CreateProject())
	api.Post("/projects/:project_id/keys", apiKeyController.CreateAPIKey())
	api.Post("/ingest", documentController.IngestDocument())
}
