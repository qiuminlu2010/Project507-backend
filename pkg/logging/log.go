package logging

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"runtime"
// )

// type Level int

// var (
// F *os.File

// DefaultPrefix      = ""
// DefaultCallerDepth = 2

// logger1    *log.Logger
// logPrefix  = ""
// levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
// )

// const (
// 	DEBUG Level = iota
// 	INFO
// 	WARNING
// 	ERROR
// 	FATAL
// )

// func Setup() {
// 	filePath := getLogFileFullPath()
// 	F = openLogFile(filePath)

// 	logger1 = log.New(F, DefaultPrefix, log.LstdFlags)
// }

// func Debug(v ...interface{}) {
// 	setPrefix(DEBUG)
// 	logger1.Println(v...)
// }

// func Info(v ...interface{}) {
// 	setPrefix(INFO)
// 	logger1.Println(v...)
// }

// func Warn(v ...interface{}) {
// 	setPrefix(WARNING)
// 	logger1.Println(v...)
// }

// func Error(v ...interface{}) {
// 	setPrefix(ERROR)
// 	logger1.Println(v...)
// }

// func Fatal(v ...interface{}) {
// 	setPrefix(FATAL)
// 	logger1.Fatalln(v...)
// }

// func setPrefix(level Level) {
// 	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
// 	if ok {
// 		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
// 	} else {
// 		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
// 	}

// 	logger1.SetPrefix(logPrefix)
// }
