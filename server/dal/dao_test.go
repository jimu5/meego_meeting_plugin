package dal

import (
	"context"
	"testing"

	"meego_meeting_plugin/config"
)

func TestCalendarBind(t *testing.T) {
	ctx := context.Background()
	t.Run("", func(t *testing.T) {
		CalendarBind.GetBindCalendarByProjectKeyAndWorkItemTypeKeyAndWorkItemID(ctx, "65805ac8683445a834840b3d", "story", 26962565)
	})
}

func TestGetUnprocessedTasks(t *testing.T) {
	ctx := context.Background()
	config.InitConfig()
	db = InitDB()
	t.Run("", func(t *testing.T) {
		PendingTask.GetUnprocessedTasksByMeegoUserKey(ctx, "")
	})
	t.Run("", func(t *testing.T) {
		PendingTask.GetUnprocessedTasks(ctx)
	})
}
