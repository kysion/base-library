package base_model

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/base-library/base_model/base_enum"
)

/*
	Hook 通信通用模型
*/

type HookHostInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type HookModel struct {
	Ctx             context.Context `json:"-"`
	Source          *HookHostInfo   `json:"source"`                 // 源
	BusinessTypeStr string          `json:"businessType" dc:"业务类型"` // 业务类型
	Data            interface{}     `json:"data" dc:"数据"`           // 数据bytes
}

// GetAddr 获取通信地址
func (m *HookHostInfo) GetAddr() string {
	return m.Host + ":" + gconv.String(m.Port)
}

func (m *HookModel) BusinessType() base_enum.HookBusinessType {
	return base_enum.Hook.BusinessType.New(m.BusinessTypeStr)
}
