package dal

import (
	"context"
	"testing"
)

func TestCalendarBind(t *testing.T) {
	ctx := context.Background()
	t.Run("", func(t *testing.T) {
		CalendarBind.GetBindCalendarByProjectKeyAndWorkItemTypeKeyAndWorkItemID(ctx, "65805ac8683445a834840b3d", "story", 26962565)
	})
}
