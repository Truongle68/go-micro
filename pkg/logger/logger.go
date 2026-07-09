package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Interface interface {
	Debug(msg any, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg any, args ...any)
	Fatal(msg any, args ...any)
}

// adapter
type Logger struct {
	logger *zerolog.Logger // adaptee
}

var _ Interface = (*Logger)(nil)

func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "info":
		l = zerolog.InfoLevel
	case "warn":
		l = zerolog.WarnLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{
		logger: new(logger),
	}
}

func (l *Logger) Debug(msg any, args ...any) {
	l.msg(zerolog.DebugLevel, msg, args...)
}
func (l *Logger) Info(msg string, args ...any) {
	l.msg(zerolog.InfoLevel, msg, args...)
}
func (l *Logger) Warn(msg string, args ...any) {
	l.msg(zerolog.WarnLevel, msg, args...)
}
func (l *Logger) Error(msg any, args ...any) {
	l.msg(zerolog.ErrorLevel, msg, args...)
}
func (l *Logger) Fatal(msg any, args ...any) {
	l.msg(zerolog.FatalLevel, msg, args...)

	os.Exit(1)
}

func (l *Logger) log(level zerolog.Level, msg string, args ...any) {
	if len(args) == 0 {
		l.logger.WithLevel(level).Msg(msg)
	} else {
		l.logger.WithLevel(level).Msgf(msg, args...)
	}
}

func (l *Logger) msg(level zerolog.Level, msg any, args ...any) {
	switch m := msg.(type) {
	case error:
		l.log(level, m.Error(), args...)
	case string:
		l.log(level, m, args...)
	default:
		l.log(level, fmt.Sprintf("%s message %v has unknown type %v", level, msg, m), args...)
	}
}
