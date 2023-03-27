package utils

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

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
