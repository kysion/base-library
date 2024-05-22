package format_utils

import (
	"fmt"
	"strconv"
)

func AnalyzeIDCard(idCard string) (year, month, day int, sex string) { // sex: 女、男
	// 提取出生日期部分
	birthStr := idCard[6:14]
	year, _ = strconv.Atoi(birthStr[:4])
	month, _ = strconv.Atoi(birthStr[4:6])
	day, _ = strconv.Atoi(birthStr[6:8])
	fmt.Printf("出生日期: %d 年 %d 月 %d 日\n", year, month, day)

	// 提取性别部分（第 17 位）
	genderCode, _ := strconv.Atoi(idCard[16:17])
	if genderCode%2 == 0 {
		fmt.Println("性别: 女")
		sex = "女"
	} else {
		fmt.Println("性别: 男")
		sex = "男"
	}

	return year, month, day, sex
}
