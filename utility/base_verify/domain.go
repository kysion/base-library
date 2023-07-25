package base_verify

import (
	"regexp"
)

// IsDomain 是否顶级域名
func IsDomain(domain string) bool {
	result, _ := regexp.MatchString(`^([0-9a-zA-Z][0-9a-zA-Z\-]{0,62}\.)+([a-zA-Z]{0,62})$`, domain)
	return result
}
