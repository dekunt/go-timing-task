package task

import "log"

const logTag = "[go-timing-task] "

func logInfof(format string, args ...interface{}) {
	log.Printf(logTag+format+"\n", args...)
}

func logErrorf(format string, args ...interface{}) {
	log.Printf("[Error]"+logTag+format+"\n", args...)
}
