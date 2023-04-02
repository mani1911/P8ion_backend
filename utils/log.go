package utils

import (
	"p8ion/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func InitLogger() {
	config := config.GetConfig()

	if config.AppEnv == "DEV" {
		return
	}

	var (
		fileName = config.Log.FileName
		maxSize  = config.Log.MaxSize
		logLevel = config.Log.Level
	)

	if config.Log.FileName == "" {
		fileName = "./log.log"
	}

	if config.Log.MaxSize == 0 {
		maxSize = 50
	}

	if config.Log.Level == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}

	Logger = &logrus.Logger{
		Out: &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    maxSize, // megabytes
			MaxBackups: 3,
		},
		Level: level,
		Formatter: &logrus.JSONFormatter{
			// Time stamp in DD-MM-YYYY HH:MM:SS format
			TimestampFormat: "02-01-2006 15:04:05",
		},
	}

	Logger.Info("Logger started")
}

func GetLogger() *logrus.Logger {
	return Logger
}

func NewLogger(fileName string) *logrus.Logger {
	return &logrus.Logger{
		Out: &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    50, // megabytes
			MaxBackups: 3,
		},
		Level: logrus.InfoLevel,
		Formatter: &logrus.JSONFormatter{
			// Time stamp in DD-MM-YYYY HH:MM:SS format
			TimestampFormat: "02-01-2006 15:04:05",
		},
	}
}
