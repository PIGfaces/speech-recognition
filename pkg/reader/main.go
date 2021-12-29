package reader

import (
	"bufio"
	"fmt"
	"os"
)

type Reader struct {
	Path       string
	ReadDataCh chan<- string
}

func (r *Reader) Read() {
	file, err := os.Open(r.Path)
	if err != nil {
		panic("open file error!" + err.Error())
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		text := fileScanner.Text()
		// fmt.Println("read url: ", text)
		r.ReadDataCh <- text
	}

	if err := fileScanner.Err(); err != nil {
		panic("read file Error" + err.Error())
	}
	fmt.Printf("[INFO] %c[42;37m读取任务完成%c[0m, 等待任务识别结束...\n", 0x1b, 0x1b)
	// 读取完成就关闭
	close(r.ReadDataCh)
}
