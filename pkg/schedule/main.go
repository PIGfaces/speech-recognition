package schedule

import (
	"run-api/pkg/reader"
	"run-api/pkg/writer"
	"sync"
	"time"
)

type (
	ScheduleInter interface {
		SpeechSchedule(taskPath string)
		closeAllResultCh()
		concurrentRecognition(url string, concurCtl <-chan emptyStruct)
	}

	CommSchedule struct {
		Writors       []*writer.Writor
		ConcurrentNum int
		ScheduleInter
	}

	emptyStruct struct{}
)

// SpeechSchedule 是执行并发任务的调度入口
func (sdl *CommSchedule) SpeechSchedule(taskPath string) {
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

func (sdl *CommSchedule) writeResult(parent *sync.WaitGroup) {
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
