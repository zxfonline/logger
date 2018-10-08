//+build debug

package logger

import (
	"fmt"
	"runtime/debug"

	"github.com/alecthomas/repr"
)

func LogDebug(format string, v ...interface{}) {
	_logFormat(Blue("[DEBUG] "), format, v...)
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
		fileLogger.Println(s) //调试环境下，直接写文件
	}
	consoleLog(s)
}

func Println(v interface{}) {
	if v == nil {
		return
	}

	// _logFormat("[Stack] ", "%v,\n%s", v, debug.Stack())
	LogDebug(repr.String(v))
}
