package base_hook

import (
	"github.com/gogf/gf/v2/container/garray"
	"github.com/kysion/base-library/base_model"
)

type BaseHook[T any, F any] struct {
	hookArr garray.Array
}

// InstallHook 安装Hook
func (s *BaseHook[T, F]) InstallHook(filter T, hookFunc F) {
	item := base_model.KeyValueT[T, F]{Key: filter, Value: hookFunc}
	s.hookArr.Append(item)
}

// UnInstallHook 卸载Hook
func (s *BaseHook[T, F]) UnInstallHook(filter T, f ...func(filter T, key T) bool) {
	newFuncArr := garray.NewArray()
	s.hookArr.Iterator(func(key int, value interface{}) bool {
		item := value.(base_model.KeyValueT[T, F])

		if len(f) > 0 && f[0](filter, item.Key) == false {
			newFuncArr.Append(value)
		}

		return true
	})
	s.hookArr = *newFuncArr
}

// ClearAllHook 清除Hook
func (s *BaseHook[T, F]) ClearAllHook() {
	s.hookArr.Clear()
}

// Iterator 遍历Hook
func (s *BaseHook[T, F]) Iterator(f func(key T, value F)) {
	s.hookArr.Iterator(func(key int, value interface{}) bool {
		item := value.(base_model.KeyValueT[T, F])

		f(item.Key, item.Value)

		return true
	})
}

func (s *BaseHook[T, F]) Where(filter T, f func(filter T, key T) bool) []F {
	result := make([]F, 0)

	s.hookArr.Iterator(func(key int, value interface{}) bool {
		item := value.(base_model.KeyValueT[T, F])

		if f(filter, item.Key) {
			result = append(result, item.Value)
		}

		return true
	})

	return result
}
