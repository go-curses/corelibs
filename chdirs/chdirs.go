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

package chdirs

import (
	"os"
	"sync"
)

var (
	gPushPop = struct {
		stack []string
		sync.RWMutex
	}{}
)

// Push notes the current working directory and changes directory to the given
// path, use Pop to return to the previous working directory
func Push(path string) (err error) {
	gPushPop.Lock()
	defer gPushPop.Unlock()
	var cwd string
	if cwd, err = os.Getwd(); err != nil {
		return
	} else if err = os.Chdir(path); err != nil {
		return
	}
	gPushPop.stack = append(gPushPop.stack, cwd)
	return
}

// Pop removes the last working directory from the stack and changes directory
// to it
func Pop() (err error) {
	gPushPop.Lock()
	defer gPushPop.Unlock()
	if last := len(gPushPop.stack) - 1; last >= 0 {
		path := gPushPop.stack[last]
		if last == 0 {
			gPushPop.stack = make([]string, 0)
		} else {
			gPushPop.stack = gPushPop.stack[:last]
		}
		err = os.Chdir(path)
	}
	return
}

// Stack returns the current Push stack
func Stack() (stack []string) {
	gPushPop.RLock()
	defer gPushPop.RUnlock()
	stack = gPushPop.stack[:]
	return
}