package task

import (
	"testing"
	"os"
	"time"
)

func TestSyncTaskService_Run(t *testing.T) {

	server := SyncTaskService{}

	//定义任务
	key := "Task1"
	task1Func := func() {
		logInfof("task1...")
		time.Sleep(3 * time.Second) //用于验证优雅退出，任务正在执行时，会等该任务执行完再退出整个程序
		logInfof("task1 done.")
	}

	//环境变量配置
	os.Setenv("tf"+key, "0-59/2 * * * * *") //通过环境变量tf{KEY}配置：每2秒执行一次
	//os.Setenv("onStart"+key, "1")           //通过环境变量onStart{KEY}配置：启动时执行

	//添加任务
	server.AddTask(key, task1Func)

	go func() {
		//5秒后退出
		for i := 0; i < 10; i++ {
			logInfof("main-work ...")
			time.Sleep(500 * time.Millisecond)
		}
		server.Stop()
	}()
	//启动服务
	server.Run()

	/**
	执行结果：
	2019/06/01 15:00:00 [go-timing-task] start...
	2019/06/01 15:00:00 [go-timing-task] main-work ...
	2019/06/01 15:00:00 [go-timing-task] task1...
	2019/06/01 15:00:01 [go-timing-task] main-work ...
	2019/06/01 15:00:01 [go-timing-task] main-work ...
	2019/06/01 15:00:02 [go-timing-task] main-work ...
	2019/06/01 15:00:02 [go-timing-task] main-work ...
	2019/06/01 15:00:03 [go-timing-task] main-work ...
	2019/06/01 15:00:03 [go-timing-task] task1 done.
	2019/06/01 15:00:03 [go-timing-task] main-work ...
	2019/06/01 15:00:04 [go-timing-task] main-work ...
	2019/06/01 15:00:04 [go-timing-task] task1...
	2019/06/01 15:00:04 [go-timing-task] main-work ...
	2019/06/01 15:00:05 [go-timing-task] main-work ...
	2019/06/01 15:00:06 [go-timing-task] got stop signal
	2019/06/01 15:00:07 [go-timing-task] task1 done.
	2019/06/01 15:00:07 [go-timing-task] got exit signal
	2019/06/01 15:00:07 [go-timing-task] stopped
	 */
}
