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

package path

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// IsHidden returns true if the base of the path starts with a period
func IsHidden(path string) bool {
	name := filepath.Base(path)
	return len(name) > 0 && name[0] == '.'
}

// IsBackup returns true if the path ends with a tilde
func IsBackup(path string) bool {
	last := len(path) - 1
	return last >= 0 && path[last] == '~'
}

// IsHiddenPath returns true if any of the path segments starts with a period
func IsHiddenPath(path string) (hidden bool) {
	var err error
	var absPath string
	if absPath, err = Abs(path); err != nil {
		panic(err) // untestable
	}
	for _, part := range strings.Split(absPath, "/") {
		if len(part) >= 2 && part != ".." {
			if hidden = part[0] == '.'; hidden {
				return
			}
		}
	}
	return
}

// IsPlainText returns true if the path is a file and the file's mimetype is
// of `text/plain` type
func IsPlainText(path string) (isPlain bool) {
	if IsRegularFile(path) {
		if kind, err := mimetype.DetectFile(path); err == nil {
			if isPlain = kind.Is("text/plain"); isPlain {
				return
			}
			for parent := kind.Parent(); parent != nil; parent = parent.Parent() {
				if isPlain = parent.Is("text/plain"); isPlain {
					return
				}
			}
		}
	}
	return
}

// Exists returns true if the path is present on the local filesystem (could be
// a directory or any type of file)
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		// path does *not* exist
	}
	// Schrödinger: file may or may not exist. See err for details.
	return false
}

// IsFile returns true if the path is an existing file (could be a char device,
// pipe or other unix goodness, but not a directory)
func IsFile(path string) bool {
	if info, err := os.Stat(path); err == nil {
		return info.IsDir() == false
	} else if errors.Is(err, os.ErrNotExist) {
		// path does *not* exist
	}
	// Schrödinger: file may or may not exist. See err for details.
	return false
}

// IsRegularFile returns true if the path is an existing regular file
func IsRegularFile(path string) bool {
	if info, err := os.Stat(path); err == nil {
		return info.IsDir() == false && info.Mode().IsRegular()
	} else if errors.Is(err, os.ErrNotExist) {
		// path does *not* exist
	}
	// Schrödinger: file may or may not exist. See err for details.
	return false
}

// IsDir returns true if the path is an existing directory
func IsDir(path string) bool {
	if info, err := os.Stat(path); err == nil {
		return info.IsDir()
	} else if errors.Is(err, os.ErrNotExist) {
		// path does *not* exist
	}
	// Schrödinger: file may or may not exist. See err for details.
	return false
}

// HasExt returns true if either the primary or secondary file extension
// matches the one given
func HasExt(path, extension string) (present bool) {
	if extension == "" {
		return
	} else if extension[0] == '.' {
		extension = extension[1:]
	}
	primary, secondary := ExtExt(path)
	if present = secondary != "" && secondary == extension; present {
	} else if present = primary != "" && primary == extension; present {
	}
	return
}

// HasAnyExt returns true
func HasAnyExt(path string, extensions ...string) (present bool) {
	// TODO: optimize HasAnyExt
	for _, extension := range extensions {
		if present = HasExt(path, extension); present {
			return
		}
	}
	return
}