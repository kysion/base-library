package file_utils

import (
	"io"
	"log"
	"os"
)

// GetFile 根据文件路径获取文件内容并返回字符串格式内容
func GetFile(path string) []byte {
	file, err := os.Open(path)

	var data []byte
	buf := make([]byte, 1024)
	for {
		// 将文件中读取的byte存储到buf中
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}
		// 将读取到的结果追加到data切片中
		data = append(data, buf[:n]...)
	}

	if err != nil {
		log.Fatal(err)
	}

	// 将data切片转为字符串即使文件内容
	return data
}
