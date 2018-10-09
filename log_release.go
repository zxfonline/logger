//+build !debug

package logger

import (
	"fmt"
	"runtime/debug"
)

func LogDebug(format string, v ...interface{}) {
	//_logFormat(Blue("[DEBUG] "), format, v...)
}

func LogInfo(format string, v ...interface{}) {
	_logFormat(Green("[INFO ] "), format, v...)
}

func LogWarn(format string, v ...interface{}) {
	_logFormat(Yellow("[WARN ] "), format, v...)
}

func LogError(format string, v ...interface{}) {
	_logFormat(Red("[ERROR] "), format, v...)
	_logFormat(Red("[Stack] "), "%s", debug.Stack())
}

func LogPanic(format string, v ...interface{}) {
	_logFormat(Magenta("[PANIC] "), format, v...)
	panic("")
}

func _logFormat(prefix string, format string, v ...interface{}) {
	s := fmt.Sprintf(prefix+format, v...)
	if fileLogger != nil {
		select {
		case logchan <- s: //正式环境下，异步写文件
			if wait := len(logchan); wait > cap(logchan)/10*6 && wait%100 == 0 {
				LogWarn("logger logchan process,waitchan:%d/%d.", wait, cap(logchan))
			}
		default: //阻塞了，写入默认输出文件
			fmt.Println(s)
		}
	} else { //没有初始化，调用默认打印
		fmt.Println(s)
	}
}

func Println(v interface{}) {

}
