package logger

import (
	"os"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init (level logrus.Level) *logrus.Logger {
	log = logrus.New()
	log.SetFormatter( &logrus.TextFormatter{
		TimestampFormat: "20025-01-01 00:00:00",
	})
	// log.SetOutput(os.Stdout)
	log.SetLevel(level)
	
	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(file)
	return log
}