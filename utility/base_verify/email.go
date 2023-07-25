package base_verify

import "regexp"

// 识别电子邮箱
func isEmail(email string) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)

	return result
}
