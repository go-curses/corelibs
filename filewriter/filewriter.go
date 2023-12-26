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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	// DefaultFileMode is the file mode used by NewFileWriter
	DefaultFileMode os.FileMode = 0644
	// DefaultOpenFlag is the os.OpenFile flag used by NewFileWriter
	DefaultOpenFlag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	// DefaultTempFile is the os.CreateTemp pattern used by NewFileWriter
	DefaultTempFile = "fw-*.tmp"
)

var _ Writer = (*CWriter)(nil)

// Writer is an io.WriteCloser that does not keep file handles open any
// longer than necessary and provides additional methods for interacting
// with the underlying file.
type Writer interface {
	io.Writer
	io.Closer

	File() (file string)
	Remove() (err error)
	ReadFile() (data []byte, err error)
	WalkFile(fn func(line string) (stop bool)) (stopped bool)

	sync.Locker
}

// CWriter implements the Writer interface.
type CWriter struct {
	file string
	flag int
	mode os.FileMode

	sync.RWMutex
}

// NewFileWriter constructs a new CWriter instance, using the given file path
// and with the default settings. If the file does not exist and the os.O_CREATE
// flag is present, the file will be created on the first write operation. If
// the file argument is empty, a temp file is created instead using
// DefaultTempFile as the os.CreateTemp pattern value.
//
// `file` is the local filesystem output destination
func NewFileWriter(file string) (w Writer, err error) {
	if file != "" {
		if file, err = filepath.Abs(file); err != nil {
			err = fmt.Errorf("error resolving absolute path to %q: %w", file, err)
			return
		}
	} else {
		var fh *os.File
		if fh, err = os.CreateTemp("", DefaultTempFile); err != nil {
			err = fmt.Errorf("error creating temp file: %w", err)
			return
		}
		file = fh.Name()
		if err = fh.Close(); err != nil {
			err = fmt.Errorf("error closing temp file: %w", err)
			return
		}
	}
	w = &CWriter{
		file: file,
		flag: DefaultOpenFlag,
		mode: DefaultFileMode,
	}
	return
}

// NewTempFileWriter is the same as NewFileWriter with the exception that the
// file is a temporary file created with os.CreateTemp and the argument is the
// pattern to use when creating it.
//
// `pattern` is the os.CreateTemp pattern argument to use
func NewTempFileWriter(pattern string) (w Writer, err error) {
	var fh *os.File
	if fh, err = os.CreateTemp("", pattern); err != nil {
		err = fmt.Errorf("error creating temp file: %w", err)
		return
	}
	pattern = fh.Name()
	if err = fh.Close(); err != nil {
		err = fmt.Errorf("error closing temp file: %w", err)
		return
	}
	w = &CWriter{
		file: pattern,
		flag: DefaultOpenFlag,
		mode: DefaultFileMode,
	}
	return
}

// SetFlag is a chainable method for setting the file flags used to open a new
// file handle each time Write is called
//
// `flag` is the os.OpenFile flags setting
func (w *CWriter) SetFlag(flag int) *CWriter {
	w.Lock()
	defer w.Unlock()
	w.flag = flag
	return w
}

// SetMode is a chainable method for setting the file mode used to open a new
// file handle each time Write is called
//
// `mode` is the file mode setting
func (w *CWriter) SetMode(mode os.FileMode) *CWriter {
	w.Lock()
	defer w.Unlock()
	w.mode = mode
	return w
}

// Write opens the log file, writes the data given and returns the bytes written
// and the error state after closing the open file handle.
func (w *CWriter) Write(p []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	var fh *os.File
	if fh, err = os.OpenFile(w.file, w.flag, w.mode); err != nil {
		return
	}
	defer func() {
		_ = fh.Close()
	}()
	n, err = fh.Write(p)
	return
}

// WriteString is a convenience wrapper around Write()
func (w *CWriter) WriteString(s string) (n int, err error) {
	n, err = w.Write([]byte(s))
	return
}

// Close is a non-operation, added to fulfill the io.WriteCloser interface
func (w *CWriter) Close() (err error) {
	return
}

// File returns the file given to NewFileWriter
func (w *CWriter) File() (file string) {
	file = w.file
	return
}

// Remove deletes the file
func (w *CWriter) Remove() (err error) {
	w.Lock()
	defer w.Unlock()
	err = os.Remove(w.file)
	return
}

// ReadFile reads and returns the entire file
func (w *CWriter) ReadFile() (data []byte, err error) {
	w.RLock()
	defer w.RUnlock()
	data, err = os.ReadFile(w.file)
	return
}

// WalkFile opens the output file and scans one line at a time, calling the
// given `fn` for each. If the `fn` returns true, the walk stops. WalkFile
// returns true if the walk was stopped.
func (w *CWriter) WalkFile(fn func(line string) (stop bool)) (stopped bool) {
	w.RLock()
	w.RUnlock()
	var fh *os.File
	var err error
	if fh, err = os.Open(w.file); err != nil {
		// logically, this should not happen
		panic(err)
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