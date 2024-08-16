package dal

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm/clause"
	"meego_meeting_plugin/model"
)

type CalendarBindDao struct {
}

func NewCalendarBindDao() CalendarBindDao {
	return CalendarBindDao{}
}

func (c CalendarBindDao) MGetCalendarBindByWorkItemIDs(ctx context.Context, workItemIDs []int64) ([]*model.CalendarBind, error) {
	result := make([]*model.CalendarBind, 0)
	err := db.WithContext(ctx).Where("work_item_id in (?)", workItemIDs).Find(&result).Error
	return result, err
}

func (c CalendarBindDao) GetBindCalendarByProjectKeyAndWorkItemTypeKeyAndWorkItemID(ctx context.Context, projectKey, workItemTypeKey string,
	workItemID int64) ([]*model.CalendarBind, error) {
	result := make([]*model.CalendarBind, 0)
	err := db.WithContext(ctx).Where("project_key = ? and work_item_type_key = ? and work_item_id = ? and bind = 1", projectKey,
		workItemTypeKey, workItemID).Order("calendar_event_start_time desc").Find(&result).Error
	if err != nil {
		log.Errorf("[CalendarBind] GetBindCalendarByProjectKeyAndWorkItemTypeKeyAndWorkItemID err, err: %v", err)
		return result, err
	}
	return result, nil
}

func (c CalendarBindDao) CountMeetingByCalendarEventID(ctx context.Context, calendarEventIDs []string) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(&model.VCMeeting{}).Where("calendar_event_id in (?)", calendarEventIDs).Count(&count).Error
	if err != nil {
		log.Errorf("[CalendarBind] CountMeetingByCalendarEventID err, err: %v", err)
		return count, err
	}
	return count, nil
}

// TODO: 这个接口后面改成分页形式
func (c CalendarBindDao) MGetMeetingByCalendarEventID(ctx context.Context, calendarEventIDs []string) ([]*model.VCMeeting, error) {
	result := make([]*model.VCMeeting, 0)
	err := db.WithContext(ctx).Where("calendar_event_id in (?)", calendarEventIDs).Find(&result).Error
	if err != nil {
		log.Errorf("[CalendarBind] MGetMeetingByCalendarEventID err, err: %v", err)
		return nil, err
	}
	return result, nil
}

func (c CalendarBindDao) GetCalendarBindByWorkItemIDAndCalendarEventID(ctx context.Context, workItemID int64, calendarEventID string) (model.CalendarBind, error) {
	result := model.CalendarBind{}
	err := db.WithContext(ctx).Where("work_item_id = ? and calendar_event_id = ?", workItemID, calendarEventID).First(&result).Error
	if err != nil {
		log.Errorf("[CalendarBind] GetCalendarBindByWorkItemIDAndCalendarEventID, err, err: %v", err)
		return result, err
	}
	return result, nil
}

func (c CalendarBindDao) MGetCalendarMeetingsByCalendarEventID(ctx context.Context, eventID string) ([]*model.VCMeeting, error) {
	result := make([]*model.VCMeeting, 0)
	err := db.WithContext(ctx).Where("calendar_event_id in (?)", eventID).Find(&result).Error
	return result, err
}

func (c CalendarBindDao) CreateOrUpdateCalendarBind(ctx context.Context, bind *model.CalendarBind, operator string) error {
	if bind == nil {
		return nil
	}
	bind.UpdateBy = operator
	if bind.ID == 0 {
		// FIXME: createby 可能会有问题
		bind.CreateBy = operator
	}
	err := db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "work_item_id"}, {Name: "calendar_event_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"calendar_event_data", "bind", "calendar_event_start_time", "update_by"})}).
		Create(bind).Error
	if err != nil {
		log.Errorf("[CalendarBind] CreateOrUpdateCalendarBind err, err: %v", err)
	}
	return err
}

func (c CalendarBindDao) CreateOrUpdateCalendarMeetings(ctx context.Context, meetings []*model.VCMeeting, opeartor string) error {
	if len(meetings) == 0 {
		return nil
	}
	for _, m := range meetings {
		if m == nil {
			continue
		}
		m.UpdateBy = opeartor
		if m.ID == 0 {
			// FIXME: create by 可能会有问题
			m.CreateBy = opeartor
		}
	}
	err := db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "calendar_event_id"}, {Name: "meeting_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"meeting_data", "record_info", "update_by"}),
		}).
		Create(meetings).Error
	if err != nil {
		log.Errorf("[CalendarBind] CreateOrUpdateCalendarMeetings err, err: %v", err)
	}
	return err
}
func (c CalendarBindDao) UnbindByCalendarEventIDAndWorkItemID(ctx context.Context, calendarEventID string, workItemID int64) error {
	err := db.WithContext(ctx).Model(&model.CalendarBind{}).Where("calendar_event_id = ? and work_item_id = ?", calendarEventID, workItemID).
		Update("bind", false).Error
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c CalendarBindDao) GetBindByCalendarEventID(ctx context.Context, calendarEventID string) (model.CalendarBind, error) {
	bind := model.CalendarBind{}
	err := db.WithContext(ctx).Where("calendar_event_id = ?", calendarEventID).First(&bind).Error
	if err != nil {
		log.Error(bind)
	}
	return bind, err
}

func (c CalendarBindDao) GetRealNoRecordMeetingByMeetingIDs(ctx context.Context, meetingIDs []string) ([]*model.VCMeeting, error) {
	result := make([]*model.VCMeeting, 0)
	if len(meetingIDs) == 0 {
		return result, nil
	}
	err := db.WithContext(ctx).Where("meeting_id in (?) and record_info = '{}'", meetingIDs).Find(&result).Error
	if err != nil {
		log.Error(err)
		return result, err
	}
	// 不需要 filter, 现在好像没有返
	//filterResult := make([]*model.VCMeeting, 0, len(result))
	//for _, m := range result {
	//	if m == nil || m.MeetingData == nil || m.MeetingData.Ability == nil {
	//		continue
	//	}
	//	if m.MeetingData.Ability.UseRecording != nil && *m.MeetingData.Ability.UseRecording {
	//		filterResult = append(filterResult, m)
	//	}
	//}
	return result, nil
}

func (c CalendarBindDao) UpdateMeetingsRecordInfo(ctx context.Context, meetings []*model.VCMeeting) error {
	if len(meetings) == 0 {
		return nil
	}
	// 过滤掉 record 没值的
	filterMeetings := make([]*model.VCMeeting, 0, len(meetings))
	for _, m := range meetings {
		if m == nil || m.RecordInfo == nil || m.RecordInfo.Url == nil || len(*m.RecordInfo.Url) == 0 {
			continue
		}
		filterMeetings = append(filterMeetings, m)
	}
	if len(filterMeetings) == 0 {
		return nil
	}

	log.Info("handle UpdateMeetingsRecordInfo")
	err := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"record_info"}),
	}).Create(filterMeetings).Error
	if err != nil {
		log.Errorf("[UpdateMeetingsRecordInfo] err, err: %v", err)
		return err
	}
	return nil
}
