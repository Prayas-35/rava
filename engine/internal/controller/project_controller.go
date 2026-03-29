package controller

import (
	"context"
	"errors"

	"github.com/Prayas-35/ragkit/engine/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ProjectController struct{}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

func (pc *ProjectController) CreateProject() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authenticated user")
		}

		var body struct {
			Name        string `json:"name"`
			AgentPrompt string `json:"agent_prompt"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if body.Name == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name is required")
		}

		req := service.CreateProjectRequest{
			UserID:      userID,
			Name:        body.Name,
			AgentPrompt: body.AgentPrompt,
		}

		project, err := service.CreateProject(context.Background(), req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(project)
	}
}

func (pc *ProjectController) ListProjects() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authenticated user")
		}

		projects, err := service.ListProjectsByUser(context.Background(), userID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(projects)
	}
}

func (pc *ProjectController) UpdateAgentPrompt() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authenticated user")
		}

		projectID := c.Params("project_id")
		if projectID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "project_id is required")
		}

		var req struct {
			AgentPrompt string `json:"agent_prompt"`
		}
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if req.AgentPrompt == "" {
			return fiber.NewError(fiber.StatusBadRequest, "agent_prompt is required")
		}

		err := service.UpdateProjectAgentPromptForUser(context.Background(), projectID, userID, req.AgentPrompt)
		if err != nil {
			if errors.Is(err, service.ErrProjectNotFoundOrUnauthorized) {
				return fiber.NewError(fiber.StatusNotFound, err.Error())
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
