package schedule

import (
	"fmt"
	"run-api/pkg/recognition"
	"sync"
)

type (
	TcentAsrSchedule struct {
		SecretIds []*recognition.TencentAsrHdler
		CommSchedule
	}
)

func (sdl *TcentAsrSchedule) closeAllResultCh() {
	for _, sySpch := range sdl.SecretIds {
		fmt.Printf("[INFO] %c[43;30m%s%c[0m %c[44;37m识别完成%c[0m, 等待结果保存...\n", 0x1b, sySpch.SecretId, 0x1b, 0x1b, 0x1b)
		close(sySpch.ResultDataCh)
	}
}

func (sdl *TcentAsrSchedule) concurrentRecognition(url string, concurCtl <-chan emptyStruct) {
	defer func() {
		<-concurCtl
	}()
	wg := sync.WaitGroup{}

	for _, spch := range sdl.SecretIds {
		wg.Add(1)
		go spch.Recognition(url, &wg)
	}
	wg.Wait()
}
