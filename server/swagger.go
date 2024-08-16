package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "meego_meeting_plugin/docs"
)

func WithSwagger(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	//app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
	//	DeepLinking: false,
	//	// Expand ("list") or Collapse ("none") tag groups by default
	//	DocExpansion: "none",
	//	// Prefill OAuth ClientId on Authorize popup
	//}))
}
