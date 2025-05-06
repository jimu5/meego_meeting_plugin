package common

import (
	"fmt"
	"testing"
	"time"
)

func TestFunc(t *testing.T) {
	tm, err := MillisecondToSecond("1746622800000")
	if err != nil {

	}
	fmt.Println(tm)
	tm = ExpandSecondTimeStamp(tm, time.Minute)
	fmt.Println(tm)
}
