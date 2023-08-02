package base_verify

import (
	"regexp"
)

// IsEmail 识别电子邮箱
func IsEmail(email string) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)

	return result
}
