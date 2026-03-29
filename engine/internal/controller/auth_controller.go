package controller

import (
	"context"
	"errors"
	"time"

	"github.com/Prayas-35/ragkit/engine/config"
	"github.com/Prayas-35/ragkit/engine/internal/auth"
	"github.com/Prayas-35/ragkit/engine/internal/service"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ac *AuthController) SignUp() fiber.Handler {
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

		token, err := auth.GenerateToken(user.ID, user.Email, config.LoadConfig().JWTSecret, 7*24*time.Hour)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"token": token,
			"user":  user,
		})
	}
}

func (ac *AuthController) SignIn() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if req.Email == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "email and password are required")
		}

		user, err := service.AuthenticateUser(context.Background(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				return fiber.NewError(fiber.StatusUnauthorized, "invalid email or password")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		token, err := auth.GenerateToken(user.ID, user.Email, config.LoadConfig().JWTSecret, 7*24*time.Hour)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
			"user":  user,
		})
	}
}
