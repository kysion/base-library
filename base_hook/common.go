package base_hook

import "context"

type CommonHookFunc func(ctx context.Context, info interface{}) error

type CommonHookInfo struct {
	Key   int
	Value CommonHookFunc
}
