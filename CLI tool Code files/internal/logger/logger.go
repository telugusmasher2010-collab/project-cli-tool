package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	level  Level
	output io.Writer
	file   *os.File
}

var Default = New(LevelInfo, os.Stderr, "")

func New(level Level, output io.Writer, logPath string) *Logger {
	l := &Logger{level: level, output: output}
	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			l.file = f
		}
	}
	return l
}

func SetLevel(level Level) {
	Default.level = level
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("%s [%s] %s\n", time.Now().Format(time.RFC3339), level.String(), msg)
	fmt.Fprint(l.output, line)
	if l.file != nil {
		fmt.Fprint(l.file, line)
	}
}

func Debug(format string, args ...interface{}) { Default.log(LevelDebug, format, args...) }
func Info(format string, args ...interface{})  { Default.log(LevelInfo, format, args...) }
func Warn(format string, args ...interface{})  { Default.log(LevelWarn, format, args...) }
func Error(format string, args ...interface{}) { Default.log(LevelError, format, args...) }

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
