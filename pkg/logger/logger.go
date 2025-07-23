package logger

import (
	"io"
	"log"
	"os"
)

var Logger *log.Logger
var multiWriter io.Writer

func SetLogFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	multiWriter = io.MultiWriter(os.Stdout, file)
	Logger = log.New(multiWriter, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

func toInit() {
	if Logger == nil {
		Logger = log.New(os.Stdout, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func Info(args ...interface{}) {
	toInit()
	Logger.SetPrefix("[INFO ] ")
	Logger.Println(args...)
}

func Danger(args ...interface{}) {
	toInit()
	Logger.SetPrefix("[ERROR] ")
	Logger.Println(args...)
}

func Fatal(args ...interface{}) {
	toInit()
	Logger.SetPrefix("[ERROR] ")
	Logger.Fatal(args...)
}

func Warn(args ...interface{}) {
	toInit()
	Logger.SetPrefix("[WARN ] ")
	Logger.Println(args...)
}

func Debug(args ...interface{}) {
	toInit()
	Logger.SetPrefix("[DEBUG] ")
	Logger.Println(args...)
}
