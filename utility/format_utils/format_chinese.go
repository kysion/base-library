package format_utils

import (
	"fmt"
	"github.com/mozillazg/go-pinyin"
)

// ChineseToPinyin 将中文转为拼音返回
func ChineseToPinyin(chineseString string) string {

	// 初始化拼音配置，默认是带声调的拼音
	pinyinConverter := pinyin.NewArgs()
	pinyinConverter.Style = pinyin.Normal

	// 将中文转换为拼音
	pinyinList := pinyin.Convert(chineseString, &pinyinConverter)

	// fmt.Println(pinyinList) //  ----》 [[ni] [hao] [shi] [jie]]

	// 打印拼音列表  ----》  [ni] [hao] [shi] [jie]
	//for _, v := range pinyinList {
	//	fmt.Print(v)
	//}
	//fmt.Println()

	// 将拼音列表转换为字符串 ----》 nihaoshijie
	var joinedPinyin string
	for _, v := range pinyinList {
		for _, word := range v {
			joinedPinyin += word + ""
		}
	}

	fmt.Println(joinedPinyin)

	return joinedPinyin
}

func test() {
	// 中文字符串
	chineseString := "你好，世界"

	pinyinStr := ChineseToPinyin(chineseString)

	fmt.Println(pinyinStr)
}
