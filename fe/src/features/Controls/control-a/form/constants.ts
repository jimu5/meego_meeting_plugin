export const columns = [
  {
    title: "标题",
    dataIndex: "name",
  },
  {
    title: "大小",
    dataIndex: "size",
  },
  {
    title: "所有者",
    dataIndex: "owner",
  },
  {
    title: "更新日期",
    dataIndex: "updateTime",
  },
];

// 最大返回长度
export const MAX_RES_LENGTH = 50;

// 最大返回会议数量
export const MAX_RES_MEETING_LENGTH = 100;

export const mockData = [
  {
    id: "1",
    bind_operator: "string",
    calendar_event_desc: "会议描述",
    calendar_event_id: "string",
    calendar_event_name: "上线预告",
    calendar_id: "string",
    calendar_event_app_link: "http://www.xxx.com",
    calendar_event_organizer: {
      display_name: "李雷",
      user_id: "1111",
    },
    meeting_host_user: {
      id: {
        open_id: "string",
        union_id: "string",
        user_id: "string",
      },
      user_role: 0,
      user_type: 0,
    },
    meeting_minute_url: "string",
    meeting_participant_count: "0",
    meeting_participant_count_accumulated: "0",
    meeting_record_url: "string",
    meeting_status: 1,
    meeting_time: {
      create_time: "1608885566",
      end_time: "1608888867",
      start_time: "1608883322",
    },
    meeting_topic: "string",
    project_key: "string",
    work_item_id: "string",
    work_item_type_key: "string",
  },
  {
    id: "2",
    bind_operator: "string",
    calendar_event_desc: "技术评审描述",
    calendar_event_id: "string1",
    calendar_event_name: "技术评审",
    calendar_event_app_link: "http://www.xxx.com",
    calendar_id: "string1",
    calendar_event_organizer: {
      display_name: "李雷",
      user_id: "1111",
    },
    meeting_host_user: {
      id: {
        open_id: "string",
        union_id: "string",
        user_id: "string",
      },
      user_role: 0,
      user_type: 0,
    },
    meeting_minute_url: "string",
    meeting_participant_count: "0",
    meeting_participant_count_accumulated: "0",
    meeting_record_url: "string",
    meeting_status: 3,
    meeting_time: {
      create_time: "1608885566",
      end_time: "1608888867",
      start_time: "1608883322",
    },
    meeting_topic: "string",
    project_key: "string",
    work_item_id: "string",
    work_item_type_key: "string",
  },
];
