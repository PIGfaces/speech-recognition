package main

import (
	"flag"
	"fmt"
	"run-api/pkg/handler"
	"run-api/pkg/schedule"
)

var (
	resultDir     string
	taskPath      string
	secretIds     string
	privKeyPath   string
	concurrentNum int
)

func init() {
	flag.StringVar(&resultDir, "SaveResultDir", "./", "保存结果的目录, 若无则默认保存在当前目录下, 如: -SaveResultDir=./[secretId]_result.txt")
	flag.StringVar(&taskPath, "TaskFilePath", "", "任务文件路径:  保存有 url 的 txt 文件, 如: -TaskFilePath=./[secretId]_result.txt")
	flag.StringVar(&secretIds, "SecretIds", "", "需要执行任务的 secreId, 有多个则用',' 隔开; 如: -SecretIds=secretId1,secretId2")
	flag.StringVar(&privKeyPath, "PrivKeyPath", "", "secretId 对应的私钥路径, 必须有! 如: -PrivKeyPath=./rsa_private_key.pem")
	// flag.IntVar(, "PrivKeyPath", "", "secretId 对应的私钥路径, 必须有! 如: -PrivKeyPath=./rsa_private_key.pem")
	flag.IntVar(&concurrentNum, "ConcurrentNum", 10, "每个 secretId 并发数量, 默认 QPS 为 10")
}

func main() {
	flag.Parse()
	// 解析命令行
	parseArg()
	// 执行任务
	runTask()
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

// 执行任务
func runTask() {
	// 拆分需要执行的任务
	var schedule schedule.ScheduleInter = handler.RecordSpeechEngin{}.GetTaskSchedule(resultDir, secretIds, privKeyPath, concurrentNum)
	schedule.SpeechSchedule(taskPath)
	fmt.Printf("[DONE] %c[31;47m所有任务完成%c[0m\n", 0x1b, 0x1b)
}
