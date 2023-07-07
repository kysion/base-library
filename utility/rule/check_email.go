package rule

import "github.com/gogf/gf/v2/text/gregex"

// IsEmail 校验是否是邮箱
func IsEmail(info string) bool {
	ok := gregex.IsMatchString(
		`^[a-zA-Z0-9_\-\.]+@[a-zA-Z0-9_\-]+(\.[a-zA-Z0-9_\-]+)+$`,
		info,
	)

	if !ok {
		return false
	}

	return true
}
