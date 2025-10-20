package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewFiberAPP() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:     false,
		BodyLimit:   1024 * 1024 * 1024,
		ProxyHeader: "X-Real-IP",
	})
	app.Use(logger.New(logger.ConfigDefault))
	log.SetLevel(log.LevelInfo)
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(cors.New())

	app.Get("/status", monitor.New())
	SetupAPIRouter(app)

	return app
}
