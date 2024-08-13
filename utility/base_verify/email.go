package base_verify

import (
	"github.com/gogf/gf/v2/text/gregex"
)

// IsEmail 识别电子邮箱
func IsEmail(email string) bool {
	// 不能匹配 .  supen.huang@qq.com
	//result, _ := regexp.Compile(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`)

	result := gregex.IsMatchString(
		`^[a-zA-Z0-9_\-\.]+@[a-zA-Z0-9_\-]+(\.[a-zA-Z0-9_\-]+)+$`,
		email,
	)

	// 正则匹配邮箱
	//emailRegex, _ := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	//if emailRegex.Match([]byte(email)) {
	//	return true
	//} else {
	//	return false
	//}

	return result
}
