package main

import (
	"log"
	"os"
)

func logInit() (*log.Logger, *log.Logger, error) {

	fileInfo, err1 := openLogFile("./myinfo.log")
	fileError, err2 := openLogFile("./myerror.log")

	if err1 != nil {
		return nil, nil, err1
	}
	if err1 != nil {
		return nil, nil, err2
	}

	//init the log writers with format
	infoLog := log.New(fileInfo, "[info]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	errorLog := log.New(fileError, "[error]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)

	return infoLog, errorLog, nil
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
