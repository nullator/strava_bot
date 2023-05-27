package logger

import (
	"log"
)

type LoggerInterface interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Logger struct {
	name   string
	logger *log.Logger
}

var _ LoggerInterface = (*Logger)(nil)

func New(name string, l *log.Logger) *Logger {
	return &Logger{name: name, logger: l}
}

func (l *Logger) Info(args ...interface{}) {
	a := make([]interface{}, 0)
	a = append(a, "INFO - ")
	a = append(a, args...)
	l.logger.Print(a...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	a := make([]interface{}, 0)
	a = append(a, "INFO - ")
	a = append(a, args...)
	l.logger.Printf(format, a...)
}

func (l *Logger) Warn(args ...interface{}) {
	a := make([]interface{}, 0)
	a = append(a, "WARNING - ")
	a = append(a, args...)
	l.logger.Print(a...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	a := make([]interface{}, 0)
	a = append(a, "WARNING - ")
	a = append(a, args...)
	l.logger.Printf(format, a...)
}

func (l *Logger) Error(args ...interface{}) {
	a := make([]interface{}, 0)
	a = append(a, "ERROR - ")
	a = append(a, args...)
	l.logger.Print(a...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	a := make([]interface{}, 0)
	a = append(a, "ERROR - ")
	a = append(a, args...)
	l.logger.Printf(format, a...)
}

func (l *Logger) Fatal(args ...interface{}) {
	a := make([]interface{}, 0)
	a = append(a, "FATAL ERROR - ")
	a = append(a, args...)
	l.logger.Fatal(a...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	a := make([]interface{}, 0)
	a = append(a, "FATAL ERROR - ")
	a = append(a, args...)
	l.logger.Fatalf(format, a...)
}
