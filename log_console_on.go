//+build console

package logger

import (
	"log"
)

func consoleLog(format string, v ...interface{}) {
	log.Printf(format, v...)
}
