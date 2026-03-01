package logger

import (
	"log"
	"os"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func Init() {
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "INFO:  ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debug(message string) {
	debugLogger.Println(message)
}

func Info(message string) {
	infoLogger.Println(message)
}

func Error(message string) {
	errorLogger.Println(message)
}