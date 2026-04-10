package logging

import (
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(debug bool) *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if debug {
		opts.Level = slog.LevelDebug
	}
	return &Logger{
		slog.New(slog.NewTextHandler(os.Stdout, opts)),
	}

}

func (l *Logger) DebugSQL(q string) {
	q = strings.ReplaceAll(q, "\t", "")
	q = strings.ReplaceAll(q, "\n", "")

	l.Debug(q)
}
