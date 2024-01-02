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
	"sync"
)

// Level specifies the verbosity of a Notifier
type Level = int

const (
	// Quiet configures the Notifier to produce no output
	Quiet Level = iota
	// Error configures the Notifier to only produce error output
	Error
	// Info configures the Notifier to produce normal and error output
	Info
	// Debug configures the Notifier to produce all outputs
	Debug
)

// Builder is the interface for constructing new Notifier instances
type Builder interface {
	// SetLevel modifies the level setting
	SetLevel(level Level) Builder
	// SetOut overrides the default os.Stdout setting
	SetOut(w io.Writer) Builder
	// SetErr overrides the default os.Stderr setting
	SetErr(w io.Writer) Builder
	// Make produces the built Notifier instance
	Make() (n Notifier)
}

type Notifier interface {
	// Level returns the verbosity of this Notifier
	Level() (level Level)
	// Stdout returns the io.Writer associated with normal output
	Stdout() (w io.Writer)
	// Stderr returns the io.Writer associated with error output
	Stderr() (w io.Writer)
	// ModifyLevel modifies the configured level setting
	ModifyLevel(level Level) Notifier
	// ModifyOut modifies the configured output writer
	ModifyOut(w io.Writer) Notifier
	// ModifyErr modifies the configured error writer
	ModifyErr(w io.Writer) Notifier
	// Debug passes the arguments to fmt.Fprintf using normal output and only
	// has effect when the Level is set to Debug
	Debug(format string, argv ...interface{})
	// Info passes the arguments to fmt.Fprintf using normal output and only
	// has effect when the Level is set to Info or Debug
	Info(format string, argv ...interface{})
	// Error passes the arguments to fmt.Fprintf using normal output and only
	// has effect when the Level is set to Info or Debug
	Error(format string, argv ...interface{})
}

type CNotifier struct {
	level  Level
	stdout io.Writer
	stderr io.Writer

	sync.RWMutex
}

// New constructs a new Builder instance
func New(level Level) Builder {
	return &CNotifier{
		level:  level,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (n *CNotifier) SetLevel(level Level) Builder {
	n.level = level
	return n
}

func (n *CNotifier) SetOut(w io.Writer) Builder {
	n.stdout = w
	return n
}

func (n *CNotifier) SetErr(w io.Writer) Builder {
	n.stderr = w
	return n
}

func (n *CNotifier) Make() Notifier {
	return n
}

func (n *CNotifier) ModifyLevel(level Level) Notifier {
	n.Lock()
	defer n.Unlock()
	n.level = level
	return n
}

func (n *CNotifier) ModifyOut(stdout io.Writer) Notifier {
	n.Lock()
	defer n.Unlock()
	n.stdout = stdout
	return n
}

func (n *CNotifier) ModifyErr(stderr io.Writer) Notifier {
	n.Lock()
	defer n.Unlock()
	n.stderr = stderr
	return n
}

func (n *CNotifier) Level() (level Level) {
	n.RLock()
	defer n.RUnlock()
	return n.level
}

func (n *CNotifier) Stdout() (w io.Writer) {
	n.RLock()
	defer n.RUnlock()
	w = n.stdout
	return
}

func (n *CNotifier) Stderr() (w io.Writer) {
	n.RLock()
	defer n.RUnlock()
	w = n.stderr
	return
}

func (n *CNotifier) Debug(format string, argv ...interface{}) {
	n.Lock()
	defer n.Unlock()
	if n.stdout != nil && n.level > Info {
		if len(argv) == 0 {
			argv = append(argv, format)
			format = "%s"
		}
		_, _ = fmt.Fprintf(n.stdout, format, argv...)
	}
}

func (n *CNotifier) Info(format string, argv ...interface{}) {
	n.Lock()
	defer n.Unlock()
	if n.stdout != nil && n.level > Error {
		if len(argv) == 0 {
			argv = append(argv, format)
			format = "%s"
		}
		_, _ = fmt.Fprintf(n.stdout, format, argv...)
	}
}

func (n *CNotifier) Error(format string, argv ...interface{}) {
	n.Lock()
	defer n.Unlock()
	if n.stderr != nil && n.level > Quiet {
		if len(argv) == 0 {
			argv = append(argv, format)
			format = "%s"
		}
		_, _ = fmt.Fprintf(n.stderr, format, argv...)
	}
}
