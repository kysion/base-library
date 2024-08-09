package sys_enum_hook

import "github.com/kysion/base-library/utility/enum"

// 业务类型：

type BusinessTypeEnum enum.IEnumCode[string]

type businessType struct {
	Default BusinessTypeEnum

	// 可拓展.....
}

var BusinessType = businessType{
	Default: enum.New[BusinessTypeEnum]("default", "default"),

	// 可拓展.....
}

func (e *businessType) New(code string, description ...string) BusinessTypeEnum {
	if code == e.Default.Code() {
		return e.Default
	}

	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}

	return enum.New[BusinessTypeEnum](code, desc)
}
