package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

/**
 * @Description: 初始化日志
 * @param fileName 日志文件路径
 * @param level 日志级别
 * 				日志级别判断规则是：‌只有日志消息的级别 ≥ 当前设定级别时，才会被输出。‌
 * 				日志级别层级关系（从低到高）：Trace < Debug < Info < Warn < Error < Fatal < Panic
 *				info为例（Trace、Debug—— ‌不会打印‌（因为低于 Info）；Info、Warn、Error、Fatal、Panic —— ‌都会打印‌）
 * @param out 日志输出方式 console/file/file+console
 * @return *logrus.Logger
 */
func InitLog(fileName, level, out string) *logrus.Logger {
	Log := logrus.New()
	Log.SetFormatter(new(LogFormatter))
	Log.SetReportCaller(true)

	// 根据配置文件设置日志级别
	if level != "" {
		l, err := logrus.ParseLevel(level)
		if err != nil {
			panic(any(fmt.Sprintf("日志级别不存在: %s", level)))
		}
		Log.SetLevel(l)
	} else {
		Log.SetLevel(logrus.DebugLevel)
	}
	var file *os.File
	if out == "file" || out == "file+console" {
		if fileName == "" {
			fileName = "./logs/info.log"
		}
		// 创建目录
		dir := filepath.Dir(fileName)
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(any(fmt.Sprintf("创建日志目录失败: %s", err.Error())))
		}
		// 创建日志文件
		var err error
		file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend|0666)
		if err != nil {
			panic(any(fmt.Sprintf("创建日志文件失败: %s", err.Error())))
		}
	}
	if out == "console" {
		Log.Out = os.Stdout
	} else if out == "file" {
		Log.Out = file
	} else if out == "file+console" {
		// 将 stdout 和文件组合成一个多路 Writer
		mw := io.MultiWriter(os.Stdout, file)
		Log.Out = mw
	}

	return Log
}

type LogFormatter struct{}

func (l *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05.000")
	level := entry.Level
	logMsg := fmt.Sprintf("%s [%s]", timestamp, strings.ToUpper(level.String()))
	// 如果存在调用信息，且为error级别以上记录文件及行号
	if caller := entry.Caller; caller != nil {
		var fp string
		// 全路径切割，只获取项目相关路径，
		// 即/Users/hml/Desktop/project/go/pandax/test.go只获取/test.go
		ps := strings.Split(caller.File, "pandax/")
		if len(ps) >= 2 {
			fp = ps[1]
		} else {
			fp = ps[0]
		}
		logMsg = logMsg + fmt.Sprintf(" [%s:%d]", fp, caller.Line)
	}
	for k, v := range entry.Data {
		logMsg = logMsg + fmt.Sprintf(" [%s=%v]", k, v)
	}
	logMsg = logMsg + fmt.Sprintf(" : %s\n", entry.Message)
	return []byte(logMsg), nil
}
