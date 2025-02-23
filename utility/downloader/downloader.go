package downloader

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/kysion/base-library/utility/base_funs"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// FilePart表示文件的一个分片，包含其索引、起始位置、结束位置以及数据内容。
type FilePart struct {
	Index int
	Start uint64
	End   uint64
	Data  []byte
}

// Downloader定义了文件下载器的结构，包含需要进行下载的文件的URL、文件名、输出路径、任务数量以及文件分片信息。
type Downloader struct {
	sync.Mutex  // 添加互斥锁以保证并发安全
	url         string
	fileName    string
	outputPath  string
	fileSize    uint64
	taskNum     int
	dividedFile []FilePart
}

// NewDownloader创建并返回一个新的Downloader实例。
// 如果输出路径为空，则使用当前工作目录。
// 如果输出路径不存在，则尝试创建该路径。
func NewDownloader(url string, fileName string, path string, taskNum int) (*Downloader, error) {
	if path == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			log.Println("无法获取当前工作目录：", path)
			return nil, err
		} else {
			log.Println("输出目录：", currentPath)
			path = currentPath
		}
	}

	if !gfile.Exists(path) {
		err := gfile.Mkdir(path)
		if err != nil {
			return nil, err
		}
		err = gfile.Chmod(path, 0755) // 更改为更合理的权限
		if err != nil {
			return nil, err
		}
	}

	return &Downloader{
		url:         url,
		fileName:    fileName,
		outputPath:  path,
		fileSize:    0,
		taskNum:     taskNum,
		dividedFile: make([]FilePart, taskNum), // 初始化数组
	}, nil
}

// Download方法执行文件的多线程下载。
// 它首先获取文件的总大小，然后将其分成多个部分，每个部分使用一个独立的goroutine进行下载。
// 最后，将所有下载的分片合并成一个完整的文件。
func (d *Downloader) Download() error {
	fileSize, err := d.GetFileSize()
	if err != nil {
		return err
	}
	d.fileSize = fileSize
	taskLength := fileSize / uint64(d.taskNum)
	for i := 0; i < d.taskNum; i++ {
		start := uint64(i) * taskLength
		end := start + taskLength - 1
		if i == d.taskNum-1 {
			end = d.fileSize - 1
		}
		d.dividedFile[i] = FilePart{Index: i, Start: start, End: end}
	}

	wg := sync.WaitGroup{}
	for _, t := range d.dividedFile {
		wg.Add(1)
		go func(task FilePart) {
			defer wg.Done()
			if err := d.downloadPart(task); err != nil {
				// 只在发生错误时打印日志
				log.Printf("文件下载失败 %v %v", err, task)
			}
		}(t)
	}
	wg.Wait()

	return d.merge()
}

// GetFileSize通过发送HEAD请求获取远程文件的大小。
func (d *Downloader) GetFileSize() (uint64, error) {
	var client = http.Client{}
	resp, err := client.Head(d.url)
	if err != nil {
		log.Println("获取长度失败, ", err.Error())
		return 0, err
	}
	if resp.StatusCode > 299 {
		return 0, errors.New(fmt.Sprintf("Can't process, response is %v", resp.StatusCode))
	}

	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return 0, errors.New("服务器不支持文件断点续传")
	}

	fileSize := resp.ContentLength
	log.Printf("要下载的文件大小为 %v\n", base_funs.ByteCountIEC(fileSize))
	return uint64(fileSize), nil
}

// downloadPart下载文件的一个特定分片。
// 它通过发送带有Range头的GET请求来实现。
func (d *Downloader) downloadPart(p FilePart) error {
	client := http.Client{}
	log.Printf("开始[%d]下载from:%v to:%v\n", p.Index, base_funs.ByteCountIEC(p.Start), base_funs.ByteCountIEC(p.End))
	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%v-%v", p.Start, p.End))
	resp, err := client.Do(req)

	if err != nil {
		log.Println("请求响应失败", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("确保关闭响应体", err)
			return
		}
	}(resp.Body) // 确保关闭响应体

	if resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("服务器状态码: %v", resp.StatusCode))
	}

	var data []byte
	buf := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, err := resp.Body.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		data = append(data, buf[:n]...)
	}
	if uint64(len(data)) != (p.End - p.Start + 1) {
		return errors.New("分片下载长度不正确")
	}
	d.Lock() // 加锁以保证并发安全
	p.Data = data
	d.dividedFile[p.Index] = p
	d.Unlock() // 解锁
	return nil
}

// merge方法将所有下载的文件分片合并成一个完整的文件。
func (d *Downloader) merge() error {
	path := filepath.Join(d.outputPath, d.fileName)
	log.Println("下载完毕，开始合并文件", path)
	file, err := os.Create(path)
	if err != nil {
		log.Println("文件创建失败")
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Println("文件关闭失败", err)
		}
	}(file)
	var totalSize uint64 = 0
	for i := 0; i < d.taskNum; i++ {
		writeLen, err := file.Write(d.dividedFile[i].Data)
		if err != nil {
			log.Printf("合并文件时失败, %v\n", err)
			return err
		}
		totalSize += uint64(writeLen)
	}
	if totalSize != d.fileSize {
		return errors.New("文件不完整")
	}
	return nil
}
