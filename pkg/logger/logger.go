package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type logger struct {
	logger *log.Logger
	level  Level
}

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

func New(level Level) Logger {
	return &logger{
		logger: log.New(os.Stdout, "", 0),
		level:  level,
	}
}

func (l *logger) Debug(args ...interface{}) {
	if l.level <= Debug {
		l.log("DEBUG", fmt.Sprint(args...))
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.level <= Debug {
		l.log("DEBUG", fmt.Sprintf(format, args...))
	}
}

func (l *logger) Info(args ...interface{}) {
	if l.level <= Info {
		l.log("INFO", fmt.Sprint(args...))
	}
}

func (l *logger) Infof(format string, args ...interface{}) {
	if l.level <= Info {
		l.log("INFO", fmt.Sprintf(format, args...))
	}
}

func (l *logger) Warn(args ...interface{}) {
	if l.level <= Warn {
		l.log("WARN", fmt.Sprint(args...))
	}
}

func (l *logger) Warnf(format string, args ...interface{}) {
	if l.level <= Warn {
		l.log("WARN", fmt.Sprintf(format, args...))
	}
}

func (l *logger) Error(args ...interface{}) {
	if l.level <= Error {
		l.log("ERROR", fmt.Sprint(args...))
	}
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if l.level <= Error {
		l.log("ERROR", fmt.Sprintf(format, args...))
	}
}

func (l *logger) Fatal(args ...interface{}) {
	if l.level <= Fatal {
		l.log("FATAL", fmt.Sprint(args...))
		os.Exit(1)
	}
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	if l.level <= Fatal {
		l.log("FATAL", fmt.Sprintf(format, args...))
		os.Exit(1)
	}
}

func (l *logger) log(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.logger.Printf("[%s] %s: %s", timestamp, level, message)
}
