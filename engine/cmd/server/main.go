package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Prayas-35/ragkit/engine/config"
	"github.com/Prayas-35/ragkit/engine/internal/database"
	"github.com/Prayas-35/ragkit/engine/internal/middleware"
	"github.com/Prayas-35/ragkit/engine/internal/rabbitmq"
	"github.com/Prayas-35/ragkit/engine/internal/routes"
	"github.com/Prayas-35/ragkit/engine/internal/vector"
	"github.com/Prayas-35/ragkit/engine/internal/worker"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(
		fiber.Config{
			AppName: "rava",
		},
	)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3001,https://altr-ai.vercel.app,https://undimly-swirly-daine.ngrok-free.dev",
		AllowMethods:     strings.Join([]string{fiber.MethodGet, fiber.MethodHead, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodPatch, fiber.MethodOptions}, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-API-Key",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400, // 24 hours
	}))

	app.Use(logger.New())

	// ETag middleware - but SKIP for authenticated requests (user-specific data)
	app.Use(etag.New(etag.Config{
		Next: func(c *fiber.Ctx) bool {
			// Skip etag for authenticated requests
			if c.Get("Authorization") != "" {
				return true // Skip etag
			}
			return false
		},
	}))

	// Cache middleware - but SKIP caching for authenticated requests
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			// Don't cache if:
			// 1. Request has Authorization header (user-specific data)
			// 2. Request has noCache query param
			// 3. Request method is not GET
			if c.Get("Authorization") != "" {
				return true // Skip cache
			}
			if c.Query("noCache") == "true" {
				return true // Skip cache
			}
			if c.Method() != fiber.MethodGet {
				return true // Skip cache
			}
			return false
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	// Rate limiting middleware - 10 requests per second per API key
	app.Use(middleware.RateLimiting())

	cfg := config.LoadConfig()

	// Initialize database connection
	database.Connect()

	rabbitmqConn := rabbitmq.Connect()
	if rabbitmqConn != nil {
		log.Println("✅ Connected to RabbitMQ")
	}
	defer rabbitmqConn.Close()

	ch, err := rabbitmqConn.Channel()
	if err != nil {
		log.Fatal("Failed to open RabbitMQ channel:", err)
	}
	defer ch.Close()

	// Ensure the ingest queue exists
	_, err = ch.QueueDeclare(
		"rag_ingest",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare rag_ingest queue:", err)
	}

	// Start background worker for ingestion
	embedder, err := vector.NewGeminiEmbedder(context.Background())
	if err != nil {
		log.Fatal("Failed to create Gemini embedder:", err)
	}

	store := &vector.Store{
		DB:       database.DB,
		Embedder: embedder,
	}

	worker.StartWorker(ch, store)

	// Register application routes
	routes.Register(app, ch, store)

	app.Get("/health",
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status": "ok",
				"time":   time.Now(),
			})
		})

	// Start server
	log.Fatal(app.Listen(":" + cfg.ServerPort))
}
