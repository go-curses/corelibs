// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package notify

import (
	"fmt"
	"io"
	"os"
)

type Level = int

const (
	Quiet Level = iota
	Info
	Debug
)

type Notifier struct {
	level  Level
	stdout io.Writer
	stderr io.Writer
}

func New(level Level) *Notifier {
	return &Notifier{
		level:  level,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (n *Notifier) Set(level Level, stdout io.Writer, stderr io.Writer) {
	n.level = level
	n.stdout = stdout
	n.stderr = stderr
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
		_, _ = fmt.Fprintf(n.stdout, format, argv...)
	}
}

func (n *Notifier) Info(format string, argv ...interface{}) {
	if n.level > Quiet {
		if len(argv) == 0 {
			argv = append(argv, format)
			format = "%s"
		}
		_, _ = fmt.Fprintf(n.stdout, format, argv...)
	}
}

func (n *Notifier) Error(format string, argv ...interface{}) {
	if n.level > Quiet {
		if len(argv) == 0 {
			argv = append(argv, format)
			format = "%s"
		}
		_, _ = fmt.Fprintf(n.stderr, format, argv...)
	}
}