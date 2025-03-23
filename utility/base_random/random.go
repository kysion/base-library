package base_random

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// GenerateRandomString 随机生成指定长度包含数字+字母的字符串
func GenerateRandomString(length int) string {
	characters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""

	// 创建本地随机生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		idx := r.Intn(len(characters))
		result += string(characters[idx : idx+1])
	}

	return result
}

// GenerateNumberString 随机生成指定长度的数字的字符串 (10位之内)
func GenerateNumberString(length int) string {
	format := "%0" + strconv.Itoa(length) + "v"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf(format, r.Int31n(1000000))
	return code
}
