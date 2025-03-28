package main

import (
	"meego_meeting_plugin/handler"
	"meego_meeting_plugin/mw"

	"github.com/gofiber/fiber/v2"
)

func SetupAPIRouter(app *fiber.App) {
	apiRoute := app.Group("/api", mw.ErrorHandle)
	v1Route := apiRoute.Group("/v1")
	{
		v1Route.All("/ping", handler.GetPing)
		larkAPI := v1Route.Group("/lark")
		larkCalendarRoute := larkAPI.Group("/calendar", mw.GetUser)
		{
			larkCalendarRoute.Get("/search", handler.CalendarSearch)
		}
		MeegoRoute := v1Route.Group("/meego")
		{
			MeegoRoute.Post("/calendar_event/bind", mw.GetUser, handler.BindCalendarEventWithWorkItem)
			MeegoRoute.Post("/calendar_event/unbind", mw.GetUser, handler.UnBindCalendarEventWithWorkItem)
			MeegoRoute.Post("/work_item_meetings", handler.ListWorkItemMeetings)
			MeegoRoute.Post("/work_item_meetings/refresh", mw.GetUser, handler.RefreshCalendar)
			MeegoRoute.Post("/work_item_meetings/chat_auto_bind", mw.GetUser, handler.ChatAutoBindCalendar)
			MeegoRoute.Get("/work_item_meetings/chat_auto_bind", handler.GetAutoBindCalendarStatus)
		}
		MeegoLarkRoute := MeegoRoute.Group("/lark")
		{
			MeegoLarkRoute.Get("/auth", handler.MeegoLarkLogin)
		}

		// 事件
		larkAPI.Post("/webhook/event", NewEventHandlerFunc(handler.LarkEventHandler))
	}

}
