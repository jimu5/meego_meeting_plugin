package common

import (
	"strconv"
	"time"
)

// 毫秒级时间戳字符串转秒级
func MillisecondToSecond(millisecond string) (string, error) {
	timeStamp, err := strconv.ParseInt(millisecond, 10, 64)
	if err != nil {
		return millisecond, err
	}
	secondStamp := time.UnixMilli(timeStamp).Unix()
	return strconv.FormatInt(secondStamp, 10), nil
}
