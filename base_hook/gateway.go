package base_hook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/kysion/base-library/base_model"
	"github.com/kysion/base-library/base_model/base_enum"
)

// GatewayHook 网关Hook
type GatewayHook func(model base_model.HookModel)

type HookMessageModel[T any] struct {
	GatewayHook      GatewayHook
	HookBusinessType base_enum.HookBusinessType
	HookFunc         T
}

// sGateway 结构体
type sGateway struct {
	GatewayHookArr []GatewayHook
}

// IGateway 接口
type IGateway interface {
	HasHookMessage() bool
	RegisterHookMessage(hook GatewayHook)
	BroadcastMessage(model base_model.HookModel)
}

var gateway = sGateway{}

func Gateway() IGateway {
	return &gateway
}

// HasHookMessage 是否有Hook消息
func (s *sGateway) HasHookMessage() bool {
	return len(s.GatewayHookArr) > 0
}

// RegisterHookMessage 注册Hook消息
func (s *sGateway) RegisterHookMessage(hook GatewayHook) {
	s.GatewayHookArr = append(s.GatewayHookArr, hook)
}

// BroadcastMessage 广播消息
func (s *sGateway) BroadcastMessage(model base_model.HookModel) {
	for _, hookItem := range s.GatewayHookArr {
		hookItem(model)
	}
}

// RegisterHookMessage 注册Hook消息,
func RegisterHookMessage[K any, F any](hook *BaseHook[K, F]) {
	gateway.RegisterHookMessage(func(model base_model.HookModel) {
		if model.BusinessType().Code() == hook.GetBusinessType().Code() {
			PublishHookMessage(context.Background(), hook, Option{})
		}
	})
}

// PublishHookMessage 发布Hook消息
func PublishHookMessage[K any, F any](ctx context.Context, hook *BaseHook[K, F], option Option, f ...func(ctx context.Context, data any, v any) error) {
	hook.Iterator(func(key K, value F) {
		_ = g.Try(ctx, func(ctx context.Context) {
			var err error
			if len(f) > 0 {
				err = f[0](ctx, option.Data, value)
				if err != nil {
					fmt.Println(err)
				}
				return
			} else {
				var vFunc interface{} = value
				err = vFunc.(CommonHookFunc)(ctx, option.Data)
			}

			if err != nil {
				fmt.Println(err)
			}
		})
	}, option)
}
