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
)

// ListMatching is a general purpose filesystem listing function, used by all
// other List functions in this package. It accepts a custom matcher func used
// to determine whether to include the specific path or not and it returns a
// new *CListPath instance which contains all the results of the list process.
// ListMatching will return at the first error
func ListMatching(path string, includeHidden, recurse bool, matcher func(dir bool, path string) (matched bool)) (list *CListPath, err error) {
	path = Clean(path)
	list = NewListPath(path)

	// check if path is actually a file (or char device, pipe, etc.)
	if IsFile(path) {
		if IsHidden(path) {
			if includeHidden && matcher(false, path) {
				list.AddHiddenFile(path)
			}
			return
		} else if matcher(false, path) {
			list.AddNormalFile(path)
		}
		return
	} else if !IsDir(path) {
		// path is neither a file nor a directory
		err = errors.New("path not found")
		return
	}

	var entries []os.DirEntry
	if entries, err = os.ReadDir(path); err == nil {
		for _, info := range entries {
			dir := info.IsDir()
			cleaned := Clean(Join(path, info.Name()))
			if IsHidden(info.Name()) {
				if includeHidden && matcher(dir, cleaned) {
					if dir {
						list.AddHiddenDir(cleaned)
					} else {
						list.AddHiddenFile(cleaned)
					}
				}
			} else if matcher(dir, cleaned) {
				if dir {
					list.AddNormalDir(cleaned)
				} else {
					list.AddNormalFile(cleaned)
				}
			}
			if dir && recurse {
				var other *CListPath
				if other, err = ListMatching(cleaned, includeHidden, recurse, matcher); err != nil {
					return
				}
				list.AddListPath(cleaned, other)
			}
		}
	}
	return
}

// List returns a list of directories and files, sorted in natural order with
// hidden things first and directories grouped before files
func List(path string, includeHidden bool) (paths []string, err error) {
	var pl *CListPath
	if pl, err = ListMatching(path, includeHidden, false, func(dir bool, path string) (matched bool) {
		matched = true // match everything
		return
	}); err != nil {
		return
	}
	paths = pl.List()
	return
}

// ListDirs is similar to List except returning only directories
func ListDirs(path string, includeHidden bool) (paths []string, err error) {
	var pl *CListPath
	if pl, err = ListMatching(path, includeHidden, false, func(dir bool, path string) (matched bool) {
		matched = dir // match only dirs
		return
	}); err != nil {
		return
	}
	paths = pl.List()
	return
}

// ListFiles is similar to ListDirs except returning only files
func ListFiles(path string, includeHidden bool) (paths []string, err error) {
	var pl *CListPath
	if pl, err = ListMatching(path, includeHidden, false, func(dir bool, path string) (matched bool) {
		matched = !dir // match only files
		return
	}); err != nil {
		return
	}
	paths = pl.List()
	return
}

// ListAllDirs is similar to ListDirs except returning all directories found
// recursively
func ListAllDirs(path string, includeHidden bool) (paths []string, err error) {
	var pl *CListPath
	if pl, err = ListMatching(path, includeHidden, true, func(dir bool, path string) (matched bool) {
		matched = dir // match only dirs
		return
	}); err != nil {
		return
	}
	paths = pl.List()
	return
}

// ListAllFiles is similar to ListAllDirs except returning only files
func ListAllFiles(path string, includeHidden bool) (paths []string, err error) {
	var pl *CListPath
	if pl, err = ListMatching(path, includeHidden, true, func(dir bool, path string) (matched bool) {
		matched = !dir // match only files
		return
	}); err != nil {
		return
	}
	paths = pl.List()
	return
}