package base_funs

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/kysion/base-library/base_model"
	"reflect"
	"sync"
	"time"
)

func If[R any](condition bool, trueVal, falseVal R) R {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}

// SearchFilterEx 支持增加拓展字段，会提炼拓展字段的最新过滤题条件模型返回，并支持从原始的过滤模型剔除不需要的条件
func SearchFilterEx(search *base_model.SearchParams, fields ...string) *base_model.SearchParams {
	result := &base_model.SearchParams{}
	newFilter := make([]base_model.FilterInfo, 0)
	newSearchFilter := make([]base_model.FilterInfo, 0)
	newSearchFilterStr := garray.NewStrArray()

	for _, info := range search.Filter {
		//count := len(result.Filter)
		ss := true
		for _, field := range fields {
			split := gstr.Split(field, ".")
			if gstr.ToLower(gstr.CaseSnake(split[0])) == gstr.ToLower(gstr.CaseSnake(info.Field)) && len(split) > 1 {
				ss = false
			}
			if gstr.ToLower(gstr.CaseSnake(info.Field)) == gstr.ToLower(gstr.CaseSnake(split[0])) {
				newFilter = append(newFilter, info)
				break
			}

		}
		//if count == len(result.Filter) {
		//  newFilter = append(newFilter, info)
		//}

		if ss {
			if !newSearchFilterStr.Contains(info.Field) {
				newSearchFilterStr.Append(info.Field)
				newSearchFilter = append(newSearchFilter, info)
			}
		}
	}

	search.Filter = newSearchFilter
	result.Filter = newFilter

	return result
}

// ByteCountSI 以1000作为基数
func ByteCountSI[T int64 | uint64](b T) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

// ByteCountIEC 以1024作为基数
func ByteCountIEC[T int64 | uint64](b T) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func RemoveSliceAt[T int | int64 | string | uint | uint64](slice []T, elem T) []T {
	if len(slice) == 0 {
		return slice
	}

	for i, v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveSliceAt(slice, elem)
			break
		}
	}
	return slice
}

func AttrBuilder[T any, TP any](ctx context.Context, key string, builder ...func(data TP)) context.Context {
	key = key + "::" + reflect.ValueOf(new(T)).Type().String() + "::" + reflect.ValueOf(new(TP)).Type().String()
	key = gstr.Replace(key, "*", "")
	def := func(data TP) {}

	if len(builder) > 0 {
		def = builder[0]
	}

	return context.WithValue(ctx, key,
		base_model.KeyValueT[string, func(data TP)]{
			Key:   key,
			Value: def,
		},
	)
}

// union_main_id::co_model.EmployeeRes::[]co_model.Team
func AttrMake[T any, TP any](ctx context.Context, key string, builder func() TP) {
	key = key + "::" + reflect.ValueOf(new(T)).Type().String() + "::" + reflect.ValueOf(new(TP)).Type().String()
	key = gstr.Replace(key, "*", "")
	v := ctx.Value(key)

	data, has := v.(base_model.KeyValueT[string, func(data TP)])
	if v != nil && has {
		data.Value(builder())
	}
}

/*
 Debounce 防抖函数
	该函数是一个防抖函数，用于延迟函数的执行。
	它返回一个闭包函数，接收一个函数参数f。
	在闭包函数中，通过互斥锁保证并发安全，停止之前的定时器，创建一个新的定时器，并在指定的时间间隔后执行函数f。
	这样可以避免在短时间内频繁调用函数f，达到防抖的效果。
 通俗理解：
	可以想象成一个需要“冷静一下”的机制。
	就好比你在不停地按一个按钮，但系统只在你最后一次按完并经过一段安静时间（防抖时间间隔）后，才真正去执行相应的操作。
	在这期间，不管你按得多频繁，只有最后一次按下去且等待一段时间没再按，才会触发实际动作。
*/
// Debounce 防抖函数
func Debounce(interval time.Duration) func(f func()) { // interval: 防抖时间间隔

	var l sync.Mutex
	var timer *time.Timer

	return func(f func()) {
		l.Lock()
		defer l.Unlock()
		// 使用lock保证d.timer更新之前一定先Stop.

		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(interval, f)
	}
}

//// Debounce 防抖函数 （优化版）
//func Debounce(interval time.Duration) func(f func()) {
//	var l sync.Mutex
//	ctx, cancel := context.WithCancel(context.Background())
//	var timer *time.Timer
//
//	fmt.Println(ctx)
//	return func(f func()) {
//		l.Lock()
//		defer l.Unlock()
//
//		// 使用lock保证d.timer更新之前一定先Stop.
//		if timer != nil {
//			timer.Stop()
//			timer = nil
//		}
//		//timer = time.AfterFunc(interval,f)
//		timer = time.AfterFunc(interval, func() {
//			// 在函数内部释放context，确保资源可以被清理
//			cancel()
//			f()
//		})
//
//		// 立即取消context以停止旧的timer（如果存在），这有助于减少资源泄露
//		// 注意：这并不能保证timer在启动前一定被取消，因此仍存在一定的竞争条件
//		cancel()
//		ctx, cancel = context.WithCancel(context.Background())
//	}
//}
