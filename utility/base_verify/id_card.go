package base_verify

import "regexp"

// IsIdCard 判断是否身份证号
func IsIdCard(idcard string) bool {
	// 18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	result, _ := regexp.MatchString(`(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)`, idcard)
	return result
}
