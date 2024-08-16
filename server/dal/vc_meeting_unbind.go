package dal

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"meego_meeting_plugin/model"
)

type VCMeetingUnbindDao struct {
}

func NewVCMeetingUnbindDao() VCMeetingUnbindDao {
	return VCMeetingUnbindDao{}
}

func (v VCMeetingUnbindDao) GetVCMeetingUnbindInfoByWorkItemID(ctx context.Context, workItemID int64) ([]*model.VCMeetingUnBind, error) {
	result := make([]*model.VCMeetingUnBind, 0)
	err := db.WithContext(ctx).Where("work_item_id = ?", workItemID).Find(&result).Error
	if err != nil {
		log.Error(err)
	}
	return result, err
}

// 含义不纯粹了, 先这样吧
func (v VCMeetingUnbindDao) SaveUnbindVCMeetings(ctx context.Context, workItemID int64, meetings []*model.VCMeeting) error {
	info := make([]*model.VCMeetingUnBind, 0, len(meetings))
	for _, m := range meetings {
		if m == nil {
			continue
		}
		info = append(info, &model.VCMeetingUnBind{
			MeetingID:       m.MeetingID,
			CalendarEventID: m.CalendarEventID,
			CalendarID:      m.CalendarID,
			WorkItemID:      workItemID,
		})
	}
	err := db.WithContext(ctx).Create(info).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

func (v VCMeetingUnbindDao) DeleteMeetingsByWorkItemIDAndMeetingIDs(ctx context.Context, workItemID int64, meetingIDs []string) error {
	err := db.WithContext(ctx).Where("work_item_id =? and meeting_id in (?)", workItemID, meetingIDs).Unscoped().Delete(&model.VCMeetingUnBind{}).Error
	if err != nil {
		log.Errorf("[DeleteMeetingsByWorkItemIDAndMeetingIDs] delete VCMeetingUnBind err, err: %v", err)
		return err
	}
	return err
}

func (v VCMeetingUnbindDao) DeleteMeetingsByWorkItemID(ctx context.Context, workItemID int64) error {
	err := db.WithContext(ctx).Where("work_item_id =?", workItemID).Unscoped().Delete(&model.VCMeetingUnBind{}).Error
	if err != nil {
		log.Errorf("[DeleteMeetingsByWorkItemID] delete VCMeetingUnBind err, err: %v", err)
		return err
	}
	return nil
}
