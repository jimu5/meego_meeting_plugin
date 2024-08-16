package main

import (
	"github.com/gofiber/fiber/v2/log"
	_ "meego_meeting_plugin/dal"
)

func main() {
	fiberApp := NewFiberAPP()
	WithSwagger(fiberApp)
	log.Fatal(fiberApp.Listen("localhost:7999"))
}
