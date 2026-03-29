package middleware

import (
	"strings"

	"github.com/Prayas-35/ragkit/engine/config"
	"github.com/Prayas-35/ragkit/engine/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func RequireJWT() fiber.Handler {
	cfg := config.LoadConfig()
	jwtSecret := cfg.JWTSecret

	return func(c *fiber.Ctx) error {
		authHeader := strings.TrimSpace(c.Get("Authorization"))
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Authorization header is required")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return fiber.NewError(fiber.StatusUnauthorized, "Authorization header must be Bearer <token>")
		}

		claims, err := auth.ParseToken(parts[1], jwtSecret)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		return c.Next()
	}
}
