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
	logger zerolog.Logger // adaptee
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
	zLogger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{
		logger: zLogger,
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

func (l *Logger) msg(level zerolog.Level, msg any, args ...any) {
	event := l.logger.WithLevel(level)

	switch m := msg.(type) {
	case error:
		event = event.Err(m)

		if len(args) > 0 {
			if format, ok := args[0].(string); ok {
				if len(args) > 1 {
					event.Msgf(format, args[1:]...)
				} else {
					event.Msg(format)
				}
				return
			}
		}
		event.Msg(m.Error())

	case string:
		if len(args) == 0 {
			event.Msg(m)
		} else {
			event.Msgf(m, args...)
		}

	default:
		event.Msg(fmt.Sprintf("%v", m))
	}
}
