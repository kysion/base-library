package base_hook

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/base-library/base_model"
	"github.com/kysion/base-library/utility/kmap"
	"reflect"
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
	// 如果是网络消息，则不进行调用
	hook.Iterator(func(key K, value F) {
		var err error
		err = g.Try(ctx, func(ctx context.Context) {
			of := reflect.ValueOf(value) // 获取回调函数的反射对象
			in := of.Type().In(1)        // 获取回调函数的第二个参数的类型
			var data interface{}
			if in.Kind() == reflect.Ptr { // 如果第二个参数是指针类型，则创建一个指针类型的对象
				// 将数据转换为指针类型
				data = reflect.New(in.Elem()).Interface()
				_ = gconv.Struct(option.Data, data)
			}
			// 调用回调函数 hookFunc
			retArr := of.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(data)})
			if len(retArr) > 0 {
				ret := retArr[0].Interface()
				if ret != nil {
					err = ret.(error)
				}
			}
		})
		if err != nil {
			g.Log().Error(ctx, err)
			errArr = append(errArr, err)
		}
	}, Option{
		Data:       option.Data,
		NetMessage: option.NetMessage,
	})
	return errArr
}
