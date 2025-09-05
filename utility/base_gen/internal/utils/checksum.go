package utils

// CalculateLuhnChecksum 计算Luhn算法校验码
func CalculateLuhnChecksum(data string) int {
	sum := 0
	length := len(data)

	for i := 0; i < length; i++ {
		// 将字符转换为数字
		if data[i] < '0' || data[i] > '9' {
			return -1 // 非法字符
		}

		digit := int(data[i] - '0')

		// 偶数位（从右向左数）加倍
		if (length-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	// 计算校验码
	return (10 - (sum % 10)) % 10
}
