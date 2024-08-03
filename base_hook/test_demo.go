package base_hook

import "context"

type User struct {
	UserId   int
	Username string
}

type UserHookFunc func(ctx context.Context, info *User) error

type UserHookInfo struct {
	Key   int // 唯一标识: 1授权 2取消授权
	Value UserHookFunc
}
