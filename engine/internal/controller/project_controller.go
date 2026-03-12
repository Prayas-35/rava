package controller

import (
	"context"

	"github.com/Prayas-35/ragkit/engine/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ProjectController struct{}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

func (pc *ProjectController) CreateProject() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req service.CreateProjectRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if req.UserID == "" || req.Name == "" {
			return fiber.NewError(fiber.StatusBadRequest, "user_id and name are required")
		}

		project, err := service.CreateProject(context.Background(), req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(project)
	}
}
