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
	"sort"

	"github.com/maruel/natural"

	"github.com/go-corelibs/maps"
)

type CListPath struct {
	Path   string
	Hidden struct {
		Dirs  []string
		Files []string
	}
	Normal struct {
		Dirs  []string
		Files []string
	}
	Others map[string]ListPaths
}

func NewListPath(path string) (l *CListPath) {
	l = &CListPath{
		Path:   path,
		Others: make(map[string]ListPaths),
	}
	return
}

func (l *CListPath) AddHiddenDir(paths ...string) {
	l.Hidden.Dirs = append(l.Hidden.Dirs, paths...)
}

func (l *CListPath) AddNormalDir(paths ...string) {
	l.Normal.Dirs = append(l.Normal.Dirs, paths...)
}

func (l *CListPath) AddHiddenFile(paths ...string) {
	l.Hidden.Files = append(l.Hidden.Files, paths...)
}

func (l *CListPath) AddNormalFile(paths ...string) {
	l.Normal.Files = append(l.Normal.Files, paths...)
}

func (l *CListPath) AddListPath(parent string, other *CListPath) {
	l.Others[parent] = append(l.Others[parent], other)
}

func (l *CListPath) Sort() {
	sort.Sort(natural.StringSlice(l.Hidden.Dirs))
	sort.Sort(natural.StringSlice(l.Normal.Dirs))
	sort.Sort(natural.StringSlice(l.Hidden.Files))
	sort.Sort(natural.StringSlice(l.Normal.Files))
	for _, others := range l.Others {
		others.Sort()
	}
}

func (l *CListPath) List() (sorted []string) {
	l.Sort() // everything sorted naturally

	// subdirectories nested, hidden ones first
	if dirs := append(l.Hidden.Dirs, l.Normal.Dirs...); len(dirs) > 0 {
		for _, dir := range dirs {
			sorted = append(sorted, dir)
			if other, ok := l.Others[dir]; ok {
				sorted = append(sorted, other.List()...)
			}
		}
	} else if len(l.Others) > 0 {
		for _, parent := range maps.SortedKeys(l.Others) {
			sorted = append(sorted, l.Others[parent].List()...)
		}
	}

	// files after directories, hidden ones first
	sorted = append(sorted, append(l.Hidden.Files, l.Normal.Files...)...)
	return
}
