package base_enum

import (
	sys_enum_captcha "github.com/kysion/base-library/base_model/base_enum/internal/captcha"
	sys_enum_hook "github.com/kysion/base-library/base_model/base_enum/internal/hook"
)

type (
	CaptchaType = sys_enum_captcha.CaptchaTypeEnum

	// HookBusinessType Hook业务类型
	HookBusinessType = sys_enum_hook.BusinessTypeEnum
)

var (
	// Captcha 验证码枚举
	Captcha = sys_enum_captcha.Captcha

	// Hook Hook枚举
	Hook = sys_enum_hook.Hook
)
