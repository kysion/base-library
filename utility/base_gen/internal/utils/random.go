package utils

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// GenerateRandomDigits 生成指定长度的随机数字字符串
func GenerateRandomDigits(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}

	// 使用crypto/rand生成安全随机数
	b := make([]byte, int(math.Ceil(float64(length)/8*5))) // base32编码每5字节生成8字符
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// 转换为base32编码并移除可能的填充字符
	encoded := base32.StdEncoding.EncodeToString(b)
	encoded = strings.TrimRight(encoded, "=")

	// 只保留数字
	var result strings.Builder
	for _, r := range encoded {
		if r >= '0' && r <= '9' {
			result.WriteRune(r)
		} else {
			// 转换字母为数字
			result.WriteRune('0' + (r-'A')%10)
		}
		if result.Len() == length {
			break
		}
	}

	// 如果长度不足，补充随机数字
	for result.Len() < length {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result.WriteRune(rune('0' + n.Int64()))
	}

	return result.String(), nil
}

// GenerateUUID 生成UUID
func GenerateUUID() (string, error) {
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		return "", err
	}

	// 设置UUID版本为4（随机生成）
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// 设置UUID变体为RFC 4122
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	// 转换为字符串
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
