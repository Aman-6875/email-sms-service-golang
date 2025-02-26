package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	Log.SetLevel(logrus.DebugLevel)

	Log.SetOutput(os.Stdout)

	Log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true, // Set to false for production
	})

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Info("Failed to log to file, using default stdout")
	}
	Log.Info("Logger initialized successfully!")
}

func GetLogger() *logrus.Logger {
	return Log
}
