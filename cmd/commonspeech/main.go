package main

import (
	"flag"
	"fmt"
	"run-api/pkg/recognition"
	"run-api/pkg/schedule"
	"run-api/pkg/writer"
	"strings"
)

var (
	resultDir     string
	taskPath      string
	secretIds     string
	privKeyPath   string
	concurrentNum int
)

func Init() {
	flag.StringVar(&resultDir, "SaveResultDir", "./", "保存结果的目录, 若无则默认保存在当前目录下, 如: -SaveResultDir=./[secretId]_result.txt")
	flag.StringVar(&taskPath, "TaskFilePath", "", "任务文件路径:  保存有 url 的 txt 文件, 如: -TaskFilePath=./[secretId]_result.txt")
	flag.StringVar(&secretIds, "SecretIds", "", "需要执行任务的 secreId, 有多个则用',' 隔开; 如: -SecretIds=secretId1,secretId2")
	flag.StringVar(&privKeyPath, "PrivKeyPath", "", "secretId 对应的私钥路径, 必须有! 如: -PrivKeyPath=./rsa_private_key.pem")
	// flag.IntVar(, "PrivKeyPath", "", "secretId 对应的私钥路径, 必须有! 如: -PrivKeyPath=./rsa_private_key.pem")
	flag.IntVar(&concurrentNum, "ConcurrentNum", 10, "每个 secretId 并发数量, 默认 QPS 为 10")
}

func main() {
	Init()
	flag.Parse()
	parseArg()
	secretIds := strings.Split(secretIds, ",")
	syncSpeechs := make([]*recognition.SyncSpeech, 0, len(secretIds))
	writors := make([]*writer.Writor, 0, len(secretIds))
	syncHandler := new(recognition.SySpchClient).GetInterface()
	client := syncHandler.GetClient(privKeyPath)

	// result, statusCode, err := client.PerformWithURL("5f042c1f1bac63001e897f27", []string{"http://172.26.2.63:18888"})
	// fmt.Println(result, statusCode, err)

	for _, sid := range secretIds {
		resultCh := make(chan map[string]string, 100)
		syncSpeechs = append(syncSpeechs, &recognition.SyncSpeech{SecretId: sid, Client: client, ResultDataCh: resultCh})
		writors = append(writors, &writer.Writor{Path: fmt.Sprintf("%s%s_result.txt", resultDir, sid), DataCh: resultCh})
	}
	handler := schedule.Schedule{SecretIds: syncSpeechs, Writors: writors, ConcurrentNum: concurrentNum}
	handler.SySpeechSchedule(taskPath)
	fmt.Printf("[DONE] %c[31;47m所有任务完成%c[0m\n", 0x1b, 0x1b)
}

func parseArg() {
	if len(secretIds) == 0 {
		panic("secretId 不能为空, 请使用 -SecretIds= 参数传入, 使用 -h 参数查看帮助")
	}
	if len(privKeyPath) == 0 {
		panic("私钥未添加, 请使用 -PrivKeyPath= 参数传入, 查看详细帮助使用 -h 参数")
	}

	if len(taskPath) == 0 {
		panic("没有任务文件, 请使用 -TaskFilePath= 参数传入, 查看详细帮助使用 -h 参数")
	}
}
