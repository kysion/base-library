package base_verify

import (
	"fmt"
	"github.com/gogf/gf/v2/text/gregex"
	"regexp"
	"strconv"
)

func IsIpv4(ipv4 string) bool {
	// 匹配规则
	result, _ := regexp.MatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, ipv4)
	return result
}

func ParseIpv4(ip string) string {
	match, err := gregex.MatchString(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`, ip)
	if err != nil || len(match) < 4 {
		return ""
	}
	return fmt.Sprintf("%s.%s.%s", match[1], match[2], match[3])
}

// ParseIpAddress 解析IP地址，返回:IP,端口.
// Eg: 192.168.1.1:80 -> 192.168.1.1, 80
func ParseIpAddress(address string) (string, int) {
	match, err := gregex.MatchString(`^(.+):(\d+)$`, address)
	if err == nil {
		i, _ := strconv.Atoi(match[2])
		return match[1], i
	}
	return "", 0
}

func IsIpv6(ipv6 string) bool {
	// 匹配规则
	result, _ := regexp.MatchString(`^([\da-fA-F]{1,4}:){7}[\da-fA-F]{1,4}$|^:((:[\da-fA-F]{1,4}){1,6}|:)$|^[\da-fA-F]{1,4}:((:[\da-fA-F]{1,4}){1,5}|:)$|^([\da-fA-F]{1,4}:){2}((:[\da-fA-F]{1,4}){1,4}|:)$|^([\da-fA-F]{1,4}:){3}((:[\da-fA-F]{1,4}){1,3}|:)$|^([\da-fA-F]{1,4}:){4}((:[\da-fA-F]{1,4}){1,2}|:)$|^([\da-fA-F]{1,4}:){5}:([\da-fA-F]{1,4})?$|^([\da-fA-F]{1,4}:){6}:$`, ipv6)
	return result
}
