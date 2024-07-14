package utils

import (
	"log"
	"os"
)

var DebugLogger *log.Logger
var NoticeLogger *log.Logger
var ErrorLogger *log.Logger

func InitLogger(prod bool) {
	if prod {
		DebugLogger = log.New(os.NewFile(0, "/dev/null"), "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	} else {
		DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	}
	NoticeLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds)
}
