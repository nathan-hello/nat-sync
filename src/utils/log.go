package utils

import (
	"log"
	"os"
)

var DebugLogger *log.Logger
var NoticeLogger *log.Logger
var ErrorLogger *log.Logger
var UserLogger *log.Logger

func InitLogger() {
	format := log.Ltime | log.Lshortfile
	DebugLogger = log.New(os.Stdout, "DEBUG: ", format)
	NoticeLogger = log.New(os.Stdout, "INFO: ", format)
	UserLogger = log.New(os.Stdout, "", format)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", format)
}
