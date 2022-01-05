package handler

import (
	"run-api/pkg/recognition"
	"run-api/pkg/schedule"
	"run-api/pkg/writer"
	"strings"
)

type (
	SpeechHandlerEngin interface {
		GetTaskSchedule(secretId, privKeyPath string, concurrentNum int) schedule.ScheduleInter
	}

	RecordSpeechEngin struct{}

	CommSpeechEngin struct{}
)

func (RecordSpeechEngin) GetTaskSchedule(secretId, privKeyPath string, concurrentNum int) schedule.ScheduleInter {
	var (
		secretIds []string         = strings.Split(secretId, ",")
		taskNum   int              = len(secretIds)
		writors   []*writer.Writor = make([]*writer.Writor, 0, taskNum)
		// tupuAPIClient   recognition.TupuSpeechClientIntf
		scheduleCtl   schedule.ScheduleInter
		speechEngins  = make([]*recognition.SyncSpeechHdler, 0, taskNum)
		tupuAPIClient = new(recognition.SySpchClient).GetInterface().GetAPIClient(privKeyPath)
	)

	for _, sid := range secretIds {
		resultCh := make(chan map[string]string, 100)
		speechEngins = append(speechEngins, &recognition.SyncSpeechHdler{CommSpeechEngin: recognition.CommSpeechEngin{SecretId: sid, ResultDataCh: resultCh}, Client: tupuAPIClient})
	}
	scheduleCtl = &schedule.SyncSchedule{SecretIds: speechEngins, CommSchedule: schedule.CommSchedule{Writors: writors, ConcurrentNum: concurrentNum}}

	return scheduleCtl
}

// 获取同步或者异步的任务体
func (CommSpeechEngin) GetTaskSchedule(secretId, privKeyPath string, concurrentNum int) schedule.ScheduleInter {
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
	}
	scheduleCtl = &schedule.AsyncSchedule{SecretIds: speechEngins, CommSchedule: schedule.CommSchedule{Writors: writors, ConcurrentNum: concurrentNum}}

	return scheduleCtl
}
