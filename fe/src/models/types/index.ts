export enum LinkType {
  MIAOJI = "miaoji",
  DOC = "doc",
  APP_LINK = "app_link",
}

// 会议状态, 可选值 1(呼叫中), 2(进行中) 3(已结束)
export enum StatusMap {
  CALLED = 1,
  IN_PROGRESS = 2,
  FINISH = 3,
}

export interface ResWrapper<T = {}> {
  message: string;
  statusCode: number;
  data: T;
}

export interface AuthRes {
  code: string;
  state: string;
}
export interface GetTokenAuthRes {
  token: string;
  expire_time: number;
}
export interface GetAuthCodeRes {
  code: string;
  redirect_uri: string;
  state: string;
}

export interface UserPluginTokenRes {
  token: string;
  expire_time: number;
  refresh_token: string;
  refresh_token_expire_time: number;
  user_key?: string;
  saas_tenant_key?: string;
}

export interface SearchCalendarParams {
  page_size: number;
  page_token?: string;
  query_word?: string;
}

export interface UpdateCalendarParams {
  project_key: string;
  work_item_id: number;
  work_item_type_key: string;
}

export interface BindCalendarParams {
  // 日程 ID
  calendar_event_id: string;
  // 日程组织者日历ID
  calendar_id: string;
  // 会议ID
  meeting_id?: string;
  // 空间信息
  project_key: string;
  work_item_id: number;
  work_item_type_key: string;
  with_after_recurring_event?: boolean;
}

export interface BindMeetingsParams {
  page_number: number;
  page_size: number;
  query_word?: string;
  // 空间信息
  project_key: string;
  work_item_id: number;
  work_item_type_key: string;
}

export interface AutoBindMeetingsParams {
  enable: boolean;
  project_key: string;
  work_item_id: number;
  work_item_type_key: string;
}

export type TimeDesc = {
  date: string;
  timestamp: string;
  timezone: string;
};

export interface ITimeRender {
  startTime: string;
  endTime: string;
}

export interface CalendarItem {
  // 日程ID
  event_id: string;
  // 日历组织者日历ID
  organizer_calendar_id: string;
  // 日程名称
  summary: string;
  // 日程开始时间
  start_time: TimeDesc;
  // 日程结束时间
  end_time: TimeDesc;
  // 日程状态
  status: string;
  // 重复日程的重复性规则；参考rfc5545；不支持COUNT和UNTIL同时出现；预定会议室重复日程长度不得超过两年。
  recurrence?: string;
}

export interface WorkItemMeeting {
  // 日程组织者信息
  calendar_event_organizer: {
    //     日程组织者姓名
    display_name: string;
    // 日程组织者user ID
    user_id: string;
  };
  // 日程的app_link
  calendar_event_app_link: string;
  // 是否是重复日程
  calendar_event_recurrence: string;
  // 关联日程操作人(这应该是一个user接口, 先暂时用 string 替代下)
  bind_operator: string;
  bind_operator_info: {
    avatar_url: "string";
    email: "string";
    meego_user_key: "string";
    name_cn: "string";
    name_en: "string";
  };
  // 日程描述(应该是对应的会议描述, 因为会议没有描述)
  calendar_event_desc: string;
  // 日程id
  calendar_event_id: string;
  // 日程名称
  calendar_event_name: string;
  // 日程组织日历ID
  calendar_id: string;
  // 会议ID
  meeting_id?: string;
  meeting_host_user: {
    // 会议主持人,用户 ID
    id: {
      open_id: string;
      union_id: string;
      user_id: string;
    };
    // 会议描述
    description: string;
    // 用户会中角色
    user_role: number;
    // 用户类型
    user_type: number;
  };
  // 会议纪要链接(看着openapi接口好像没有)
  meeting_minute_url: string;
  // 参会峰值人数
  meeting_participant_count: string;
  // 参会累计人数
  meeting_participant_count_accumulated: string;
  // 会议录制链接
  meeting_record_url: string;
  // 会议状态, 可选值 1(呼叫中), 2(进行中) 3(已结束)
  meeting_status: number;
  // 会议时间, 都是 unix 时间, 单位 sec
  meeting_time: {
    // 创建时间
    create_time: string;
    // 结束时间
    end_time: string;
    // 开始时间
    start_time: string;
  };
  // 会议主题
  meeting_topic: string;
  project_key: string;
  work_item_id: string;
  work_item_type_key: string;
}

export interface SearchCalendarRes {
  items: CalendarItem[];
  page_token: string;
}
export interface BindMeetingsParamsRes {
  meetings: WorkItemMeeting[];
  total: number;
  msg?: string;
}

export interface GetAutoBindMeetingsStatus {
  work_item_id: number;
}
export interface GetAutoBindMeetingsStatusRes {
  enable: boolean;
}
