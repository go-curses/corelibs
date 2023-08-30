package notify

import (
	"fmt"
	"os"
)

type Level = int

const (
	Quiet Level = iota
	Info
	Debug
)

type Notifier struct {
	level Level
}

func New(level Level) *Notifier {
	return &Notifier{
		level: level,
	}
}

func (n *Notifier) Set(level Level) {
	n.level = level
}

func (n *Notifier) Level() (level Level) {
	return n.level
}

func (n *Notifier) Debug(format string, argv ...interface{}) {
	if n.level > Info {
		if len(argv) == 0 {
			argv = append(argv, format)
			format = "%s"
		}
		fmt.Printf(format, argv...)
	}
}

func (n *Notifier) Info(format string, argv ...interface{}) {
	if n.level > Quiet {
		if len(argv) == 0 {
			argv = append(argv, format)
			format = "%s"
		}
		fmt.Printf(format, argv...)
	}
}

func (n *Notifier) Error(format string, argv ...interface{}) {
	if n.level > Quiet {
		if len(argv) == 0 {
			argv = append(argv, format)
			format = "%s"
		}
		fmt.Fprintf(os.Stderr, format, argv...)
	}
}
