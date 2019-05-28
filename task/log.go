package task

import "log"

const LogTag = "go-timing-task: "

func logInfof(format string, args ...interface{}) {
	log.Printf(LogTag+format+"\n", args...)
}

func logErrorf(format string, args ...interface{}) {
	log.Fatalf(LogTag+format+"\n", args...)
}
