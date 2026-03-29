package routes

import (
	"github.com/Prayas-35/ragkit/engine/internal/controller"
	"github.com/Prayas-35/ragkit/engine/internal/middleware"
	"github.com/Prayas-35/ragkit/engine/internal/vector"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Register(app *fiber.App, ch *amqp.Channel, store *vector.Store) {
	documentController := controller.NewDocumentController(ch)
	userController := controller.NewUserController()
	authController := controller.NewAuthController()
	projectController := controller.NewProjectController()
	apiKeyController := controller.NewAPIKeyController()
	generateController := controller.NewGenerateController(store)

	api := app.Group("/api")
	api.Post("/users", userController.CreateUser())
	api.Post("/auth/signup", authController.SignUp())
	api.Post("/auth/signin", authController.SignIn())

	protected := api.Group("", middleware.RequireJWT())
	protected.Get("/projects", projectController.ListProjects())
	protected.Post("/projects", projectController.CreateProject())
	protected.Get("/projects/:project_id/keys", apiKeyController.ListAPIKeys())
	protected.Post("/projects/:project_id/keys", apiKeyController.CreateAPIKey())
	protected.Put("/projects/:project_id/agent-prompt", projectController.UpdateAgentPrompt())

	api.Put("/ingest", documentController.IngestDocument())
	api.Post("/query", generateController.Query())
}
