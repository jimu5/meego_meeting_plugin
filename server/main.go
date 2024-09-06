package main

import (
	"meego_meeting_plugin/config"
	"meego_meeting_plugin/dal"
	_ "meego_meeting_plugin/dal"

	"github.com/gofiber/fiber/v2/log"
)

func main() {
	config.InitConfig()
	dal.InitDB()

	fiberApp := NewFiberAPP()
	WithSwagger(fiberApp)
	log.Fatal(fiberApp.Listen("localhost:7999"))
}
