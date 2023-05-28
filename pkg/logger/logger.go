package logger

import (
	"fmt"
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
	output := "INFO - "
	output += fmt.Sprint(args...)
	l.logger.Print(output)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	output := "INFO - "
	output += fmt.Sprintf(format, args...)
	l.logger.Print(output)
}

func (l *Logger) Warn(args ...interface{}) {
	output := "WARNING - "
	output += fmt.Sprint(args...)
	l.logger.Print(output)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	output := "WARNING - "
	output += fmt.Sprintf(format, args...)
	l.logger.Print(output)
}

func (l *Logger) Error(args ...interface{}) {
	output := "ERROR - "
	output += fmt.Sprint(args...)
	l.logger.Print(output)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	output := "ERROR - "
	output += fmt.Sprintf(format, args...)
	l.logger.Print(output)
}

func (l *Logger) Fatal(args ...interface{}) {
	output := "FATAL ERROR - "
	output += fmt.Sprint(args...)
	l.logger.Print(output)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	output := "FATAL ERROR - "
	output += fmt.Sprintf(format, args...)
	l.logger.Print(output)
}
