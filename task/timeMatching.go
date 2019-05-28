package task

import (
	"strings"
	"strconv"
	"time"
)

var lastMatchTime *time.Time
var timeFormatSplit []string

/**
 * 检查时间表达式checkStr，是否匹配当前的时间戳
 *  例如："35 0-59/20 * * * *" 依次是：秒 分 时 日 月 年，可以匹配秒=35并且分=0/20/40的时间
 */
func TimeFormatMatch(checkStr string, workTime *time.Time) bool {
	if checkStr == "" || workTime == nil {
		return false
	}
	if lastMatchTime == nil || !lastMatchTime.Equal(*workTime) {
		timeFormat := workTime.Format("05 04 15 02 01 2006")
		timeFormatSplit = strings.Split(timeFormat, " ")
	}

	checkStrSplit := strings.Split(checkStr, " ")
	splitLen := len(checkStrSplit)
	if splitLen != len(timeFormatSplit) {
		return false
	}
	for i := 0; i < splitLen; i++ {
		if !timeSplitMatch(checkStrSplit[i], timeFormatSplit[i]) {
			return false
		}
	}
	return true
}

// 匹配: *, x, x-y, x-y/z
func timeSplitMatch(checkStr string, timeStr string) bool {
	// range匹配: *, x, x-y
	if timeSplitMatchRange(checkStr, timeStr) {
		return true
	}
	// {x-y/z} 表达式匹配
	if strings.Contains(checkStr, "/") {
		if zSplit := strings.Split(checkStr, "/"); len(zSplit) == 2 {
			if timeSplitMatchRange(zSplit[0], timeStr) {
				if z, err := strconv.Atoi(zSplit[1]); err == nil && z > 0 {
					if timeInt, _ := strconv.Atoi(timeStr); timeInt%z == 0 {
						return true
					}
				}
			}
		}
	}
	return false
}

// 匹配: *, x, x-y
func timeSplitMatchRange(checkStr string, timeStr string) bool {
	if checkStr == "*" || checkStr == timeStr {
		return true
	}
	// {X} 和 {0X} 匹配
	if len(checkStr) == 1 && len(timeStr) == 2 {
		if strings.Index(timeStr, "0") == 0 && strings.LastIndex(timeStr, checkStr) == 1 {
			return true
		}
	}
	// {x-y} 范围表达式匹配
	if strings.Contains(checkStr, "-") && !strings.Contains(checkStr, "/") {
		if xySplit := strings.Split(checkStr, "-"); len(xySplit) == 2 {
			var x, y int
			var err error
			if x, err = strconv.Atoi(xySplit[0]); err == nil {
				if y, err = strconv.Atoi(xySplit[1]); err == nil {
					timeInt, _ := strconv.Atoi(timeStr)
					if timeInt >= x && timeInt <= y {
						return true
					}
				}
			}
		}
	}
	return false
}
