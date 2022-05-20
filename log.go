package main

import (
	"log"
	"os"
)

func LogFatalError(message string, err error) {
	log.Println(message, "-", err.Error())
	os.Exit(1)
}

func LogFatalApiError(message string, err ApiError) {
	LogError(message, err)
	os.Exit(1)
}

func LogError(message string, err ApiError) {
	log.Printf("%s - %q (%s / %s)", message, err.ErrorMessage, err.ErrorCode, err.ResultCode)
}
