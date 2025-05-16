package base_funs

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/kysion/base-library/base_model"
)

// 泛型工厂函数 - 支持创建泛型结构体实例
func CreateGenericInstance[T any]() (T, error) {
	var zero T
	// 获取泛型类型信息
	typ := reflect.TypeOf(zero).Elem()

	// 验证是否为结构体类型
	if typ.Kind() != reflect.Struct {
		return zero, fmt.Errorf("type %v is not a struct", typ)
	}

	// 创建结构体实例
	instance := reflect.New(typ).Interface()

	return instance.(T), nil
}

// MaxBy 函数用于找到切片中具有最大值的元素
// 参数:
// - slice: 要查找的切片，类型为 T 的切片
// - getValue: 一个函数，用于从切片中的元素获取一个整数值
// 返回值:
// - 返回一个整数值，表示切片中具有最大值的元素的值
func MaxBy[T any](slice []T, getValue func(T) int) int {
	if len(slice) == 0 {
		return 0
	}
	maxVal := getValue(slice[0])
	for _, item := range slice[1:] {
		val := getValue(item)
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

// MinBy 函数用于找到切片中具有最小值的元素
// 参数:
// - slice: 要查找的切片，类型为 T 的切片
// - getValue: 一个函数，用于从切片中的元素获取一个整数值
// 返回值:
// - 返回一个整数值，表示切片中具有最小值的元素的值
func MinBy[T any](slice []T, getValue func(T) int) int {
	if len(slice) == 0 {
		return 0
	}
	minVal := getValue(slice[0])
	for _, item := range slice[1:] {
		val := getValue(item)
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

// RandomInt 生成一个介于 min 和 max 之间的随机整数
// min: 最小值
// max: 最大值
// 返回值: 一个介于 min 和 max 之间的随机整数
func RandomInt(min, max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return min + r.Intn(max-min+1)
}

// SetInterval 模拟 JavaScript 中的 setInterval 函数
// interval 是执行任务的时间间隔
// fn 是要周期性执行的函数
// 返回一个通道，通过关闭该通道可以停止定时器
func SetInterval(interval time.Duration, fn func()) chan struct{} {
	stop := make(chan struct{})
	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				// 检查当前时间是否大于 21 点
				currentHour := time.Now().Hour()
				currentMinute := time.Now().Minute()
				if currentHour >= 21 && currentMinute >= 25 {
					ticker.Stop()
					close(stop)
					return
				}
				fn()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	return stop
}

// Ternary 实现类似三元操作符的功能
// condition: 条件表达式
// trueVal: 条件为真时返回的值
// falseVal: 条件为假时返回的值
func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

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
	if len(slice) == 0 {
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

// FindInSlice 用于在切片中查找满足条件的元素
// 参数说明:
//   - slice: 要搜索的切片，类型为 T 的切片
//   - predicate: 判断函数，接收类型为 T 的元素，返回是否匹配
//
// 返回值:
//   - int: 找到的元素索引，如果未找到返回 -1
//   - T: 找到的元素值，如果未找到返回零值
//   - bool: 是否找到匹配的元素
//
// 注意:
//   - 如果 slice 为 nil 或 predicate 为 nil，将返回未找到的结果
//   - 对于空切片，将直接返回未找到的结果
func FindInSlice[T any](slice []T, predicate func(T) bool) (int, T, bool) {
	// 参数验证
	if slice == nil || predicate == nil || len(slice) == 0 {
		var zero T
		return -1, zero, false
	}

	// 使用 for range 遍历切片
	for i, item := range slice {
		if predicate(item) {
			return i, item, true
		}
	}

	var zero T
	return -1, zero, false
}

// HasInSlice 用于检查切片中是否存在满足条件的元素
// 参数说明:
//   - slice: 要检查的切片，类型为 T 的切片
//   - predicate: 判断函数，接收类型为 T 的元素，返回是否匹配
//
// 返回值:
//   - bool: 是否存在匹配的元素
//
// 注意:
//   - 如果 slice 为 nil 或 predicate 为 nil，将返回 false
//   - 对于空切片，将直接返回 false
//   - 当切片长度大于 1000 时，将使用并行处理提高性能
func HasInSlice[T any](slice []T, predicate func(T) bool) bool {
	// 参数验证
	if slice == nil || predicate == nil || len(slice) == 0 {
		return false
	}

	// 对于较小的切片，使用普通遍历
	if len(slice) <= 1000 {
		for _, item := range slice {
			if predicate(item) {
				return true
			}
		}
		return false
	}

	// 对于较大的切片，使用并行处理
	const numWorkers = 4
	chunkSize := (len(slice) + numWorkers - 1) / numWorkers
	found := make(chan bool, 1)
	var wg sync.WaitGroup

	// 启动多个 goroutine 并行处理
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + chunkSize
			if end > len(slice) {
				end = len(slice)
			}
			for j := start; j < end; j++ {
				if predicate(slice[j]) {
					select {
					case found <- true:
					default:
					}
					return
				}
			}
		}(i * chunkSize)
	}

	// 等待所有 goroutine 完成
	go func() {
		wg.Wait()
		select {
		case found <- false:
		default:
		}
	}()

	return <-found
}

// ExtractProperty 函数用于从切片中提取指定属性的值
// 参数:
// - slice: 要处理的切片，类型为 T 的切片
// - getProperty: 一个函数，用于从切片中的元素获取一个值
// 返回值:
// - 返回一个包含从切片中提取的值的切片，类型为 R
func ExtractProperty[T any, R any](slice []T, getProperty func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = getProperty(item)
	}
	return result
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

// SearchFilterEx 根据指定的字段列表对搜索参数进行过滤和分类
// 参数说明:
//   - search: 原始搜索参数，包含过滤条件
//   - excludeSearchFields: 需要过滤的字段列表，支持点号分隔的字段路径（如 "user.name"）
//
// 返回值:
//   - *base_model.SearchParams: 返回一个新的搜索参数对象，其中 Filter 字段包含匹配指定字段的过滤条件
//
// 处理逻辑:
//  1. 将过滤条件分为两类：匹配指定字段的条件和不匹配的条件
//  2. 对于每个过滤条件，检查其字段是否在指定的字段列表中
//  3. 支持字段路径匹配（如 "user.name" 可以匹配 "user" 开头的字段）
//  4. 确保相同字段的过滤条件不会重复
func SearchFilterEx(search *base_model.SearchParams, excludeSearchFields ...string) *base_model.SearchParams {
	// 参数验证
	if search == nil {
		return &base_model.SearchParams{}
	}

	// 初始化结果对象
	result := &base_model.SearchParams{}

	// 使用 map 替代 garray.StrArray 进行字段去重
	fieldMap := make(map[string]struct{})
	searchFilterMap := make(map[string]base_model.FilterInfo)

	// 预处理字段列表，创建字段映射
	fieldPrefixMap := make(map[string]struct{})
	for _, field := range excludeSearchFields {
		field = gstr.ToLower(gstr.CaseSnake(field))
		fieldMap[field] = struct{}{}
		// 处理字段路径
		if parts := gstr.Split(field, "."); len(parts) > 1 {
			fieldPrefixMap[parts[0]] = struct{}{}
		}
	}

	// 存储匹配指定字段的过滤条件
	newFilter := make([]base_model.FilterInfo, 0, len(search.Filter))
	// 存储不匹配指定字段的过滤条件
	newSearchFilter := make([]base_model.FilterInfo, 0, len(search.Filter))

	for _, info := range search.Filter {
		field := gstr.ToLower(gstr.CaseSnake(info.Field))
		shouldAddToSearchFilter := true

		// 检查字段是否匹配
		if _, exists := fieldMap[field]; exists {
			// 完全匹配的情况
			newFilter = append(newFilter, info)
			shouldAddToSearchFilter = false
		} else if _, exists := fieldPrefixMap[field]; exists {
			// 字段路径匹配的情况
			shouldAddToSearchFilter = false
		}

		// 处理不匹配指定字段的过滤条件
		if shouldAddToSearchFilter {
			// 使用 map 进行去重
			if existingInfo, exists := searchFilterMap[field]; !exists || existingInfo.Value != info.Value {
				searchFilterMap[field] = info
			}
		}
	}

	// 将 map 转换回切片
	for _, info := range searchFilterMap {
		newSearchFilter = append(newSearchFilter, info)
	}

	// 更新原始搜索参数的过滤条件
	search.Filter = newSearchFilter
	// 设置结果对象的过滤条件
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
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if v != elem {
			result = append(result, v)
		}
	}
	return result
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

// AttrMake 动态地创建属性值。
//
// 该函数通过反射机制，根据传入的类型信息 T 和 TP，以及一个构建器函数 builder，
// 来生成一个特定类型的属性值。这个属性值被存储在上下文 ctx 中，使用一个与类型相关的键值。
//
// 参数:
// - ctx: 上下文，用于存储属性值。
// - key: 属性的键值，用于在上下文中检索。
// - builder: 一个函数，用于创建 TP 类型的实例。
//
// 注意: 该函数假定 ctx 中存储的属性值类型与预期匹配。
func AttrMake[T any, TP any](ctx context.Context, key string, builder func() TP) {
	// 构建一个完整的键值，包括泛型类型的名称，以确保键值的唯一性。
	key = key + "::" + reflect.ValueOf(new(T)).Type().String() + "::" + reflect.ValueOf(new(TP)).Type().String()
	// 移除键值中的 "*" 字符，这是为了键值的清晰和一致性。
	key = gstr.Replace(key, "*", "")
	// 从上下文中获取与键值相关联的属性值。
	v := ctx.Value(key)

	// 初始化一个泛型数据结构，用于存储属性键值和构建的实例。
	var data base_model.KeyValueT[string, func(data TP)]
	// 尝试将获取的值断言为目标类型。
	if v, ok := v.(base_model.KeyValueT[string, func(data TP)]); ok {
		// 如果类型断言成功，初始化 data 并使用 builder 函数创建一个实例。
		data = v
		data.Value(builder())
	} else {
		// 如果类型断言失败，输出错误信息。
		fmt.Println("Type assertion failed")
	}
}

// Throttle 函数用于限制一个操作的执行频率，确保操作不会被过于频繁地执行。
// 它接受一个函数 f 和一个间隔 interval，返回一个新函数，该新函数会确保 f 最多以 interval 为间隔执行。
// 参数:
//
//	f: 要节流的函数，即需要限制执行频率的操作。
//	interval: 两次连续执行 f 之间的最短时间间隔。
//
// 返回值:
//
//	返回一个函数，该函数在调用时会根据节流逻辑决定是否执行 f。
func Throttle(f func(), interval time.Duration) func() {
	var lastTime time.Time // 记录上一次执行 f 的时间
	var mutex sync.Mutex   // 用于确保时间检查的线程安全

	return func() {
		now := time.Now()    // 获取当前时间
		mutex.Lock()         // 上锁以确保线程安全
		defer mutex.Unlock() // 在函数退出时解锁

		// 更新 lastTime 在调用 f 之前
		lastTime = now
		// 如果是第一次调用或距离上次调用已超过interval，则执行f
		if lastTime.IsZero() || now.Sub(lastTime) >= interval {
			f() // 执行传入的函数 f
			// 确保实际调用间隔至少为 interval
			time.Sleep(interval - time.Since(lastTime))
		}
	}
}

// Throttle 示例函数，用于打印消息
//func printMessage() {
//	fmt.Println("Hello, throttled execution!")
//}
//
//func main() {
//	// 创建一个节流后的函数，每两秒最多执行一次
//	throttledPrint := Throttle(printMessage, 2*time.Second)
//
//	// 每秒尝试调用一次 throttledPrint
//	for i := 0; i < 10; i++ {
//		throttledPrint()
//		time.Sleep(1 * time.Second)
//	}
//}

/*
 Debounce 防抖函数
	该函数是一个防抖函数，用于延迟函数的执行。
	它返回一个闭包函数，接收一个函数参数f。
	在闭包函数中，通过互斥锁保证并发安全，停止之前的定时器，创建一个新的定时器，并在指定的时间间隔后执行函数f。
	这样可以避免在短时间内频繁调用函数f，达到防抖的效果。
 通俗理解：
	可以想象成一个需要"冷静一下"的机制。
	就好比你在不停地按一个按钮，但系统只在你最后一次按完并经过一段安静时间（防抖时间间隔）后，才真正去执行相应的操作。
	在这期间，不管你按得多频繁，只有最后一次按下去且等待一段时间没再按，才会触发实际动作。
*/

// Debounce 创建一个防抖函数和一个停止函数。
// 防抖函数用于在多次触发时只执行一次给定的函数，且只在停止触发后至少间隔指定时间再执行。
// 停止函数用于停止防抖函数的执行。
// 参数 interval 为防抖时间间隔，即在停止触发后多久执行给定函数。
func Debounce(interval time.Duration) (func(f func()), func()) {

	// 使用互斥锁来保证并发安全，特别是在停止和重新设定计时器时。
	var l sync.Mutex
	// timer 用于实现防抖逻辑，通过停止和重新设定来控制防抖行为。
	var timer *time.Timer

	// run 函数用于执行传入的函数f，并在f多次触发时，保证f只执行一次。
	run := func(f func()) {
		l.Lock()
		defer l.Unlock()
		// 使用lock保证d.timer更新之前一定先Stop.

		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(interval, f)
	}

	// stop 函数用于停止当前的防抖函数执行。
	stop := func() {
		l.Lock()
		defer l.Unlock()
		if timer != nil {
			timer.Stop()
			timer = nil
		}
	}

	// 返回run函数和stop函数，分别用于执行防抖逻辑和停止防抖逻辑。
	return run, stop
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
