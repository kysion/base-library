package base_verify

import "regexp"

// IsUrl 是否URL
func IsUrl(tel string) bool {
	// 匹配规则
	
	result, _ := regexp.MatchString(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`, tel)

	return result
}
