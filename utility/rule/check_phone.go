package rule

import "github.com/gogf/gf/v2/text/gregex"

// IsPhone 校验是否是手机号
func IsPhone(info string) bool {
	ok := gregex.IsMatchString(
		`^13[\d]{9}$|^14[5,7]{1}\d{8}$|^15[^4]{1}\d{8}$|^16[\d]{9}$|^17[0,2,3,5,6,7,8]{1}\d{8}$|^18[\d]{9}$|^19[\d]{9}$`,
		info,
	)
	if !ok {
		return false
	}

	return true
}
