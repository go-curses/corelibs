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

package filewriter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// DefaultFileMode is the file mode used by NewFileWriter
	DefaultFileMode os.FileMode = 0644
	// DefaultOpenFlag is the os.OpenFile flag used by NewFileWriter
	DefaultOpenFlag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	// DefaultTempFile is the os.CreateTemp pattern used by NewFileWriter
	DefaultTempFile = "filewriter-*.tmp"
)

var (
	_ FileWriter = (*cWriter)(nil)
	_ io.Writer  = (*cWriter)(nil)
	_ io.Closer  = (*cWriter)(nil)
)

// Builder is the buildable interface for constructing new FileWriter
// instances
type Builder interface {
	// SetFile specifies the file path to use
	SetFile(file string) Builder
	// UseTemp specifies the pattern to use with os.CreateTemp. If the pattern
	// includes a unix directory separator, the `pattern` argument is the base
	// name and the rest is used as the `dir` argument
	UseTemp(pattern string) Builder
	// SetMode specifies the file permissions to use when creating files
	SetMode(mode os.FileMode) Builder
	// Make initializes the settings configured with other Builder methods and
	// returns the FileWriter instance
	Make() (writer FileWriter, err error)
}

// FileWriter is an io.WriteCloser that does not keep file handles open any
// longer than necessary and provides additional methods for interacting
// with the underlying file. All operations are safe for concurrent calls
type FileWriter interface {
	// File returns the underlying file name
	File() (file string)
	// Mode returns the file permissions used when creating files
	Mode() (mode os.FileMode)
	// Remove deletes the underlying file
	Remove() (err error)

	// Write opens the log file, writes the data given and returns the number
	// of bytes written and the error state after closing the open file handle
	Write(p []byte) (n int, err error)
	// WriteString is a convenience wrapper around Write
	WriteString(s string) (n int, err error)
	// Close flags this FileWriter as being closed, blocking any further Write
	// or WriteString operations; ReadFile and WalkFile can still be used
	// until Remove is called
	Close() (err error)

	// ReadFile returns the entire file contents
	ReadFile() (data []byte, err error)
	// WalkFile opens the output file and scans one line at a time, calling the
	// given `fn` for each. If the `fn` returns true, the walk stops. WalkFile
	// returns true if the walk was stopped. WalkFile will panic on any os.Open
	// error that is not os.ErrNotExist
	WalkFile(fn func(line string) (stop bool)) (stopped bool)
}

type cWriter struct {
	file string
	temp string
	tDir string
	flag int
	mode os.FileMode

	closed bool

	sync.RWMutex
}

// New constructs a new Builder instance with the DefaultOpenFlag and
// DefaultFileMode
func New() Builder {
	return &cWriter{
		flag: DefaultOpenFlag,
		mode: DefaultFileMode,
	}
}

func (w *cWriter) SetFile(file string) Builder {
	w.file = file
	return w
}

func (w *cWriter) UseTemp(pattern string) Builder {
	if len(pattern) > 0 {
		if strings.Contains(pattern, "/") {
			path := pattern
			pattern = filepath.Base(path)
			path = strings.TrimSuffix(path, pattern)
			if last := len(path) - 1; path[last] == '/' {
				path = path[:last]
			}
			w.tDir = path
		}
	}
	w.temp = pattern
	return w
}

func (w *cWriter) SetMode(mode os.FileMode) Builder {
	w.mode = mode
	return w
}

func (w *cWriter) Make() (writer FileWriter, err error) {
	if w.temp == "" && w.file == "" {
		w.temp = DefaultTempFile
	}

	if w.temp != "" {

		var fh *os.File
		if fh, err = os.CreateTemp(w.tDir, w.temp); err != nil {
			err = fmt.Errorf("error creating temp file: %w", err)
			return
		}
		w.file = fh.Name()
		_ = fh.Close()
		if err = os.Chmod(w.file, w.mode); err != nil {
			err = fmt.Errorf("error chmod temp file %q: %w", w.file, err)
			return
		}

	} else {

		file := w.file // just for the error because it gets clobbered
		if w.file, err = filepath.Abs(w.file); err != nil {
			err = fmt.Errorf("error resolving absolute path to %q: %w", file, err)
			return
		}

	}

	writer = w
	return
}

func (w *cWriter) File() (file string) {
	w.RLock()
	defer w.RUnlock()
	file = w.file
	return
}

func (w *cWriter) Mode() (mode os.FileMode) {
	w.RLock()
	defer w.RUnlock()
	mode = w.mode
	return
}

func (w *cWriter) Remove() (err error) {
	w.Lock()
	defer w.Unlock()
	if err = os.Remove(w.file); errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	return
}

func (w *cWriter) Close() (err error) {
	w.Lock()
	defer w.Unlock()
	w.closed = true
	return
}

func (w *cWriter) Write(p []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	if w.closed {
		err = os.ErrClosed
		return
	}
	var fh *os.File
	if fh, err = os.OpenFile(w.file, w.flag, w.mode); err != nil {
		return
	}
	n, err = fh.Write(p)
	_ = fh.Close()
	return
}

func (w *cWriter) WriteString(s string) (n int, err error) {
	n, err = w.Write([]byte(s))
	return
}

func (w *cWriter) ReadFile() (data []byte, err error) {
	w.RLock()
	defer w.RUnlock()
	data, err = os.ReadFile(w.file)
	return
}

func (w *cWriter) WalkFile(fn func(line string) (stop bool)) (stopped bool) {
	w.RLock()
	w.RUnlock()
	var fh *os.File
	var err error
	if fh, err = os.Open(w.file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		panic(err) // should never happen
	}
	defer fh.Close()
	s := bufio.NewScanner(fh)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		if stopped = fn(s.Text()); stopped {
			return
		}
	}
	return
}