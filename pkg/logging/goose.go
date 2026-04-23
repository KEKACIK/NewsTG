package logging

import (
	"fmt"
	"os"
)

func (l *Logger) Printf(format string, v ...any) {
	// goose обычно пишет информационные сообщения
	l.Info(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}
