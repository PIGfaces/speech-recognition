package handler

import (
	"fmt"
	"run-api/pkg/recognition"
	"run-api/pkg/schedule"
	"run-api/pkg/writer"
	"strings"
)

type (
	SpeechHandlerEngin interface {
		GetTaskSchedule(resultDir, secretId, privKeyPath string, concurrentNum int) schedule.ScheduleInter
	}

	RecordSpeechEngin struct{}

	CommSpeechEngin struct{}
)

func (CommSpeechEngin) GetTaskSchedule(resultDir, secretId, privKeyPath string, concurrentNum int) schedule.ScheduleInter {
	var (
		secretIds []string         = strings.Split(secretId, ",")
		taskNum   int              = len(secretIds)
		writors   []*writer.Writor = make([]*writer.Writor, 0, taskNum)
		// tupuAPIClient   recognition.TupuSpeechClientIntf
		scheduleCtl   schedule.ScheduleInter
		speechEngins  = make([]*recognition.SyncSpeechHdler, 0, taskNum)
		tupuAPIClient = new(recognition.SySpchClient).GetInterface().GetAPIClient(privKeyPath)
	)

	// 创建任务体
	for _, sid := range secretIds {
		resultCh := make(chan map[string]string, 100)
		speechEngins = append(speechEngins, &recognition.SyncSpeechHdler{CommSpeechEngin: recognition.CommSpeechEngin{SecretId: sid, ResultDataCh: resultCh}, Client: tupuAPIClient})
		writors = append(writors, &writer.Writor{Path: fmt.Sprintf("%scommon_%s_result.txt", resultDir, sid), DataCh: resultCh})
		fmt.Printf("[INFO] %c[41;37m任务监控%c[0m : https://g.dev.tuputech.com/d/38dF4VnMz/jie-kou-chu-li-xiang-qing?orgId=1&refresh=1m&from=now-6h&to=now&var-hostname=All&var-operatorName=All&var-secretId=%s&fullscreen&panelId=4\n", 0x1b, 0x1b, sid)
	}
	syncSchedule := &schedule.SyncSchedule{SecretIds: speechEngins, CommSchedule: schedule.CommSchedule{Writors: writors, ConcurrentNum: concurrentNum}}
	syncSchedule.ScheduleInter = syncSchedule
	scheduleCtl = syncSchedule

	return scheduleCtl
}

// 获取同步或者异步的任务体
func (RecordSpeechEngin) GetTaskSchedule(resultDir, secretId, privKeyPath string, concurrentNum int) schedule.ScheduleInter {
	var (
		secretIds []string         = strings.Split(secretId, ",")
		taskNum   int              = len(secretIds)
		writors   []*writer.Writor = make([]*writer.Writor, 0, taskNum)
		// tupuAPIClient   recognition.TupuSpeechClientIntf
		scheduleCtl   schedule.ScheduleInter
		speechEngins  = make([]*recognition.AsyncSpeechHdler, 0, taskNum)
		tupuAPIClient = new(recognition.AsySpchClient).GetInterface().GetAPIClient(privKeyPath)
	)

	for _, sid := range secretIds {
		resultCh := make(chan map[string]string, 100)
		speechEngins = append(speechEngins, &recognition.AsyncSpeechHdler{CommSpeechEngin: recognition.CommSpeechEngin{SecretId: sid, ResultDataCh: resultCh}, Client: tupuAPIClient})
		writors = append(writors, &writer.Writor{Path: fmt.Sprintf("%srecord_%s_result.txt", resultDir, sid), DataCh: resultCh})
		fmt.Printf("[INFO] %c[41;37m任务监控地址%c[0m : https://g.dev.tuputech.com/d/yo0VV-IGz/yu-yin-da-wen-jian-jie-kou-xiang-qing?orgId=1&refresh=1m&fullscreen&panelId=54&var-datasource=speech-lobby-production&var-hostname=All&var-secretId=%s\n", 0x1b, 0x1b, sid)
	}
	asyncSchedule := &schedule.AsyncSchedule{SecretIds: speechEngins, CommSchedule: schedule.CommSchedule{Writors: writors, ConcurrentNum: concurrentNum}}
	asyncSchedule.ScheduleInter = asyncSchedule
	scheduleCtl = asyncSchedule

	return scheduleCtl
}
