package basic

import "log"

type Logger interface {
    Debug(p ...interface{})
    Info(p ...interface{})
    Warn(p ...interface{})
    Error(p ...interface{})
}

type DefaultLogger struct {
}

func (dl DefaultLogger) Debug(p ...interface{}) {
    log.Println(p...)
}

func (dl DefaultLogger) Info(p ...interface{}) {
    log.Println(p...)
}

func (dl DefaultLogger) Warn(p ...interface{}) {
    log.Println(p...)
}

func (dl DefaultLogger) Error(p ...interface{}) {
    log.Println(p...)
}

type NoOpLogger struct {
}

func (dl NoOpLogger) Debug(p ...interface{}) {
}

func (dl NoOpLogger) Info(p ...interface{}) {
}

func (dl NoOpLogger) Warn(p ...interface{}) {
}

func (dl NoOpLogger) Error(p ...interface{}) {
}
