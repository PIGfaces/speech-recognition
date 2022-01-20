package handler

import (
	"fmt"
	"run-api/pkg/recognition"
	"run-api/pkg/schedule"
	"run-api/pkg/writer"
	"strings"
)

type (
	TencentAsrHandlerEngin struct{}
)

func (tcentAsr *TencentAsrHandlerEngin) GetTaskSchedule(resultDir, secretId, privKeyPath string, concurrentNum int) schedule.ScheduleInter {
	var (
		secretIds []string         = strings.Split(secretId, ",")
		taskNum   int              = len(secretIds)
		writors   []*writer.Writor = make([]*writer.Writor, 0, taskNum)
		// tupuAPIClient   recognition.TupuSpeechClientIntf
		scheduleCtl    schedule.ScheduleInter
		tcentEngins    = make([]*recognition.TencentAsrHdler, 0, taskNum)
		tcentASRClient = new(recognition.TencentAsrClient).GetInterface().GetAPIClient(privKeyPath)
	)
	// 创建任务体
	for _, sid := range secretIds {
		resultCh := make(chan map[string]string, 100)
		tcentEngins = append(tcentEngins, &recognition.TencentAsrHdler{CommSpeechEngin: recognition.CommSpeechEngin{SecretId: sid, ResultDataCh: resultCh}, Client: tcentASRClient})
		writors = append(writors, &writer.Writor{Path: fmt.Sprintf("%stencent_%s_result.txt", resultDir, sid), DataCh: resultCh})
		// fmt.Printf("[INFO] %c[41;37m任务监控%c[0m : https://g.dev.tuputech.com/d/38dF4VnMz/jie-kou-chu-li-xiang-qing?orgId=1&refresh=1m&from=now-6h&to=now&var-hostname=All&var-operatorName=All&var-secretId=%s&fullscreen&panelId=4\n", 0x1b, 0x1b, sid)
	}

	tcentAsrSchedule := &schedule.TcentAsrSchedule{SecretIds: tcentEngins, CommSchedule: schedule.CommSchedule{Writors: writors, ConcurrentNum: concurrentNum}}
	tcentAsrSchedule.ScheduleInter = tcentAsrSchedule
	scheduleCtl = tcentAsrSchedule

	return scheduleCtl
}
