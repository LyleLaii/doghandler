package logger

import (
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Logger logger instance
var Logger = logrus.New()

// InitLogger init logger
func InitLogger() {
	currentPath, _ := os.Getwd()
	logFilePath := "log"
	logFileName := "doghandler.log"

	// 日志文件
	fileName := path.Join(currentPath, logFilePath, logFileName)

	// 设置日志级别
	Logger.SetLevel(logrus.DebugLevel)

	// 设置 rotatelogs
	logWriter, _ := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 新增 Hook
	Logger.AddHook(lfHook)
}

// LogDebug log warn infomation
func LogDebug(module string, data interface{}) {
	Logger.WithFields(
		logrus.Fields{
			"module": module,
		}).Debug(data)
}

// LogInfo log info infomation
func LogInfo(module string, data interface{}) {
	Logger.WithFields(
		logrus.Fields{
			"module": module,
		}).Info(data)
}

// LogWarn log warn infomation
func LogWarn(module string, data interface{}) {
	Logger.WithFields(
		logrus.Fields{
			"module": module,
		}).Warn(data)
}

// LogError log error ingomation
func LogError(module string, data interface{}) {
	Logger.WithFields(
		logrus.Fields{
			"module": module,
		}).Error(data)
}

// LogPanic log panic information
func LogPanic(module string, data interface{}) {
	Logger.WithFields(
		logrus.Fields{
			"module": module,
		}).Panic(data)
}
