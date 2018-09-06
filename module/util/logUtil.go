package util

import (
		"time"
	"log"
	"fmt"
)

type MyLog struct {
}

func (m *MyLog) Info(msg string) {
	m.log("info", msg)
}

func (m *MyLog) Warn(msg string) {
	m.log("warn", msg)
}

func (m *MyLog) Error(msg string) {
	m.log("error", msg)
}

func (m *MyLog) Debug(msg string) {
	m.log("debug", msg)
}

func (m *MyLog) log(level string, msg string) {
	log.Printf("[%s]-[%s]-%s\n", time.Now().Format("2006-01-02 15:04:05"),
		level, msg)
	fmt.Printf("[%s]-[%s]-%s\n", time.Now().Format("2006-01-02 15:04:05"),
		level, msg)
}
