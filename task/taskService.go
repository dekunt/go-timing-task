package task

import (
	"time"
	"os"
	"os/signal"
	"syscall"
)

type SyncTaskService struct {
	taskMap map[string]*SyncTask
	server  *Service
}

type SyncTask struct {
	runOnStart         bool   //服务启动时是否执行,对应环境变量"onStart{key}=1"
	timeFormatMatching string //任务执行的时间格式,对应环境变量"tf{key}=秒 分 时 日 月 年"
	workFunc           func() //任务的执行函数
	isRunning          bool
}

func (s *SyncTaskService) AddTask(key string, workFunc func()) {
	if s.taskMap == nil {
		s.taskMap = make(map[string]*SyncTask)
	}
	s.taskMap[key] = &SyncTask{
		runOnStart:         os.Getenv("onStart"+key) == "1",
		timeFormatMatching: os.Getenv("tf" + key),
		workFunc:           workFunc,
	}
}

func (s *SyncTaskService) Run() {
	if s.server != nil {
		logErrorf("SyncTaskService is already running!!!")
	}
	logInfof("start...")
	s.server = NewService(time.Second)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logInfof("shutting down...")
		s.server.Stop()
	}()
	s.server.GoDoWork(s.onStart)
	s.server.Run(s.onTick)
	logInfof("stopped")
}

func (s *SyncTaskService) Stop() {
	s.server.Stop()
}

func (s *SyncTaskService) IsStopping() bool {
	return s.server.IsStopping()
}

//服务启动时，执行部分任务
func (s *SyncTaskService) onStart() {
	for _, task := range s.taskMap {
		if task.runOnStart {
			s.taskDoWork(task)
		}
	}
}

//每秒钟触发，检查定时任务
func (s *SyncTaskService) onTick(workTime *time.Time) {
	for _, task := range s.taskMap {
		if TimeFormatMatch(task.timeFormatMatching, workTime) {
			s.taskDoWork(task)
		}
	}
}

func (s *SyncTaskService) taskDoWork(task *SyncTask) {
	if !task.isRunning {
		task.isRunning = true
		s.server.GoDoWork(func() {
			defer func() {
				task.isRunning = false
			}()
			task.workFunc()
		})
	}
}

func (s *SyncTaskService) Go(f func()) {
	s.server.GoDoWork(f)
}
