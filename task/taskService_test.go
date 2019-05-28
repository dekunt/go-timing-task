package task

import (
	"testing"
	"os"
	"time"
)

func TestSyncTaskService_Run(t *testing.T) {

	server := SyncTaskService{}

	keyTick := "Tick"
	os.Setenv("tf"+keyTick, "0-59/2 * * * * *") //每2秒打印一个心跳日志
	server.AddTask(keyTick, func() {
		logInfof("Tick...")
		time.Sleep(3 * time.Second) //用于验证优雅退出，任务正在执行时，会等该任务执行完再退出整个程序
		logInfof("Tick... done.")
	})

	go func() {
		//5秒后退出
		time.Sleep(5 * time.Second)
		server.server.Stop()
	}()
	//启动服务
	server.Run()

	/**
	执行结果：
	xx:xx:12 go-timing-task: start...
	xx:xx:12 go-timing-task: Tick...
	xx:xx:15 go-timing-task: Tick... done.
	xx:xx:16 go-timing-task: Tick...
	xx:xx:18 go-timing-task: got stop signal
	xx:xx:19 go-timing-task: Tick... done.
	xx:xx:19 go-timing-task: got exit signal
	xx:xx:19 go-timing-task: stopped
	 */
}
