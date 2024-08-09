package format_utils

import (
	"fmt"
	"strconv"
)

func AnalyzeIDCard(idCard string) (year, month, day int, sex string) { // sex: 女、男
	//  身份证号码的解析，分析出用户画像关键数据（性别、年龄、城市），
	/*
		出生日期：从身份证号码的第 7 到 14 位可以提取出出生年月日信息。
		性别：第 17 位数字，奇数表示男性，偶数表示女性。
		户籍所在地：身份证号码的前六位对应着户籍所在地的行政区划代码，但这需要对照相关的最新的行政区划编码表来确定具体地区。（362427 --> 360827）
	*/

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
