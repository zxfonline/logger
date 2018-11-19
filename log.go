package logger

import (
	"sync"

	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/zxfonline/chanutil"
	"github.com/zxfonline/expvar"
	"github.com/zxfonline/fileutil"
	"github.com/zxfonline/timefix"
)

var (
	logFile    *os.File
	fileLogger *log.Logger
	logchan    chan string
	base_url   = "../log/"
	PRINT      = false
	_filename  string
	stopD      chanutil.DoneChan
)

//  初始化Log 文件，不调用的话，就不会写入文件
func InitLogFile(wg *sync.WaitGroup, filename, logpath string) {
	if len(logpath) == 0 {
		logpath = base_url
	}
	base_url = logpath
	_filename = filename
	var err error
	logFile, err = fileutil.OpenFile(filepath.Join(base_url, filename+"_"+time.Now().Format("20060102")+".log"), fileutil.DefaultFileFlag, fileutil.DefaultFileMode)
	if err != nil {
		log.Panicf("open log file:%s.log error:%s", filename, err)
	}
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	// log.SetOutput(logFile)
	fileLogger = log.New(logFile, "", log.Ldate|log.Lmicroseconds)
	stopD = chanutil.NewDoneChan()
	logchan = make(chan string, 51200)
	expvar.RegistChanMonitor("chanLog", logchan)
	go writeloop(wg)
}

func GetLogger() *log.Logger {
	return fileLogger
}

// CloseLogFile 关闭日志文件
func CloseLogFile() {
	stopD.SetDone()
}

func writeloop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		logFile.Close()
		wg.Done()
	}()
	//添加跟踪信息
	// proxyTrace := trace.TraceStart("Goroutine", "Logger Start", false)
	// defer trace.TraceFinish(proxyTrace)
	now := time.Now()
	pm := time.NewTimer(time.Duration(timefix.NextMidnight(now, 1).Unix()-now.Unix()) * time.Second)
	for q := false; !q; {
		select {
		case <-stopD:
			select {
			case str := <-logchan:
				fileLogger.Println(str)
			default:
				q = true
			}
		case str := <-logchan:
			// select {
			// case <-pm.C:
			// 	if logFile1, err := fileutil.OpenFile(filepath.Join(base_url, _filename+"_"+time.Now().Format("20060102")+".log"), fileutil.DefaultFileFlag, fileutil.DefaultFileMode); err != nil {
			// 		log.Printf("[ERROR] "+"open file err:%v\n", err)
			// 	} else {
			// 		now := time.Now()
			// 		fileLogger.SetOutput(logFile1)
			// 		// log.SetOutput(logFile1)
			// 		logFile.Close()
			// 		logFile = logFile1
			// 		pm.Reset(time.Duration(timefix.NextMidnight(now, 1).Unix()-now.Unix()) * time.Second)
			// 	}
			// default:
			// }
			fileLogger.Println(str)
		case <-pm.C:
			if logFile1, err := fileutil.OpenFile(filepath.Join(base_url, _filename+"_"+time.Now().Format("20060102")+".log"), fileutil.DefaultFileFlag, fileutil.DefaultFileMode); err != nil {
				log.Printf("[ERROR] "+"open file err:%v\n", err)
			} else {
				now := time.Now()
				fileLogger.SetOutput(logFile1)
				// log.SetOutput(logFile1)
				logFile.Close()
				logFile = logFile1
				pm.Reset(time.Duration(timefix.NextMidnight(now, 1).Unix()-now.Unix()) * time.Second)
			}
		}
	}
}
