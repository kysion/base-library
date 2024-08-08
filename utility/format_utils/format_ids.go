package format_utils

import (
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
)

// BuildIdsToStrIds 将Int类型的Ids ---》 转换为string类型的Ids
func BuildIdsToStrIds(ids []int64) []string {
	strIds := make([]string, len(ids))
	for i, id := range ids {
		strIds[i] = fmt.Sprint(id)
	}
	slice := garray.NewStrArrayFrom(strIds).Unique().Slice()

	return slice
}
