package service

import "strings"

func getMeetingNOByMeetingUrl(meetingUrl string) string {
	split := strings.Split(meetingUrl, "/")
	return split[len(split)-1]
}
