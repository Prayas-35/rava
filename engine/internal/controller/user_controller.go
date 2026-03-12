package controller

import (
	"context"

	"github.com/Prayas-35/ragkit/engine/internal/service"
	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc *UserController) CreateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req service.CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if req.Email == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "email and password are required")
		}

		user, err := service.CreateUser(context.Background(), req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	}
}
