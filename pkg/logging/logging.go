package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Fatal(msg string) {
	l.Error(msg)
	panic(msg)
}

func NewLogger(debug bool) *Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	logger := &Logger{
		slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			TimeFormat: "15:04:05",
		})),
	}

	return logger
}
