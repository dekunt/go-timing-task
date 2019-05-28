package task

import (
	"sync"
	"time"
	"runtime/debug"
)

type Service struct {
	stopSignal   chan bool
	exitSignal   chan bool
	waitGroup    *sync.WaitGroup
	workTime     *time.Time
	workDuration time.Duration
	stopping     bool
}

func NewService(duration time.Duration) *Service {
	s := &Service{
		stopSignal:   make(chan bool),
		exitSignal:   make(chan bool),
		waitGroup:    &sync.WaitGroup{},
		workDuration: duration,
		stopping:     false,
	}
	return s
}

func (s *Service) Run(worker func(t *time.Time)) {
runner:
	for {
		select {
		case <-s.stopSignal:
			logInfof("got stop signal")
			break runner
		default:
		}
		s.workLoop(worker)
	}
	s.stopping = true
	select {
	case <-s.exitSignal:
		logInfof("got exit signal")
		return
	}
}

func (s *Service) Stop() {
	close(s.stopSignal)
	s.waitGroup.Wait()
	close(s.exitSignal)
}

func (s *Service) IsStopping() bool {
	return s.stopping
}

func (s *Service) workLoop(worker func(t *time.Time)) {
	if s.workTime == nil {
		nowTime := time.Now()
		s.workTime = &nowTime
	}
	s.GoDoWork(func() {
		worker(s.workTime)
	})
	endTime := s.workTime.Add(s.workDuration)
	time.Sleep(endTime.Sub(time.Now()))
	s.workTime = &endTime
}

// 新建 GoRoutine 执行任务
func (s *Service) GoDoWork(worker func()) {
	s.waitGroup.Add(1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				logErrorf("%s: %s", p, debug.Stack())
			}
			s.waitGroup.Done()
		}()
		worker()
	}()
}
