package logging

import "strings"

func (l *Logger) DebugSQL(q string) {
	q = strings.ReplaceAll(q, "\t", " ")
	q = strings.ReplaceAll(q, "\n", "")
	q = strings.TrimSpace(q)

	l.Debug(q)
}
