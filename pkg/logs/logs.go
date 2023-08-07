package logs

import (
	"fmt"
	"io"
	"log"
)

type BLogger struct {
	template string
	out      io.Writer
	logger   *log.Logger
	in       chan string
}

func New(out io.Writer) *BLogger {
	b := &BLogger{
		out:      out,
		logger:   log.New(out, "", log.Ldate|log.Ltime),
		template: "%10s: %s",
		in: make(chan string),
	}
	go func() {
		for msg := range b.in {
			if b.out == nil {
				continue
			}
			b.logger.Println(msg)
		}
	}()
	return b
}

func (l *BLogger) Close() {
	close(l.in)
}

func (l *BLogger) Info(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	msg = fmt.Sprintf(l.template, "[INFO]", msg)
	l.in <- msg
}

func (l *BLogger) Warn(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	msg = fmt.Sprintf(l.template, "[WARN]", msg)
	l.in <- msg
}

func (l *BLogger) Error(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	msg = fmt.Sprintf(l.template, "[EROR]", msg)
	l.in <- msg
}
