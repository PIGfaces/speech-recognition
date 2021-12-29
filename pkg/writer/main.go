package writer

import (
	"fmt"
	"os"
	"sync"
)

type Writor struct {
	Path   string
	DataCh <-chan map[string]string
}

// Write 支持并发
func (w *Writor) Write(wg *sync.WaitGroup) {
	if w.DataCh == nil {
		panic("unknow error, mollc memory failed")
	}
	fwrite, err := os.OpenFile(w.Path, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic("create result file failed" + err.Error())
	}
	defer func() {
		fwrite.Close()
		wg.Done()
	}()
	for result := range w.DataCh {
		for url, res := range result {
			// fmt.Println("write: ", url)
			fwrite.WriteString(fmt.Sprintf("%s\t%s\n", url, res))
		}
	}
	fmt.Printf("[INFO] 保存结果到文件: %c[47;30m%s%c[0m %c[42;37m完成%c[0m...\n", 0x1b, w.Path, 0x1b, 0x1b, 0x1b)
}
