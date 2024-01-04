// Copyright (c) 2024  The Go-Curses Authors
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
	"strings"
)

// FindFileRelativeToPath looks for the named file within the directory path
// given and if not found, walks up the parent directories, checking each of
// them and returning the first named file found (as an absolute path)
//
// The intent is something similar to how the `git` command knows it's in a
// repository even though the command may be run from within a subdirectory
// of the repository
func FindFileRelativeToPath(name, path string) (file string) {
	if abs, err := Abs(path); err == nil {
		if absName := abs + "/" + name; IsFile(absName) {
			file = absName
			return
		}
		parts := strings.Split(abs, "/")
		parts = parts[1:]
		pl := len(parts)
		for i := pl - 1; i >= 0; i-- {
			combined := "/" + strings.Join(parts[0:i], "/") + "/" + name
			if IsFile(combined) {
				file = combined
				return
			}
		}
	}
	return
}

// FindFileRelativeToPwd is a convenience wrapper for FindFileRelativeToPath,
// specifying the given name and a path of "."
func FindFileRelativeToPwd(name string) (file string) {
	file = FindFileRelativeToPath(name, ".")
	return
}