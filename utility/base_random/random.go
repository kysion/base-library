package base_random

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateRandomString 随机生成指定长度包含数字+字母的字符串
func GenerateRandomString(length int) string {
	characters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		result += string(characters[rand.Intn(len(characters))])
	}

	return result
}

// GenerateNumberString 随机生成指定长度的数字的字符串 (10位之内)
func GenerateNumberString(length int) string {
	format := "%0" + string(length) + "v"
	code := fmt.Sprintf(format, rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	return code
}
