package format_utils

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

// GetWeekDay 传入指定的时间，返回具体是一周内的第几天
func GetWeekDay(t *gtime.Time) int {
	weekDay := t.Weekday().String()

	switch weekDay {
	case "Monday":
		return 1
	case "Tuesday":
		return 2
	case "Wednesday":
		return 3
	case "Thursday":
		return 4
	case "Friday":
		return 5
	case "Saturday":
		return 6
	case "Sunday":
		return 7
	}

	return 0
}

// GetQuarter 传入指定的时间，返回具体的季度
func GetQuarter(t *gtime.Time) int {
	month := gconv.Int(t.Format("n"))
	if month <= 3 {
		return 1
	} else if month > 3 && month <= 6 {
		return 2
	} else if month > 6 && month <= 9 {
		return 3
	} else if month > 9 && month <= 12 {
		return 4
	}
	return 0
}

// SecondToDuration 将秒数转化为 duration 对象
func SecondToDuration(second int) time.Duration {
	t := time.Duration(second) * time.Second

	return t
}
