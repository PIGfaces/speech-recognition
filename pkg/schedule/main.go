package schedule

import (
	"fmt"
	"run-api/pkg/reader"
	"run-api/pkg/recognition"
	"run-api/pkg/writer"
	"sync"
	"time"
)

type Schedule struct {
	SecretIds     []*recognition.SyncSpeech
	Writors       []*writer.Writor
	ConcurrentNum int
}

type emptyStruct struct{}

func (sdl *Schedule) SySpeechSchedule(taskPath string) {
	urls := make(chan string, sdl.ConcurrentNum*2)
	concurCtrl := make(chan emptyStruct, sdl.ConcurrentNum)
	taskReader := reader.Reader{Path: taskPath, ReadDataCh: urls}
	wg := sync.WaitGroup{}
	go taskReader.Read()
	go sdl.writeResult(&wg)
	for url := range urls {
		concurCtrl <- emptyStruct{}
		go sdl.concurrentRecognition(url, concurCtrl)
	}

	for {
		// fmt.Println("recognition goroutine num:", len(synCh))
		if len(concurCtrl) == 0 {
			// 这里发现所有 goroutine 关闭之后会触发写退出, 最终程序完成
			sdl.closeAllResultCh()
			break
		}
		time.Sleep(time.Second * 2)
	}
	wg.Wait()
}

func (sdl *Schedule) closeAllResultCh() {
	for _, sySpch := range sdl.SecretIds {
		fmt.Printf("[INFO] %c[43;30m%s%c[0m %c[44;37m识别完成%c[0m, 等待结果保存...\n", 0x1b, sySpch.SecretId, 0x1b, 0x1b, 0x1b)
		close(sySpch.ResultDataCh)
	}
}

func (sdl *Schedule) concurrentRecognition(url string, concurCtl <-chan emptyStruct) {
	defer func() {
		<-concurCtl
	}()
	wg := sync.WaitGroup{}

	for _, syspch := range sdl.SecretIds {
		wg.Add(1)
		go syspch.Recognition(url, &wg)
	}
	wg.Wait()
}

func (sdl *Schedule) writeResult(parent *sync.WaitGroup) {
	defer func() {
		parent.Done()
	}()

	wg := sync.WaitGroup{}
	for _, writor := range sdl.Writors {
		wg.Add(1)
		go writor.Write(&wg)
	}
	wg.Wait()
}
