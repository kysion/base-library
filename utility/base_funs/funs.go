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

// Contains 检查给定的切片中是否包含指定的元素。
//
// 参数：
//
//	slice: 要检查的切片，可以是字符串、整数、浮点数等类型。
//	element: 要查找的元素，类型与切片中的元素相同。
//
// 返回值：
//
//	如果切片中包含指定的元素，则返回 true，否则返回 false。
//
// 注意：
//
//	该函数使用泛型定义，支持多种类型的操作，提高了代码的通用性和复用性。
func Contains[T comparable](slice []T, element T) bool {
	// 检查切片是否为空或长度为0，如果是，则直接返回 false，因为切片中不可能包含指定元素。
	if slice == nil || len(slice) == 0 {
		return false
	}

	// 遍历切片中的每个元素，如果找到与指定元素相等的元素，则返回 true。
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	// 如果遍历完切片后没有找到指定元素，则返回 false。
	return false
}

// Unique 函数用于去除给定切片中的重复元素，返回一个只包含唯一元素的新切片。
// 它使用了泛型 T，其中 T 需要实现 comparable 接口，以便可以比较两个元素是否相同。
// 参数 slice: 待处理的切片，其元素类型为泛型 T。
// 返回值: 一个新切片，包含输入切片中的唯一元素。
func Unique[T comparable](slice []T) []T {
	// 使用空结构体来节省空间，用作标记已经见过的元素
	seen := make(map[T]struct{})
	// 初始化结果切片，用于存储唯一元素
	var result []T

	// 遍历输入切片中的每个元素
	for _, v := range slice {
		// 检查当前元素是否已经存在于 seen 中
		if _, exists := seen[v]; !exists {
			// 如果不存在，则将其添加到 seen 中，并添加到结果切片中
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	// 返回结果切片
	return result
}

// FilterEmpty 过滤掉字符串切片中的空字符串。
// 参数:
//
//	slice: 待过滤的字符串切片。
//
// 返回值:
//
//	一个新的不包含空字符串的字符串切片。
func FilterEmpty(slice []string) []string {
	// 初始化一个空的字符串切片，用于存储非空字符串。
	var result []string
	// 遍历输入的字符串切片。
	for _, s := range slice {
		// 如果当前字符串非空，则将其添加到结果切片中。
		if s != "" {
			result = append(result, s)
		}
	}
	// 返回最终的结果切片。
	return result
}

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
