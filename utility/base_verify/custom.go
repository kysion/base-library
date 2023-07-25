package base_verify

import "regexp"

// IsVerify 自定义校验规则
func IsVerify(pattern string, s string) bool {
	result, _ := regexp.MatchString(pattern, s)
	return result
}
