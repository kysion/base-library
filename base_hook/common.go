package base_hook

import "context"

// DefaultHookFunc 公共Hook函数【用于测试，业务层应该有自己的Hook定义】
type DefaultHookFunc func(ctx context.Context, info interface{}) error

type DefaultHookInfo struct {
	Key   int
	Value DefaultHookFunc
}

// UserHookFunc 用户Hook函数【用于测试，业务层应该有自己的Hook定义】
type UserHookFunc func(ctx context.Context, info *User) error

type UserHookInfo struct {
	Key   int // 唯一标识: 1授权 2取消授权
	Value UserHookFunc
}
type User struct {
	UserId   int
	Username string
}
