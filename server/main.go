package main

import (
	"github.com/gofiber/fiber/v2/log"
	"meego_meeting_plugin/config"
	_ "meego_meeting_plugin/dal"
)

func main() {
	config.InitConfig()

	fiberApp := NewFiberAPP()
	WithSwagger(fiberApp)
	log.Fatal(fiberApp.Listen("localhost:7999"))
}
