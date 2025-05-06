package common

import (
	"fmt"
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
	return strconv.FormatInt(int64(secondStamp), 10), nil
}

// 秒级时间戳
func ExpandSecondTimeStamp(input string, duration time.Duration) string {
	timeStamp, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return input
	}
	secondTime := time.Unix(timeStamp, 0).Add(duration).Unix()
	return fmt.Sprintf("%d", secondTime)
}
