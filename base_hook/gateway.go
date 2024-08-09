package base_hook

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/base-library/base_model"
	"github.com/kysion/base-library/utility/enum"
	"github.com/kysion/base-library/utility/kmap"
	"reflect"
	"strings"
)

// GatewayHook 网关Hook
type GatewayHook func(model base_model.HookModel)

// sGateway 结构体
type sGateway struct {
	GatewayHookMap kmap.HashMap[string, GatewayHook]
}

// IGateway 接口
type IGateway interface {
	HasHookMessage() bool
	BroadcastMessage(model base_model.HookModel)
}

var gateway = sGateway{}

func Gateway() IGateway {
	return &gateway
}

// HasHookMessage 是否有Hook消息
func (s *sGateway) HasHookMessage() bool {
	return s.GatewayHookMap.Size() > 0
}

// BroadcastMessage 处理广播消息
func (s *sGateway) BroadcastMessage(model base_model.HookModel) {
	s.GatewayHookMap.Iterator(func(key string, hookFunc GatewayHook) bool {
		hookFunc(model)
		return true
	})
}

// RegisterHookMessage 注册Hook消息,
func RegisterHookMessage[K any, F any](hook *BaseHook[K, F]) bool {
	if gateway.GatewayHookMap.Contains(hook.GetBusinessType().Code()) {
		return false
	}

	//【预注册】相当于是将业务层的hookFunc预先封装成一个函数，并且将其存储在gateway.GatewayHookMap中。然后在后续的有发布消息时，就可以通过这个函数来调用了。
	var gatewayHook GatewayHook = func(model base_model.HookModel) {
		gateway.GatewayHookMap.Iterator(func(key string, hookFunc GatewayHook) bool {
			if model.BusinessType().Code() == key {
				var option = Option{}
				_ = gconv.Struct(model.Data, &option)
				// 如果是网络消息，则不进行调用, 因这里是网络消息，所以强制改成false，防止循环调用
				option.NetMessage = false
				PublishHookMessage(context.Background(), hook, option)
			}
			return true
		})
	}

	gateway.GatewayHookMap.Set(hook.GetBusinessType().Code(), gatewayHook)
	return true
}

// PublishHookMessage 发布Hook消息
func PublishHookMessage[K any, F any](ctx context.Context, hook *BaseHook[K, F], option Option) []error {
	var errArr []error
	var dataKind = reflect.TypeOf(option.Data)
	var srcDataArr []interface{}

	// 如果option.Data是数组类型，则将数据断言为数组类型
	if dataKind.Kind() == reflect.Array || dataKind.Kind() == reflect.Slice {
		srcDataArr = option.Data.([]interface{})
	} else { // 如果option.Data不是数组类型，则强制构建一个数组，进行赋值
		srcDataArr = []interface{}{option.Data}
	}

	// 如果是网络消息，则不进行调用
	hook.Iterator(func(key K, value F) {
		// 过滤掉不匹配的hook订阅
		if option.NetMessage == false && option.HookTypeStr != reflect.TypeOf(value).String() {
			return
		}

		var err error
		err = g.Try(ctx, func(ctx context.Context) {
			of := reflect.ValueOf(value) // 获取回调函数的反射对象
			inNum := of.Type().NumIn()

			// 如果回调函数的参数个数与数据的元素个数不一致，则返回错误
			if inNum != len(srcDataArr)+1 {
				g.Log().Error(ctx, err)
				errArr = append(errArr, fmt.Errorf("回调函数的参数个数与数据的元素个数不一致"))
				return
			}

			// 如果是数组类型，则将数据转换为数组类型
			inArr := make([]reflect.Value, 0)
			inArr = append(inArr, reflect.ValueOf(ctx))
			// 将消息数据转换为hook回调函数参数对应的类型
			for i := 1; i < inNum; i++ {
				in := of.Type().In(i) // 获取回调函数的第二个参数的类型
				_ = g.Try(ctx, func(ctx context.Context) {

					var targetData interface{}

					if strings.HasSuffix(strings.Split(in.String(), ".")[0], "enum") { // 枚举类型的包名，后缀必须是 enum
						enumMap := srcDataArr[i-1].(map[string]interface{})
						desc := gconv.String(enumMap["description"])

						switch reflect.ValueOf(enumMap["code"]).Kind() { // 枚举类型对象的code类型，通常都是 int 和 string，如果是其他类型，需要在此处增加case 的判断逻辑

						case reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

							if gconv.Int64(enumMap["code"]) < 99999999 {
								code := gconv.Int(enumMap["code"])
								targetData = enum.New[enum.IEnumCode[int]](code, desc)
							} else {
								code := gconv.Int64(enumMap["code"])
								targetData = enum.New[enum.IEnumCode[int64]](code, desc)
							}
						case reflect.String:
							code := gconv.String(enumMap["code"])
							targetData = enum.New[enum.IEnumCode[string]](code, desc)
						}
					} else {
						// 根据函数参数类型，构建参数对象
						targetData = reflect.New(in.Elem()).Interface()

						// 如果是数组类型，则将数据转换为回调函数参数对应数据类型
						switch reflect.ValueOf(targetData).Kind() {
						// 如果是结构体类型，则将数据转换为回调函数参数对应数据类型
						case reflect.Struct, reflect.Ptr:
							_ = gconv.Struct(srcDataArr[i-1], targetData)
						// 如果是常量则直接赋值
						case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String, reflect.Array, reflect.Slice, reflect.Map:
							targetData = srcDataArr[i-1]
						}
					}

					// 将数据追加到参数数组中
					inArr = append(inArr, reflect.ValueOf(targetData))
				})
			}

			// 调用HOOK回调函数 hookFunc
			retArr := of.Call(inArr)

			// 如果回调函数返回了错误，则返回错误
			if len(retArr) > 0 {
				ret := retArr[0].Interface()
				if ret != nil {
					errors.As(err, &ret)
				}
			}
		})
		if err != nil {
			g.Log().Error(ctx, err)
			errArr = append(errArr, err)
		}
	}, Option{
		Data:        option.Data,
		NetMessage:  option.NetMessage,
		HookTypeStr: option.HookTypeStr,
	})
	return errArr
}
