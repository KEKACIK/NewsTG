package logging

import "strings"

func (l *Logger) DebugSQL(q string, args ...any) {
	q = strings.ReplaceAll(q, "\t", " ")
	q = strings.ReplaceAll(q, "\n", "")
	q = strings.TrimSpace(q)

	l.Debug(q, args...)
}
