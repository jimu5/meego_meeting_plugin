package service

const (
	CalendarTypePrimary = "primary"
	LarkAuthURL         = "https://open.feishu.cn/open-apis/authen/v1/index"
)

type calendarContent struct {
	Summary   string `json:"summary,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}
