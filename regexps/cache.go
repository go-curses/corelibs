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

package regexps

import (
	"regexp"
	"sync"
)

var _cache = newCache()

type entry struct {
	Regexp *regexp.Regexp
	Source string
	Error  error
}

type cache struct {
	data map[string]*entry

	sync.RWMutex
}

func newCache() (c *cache) {
	c = &cache{
		data: make(map[string]*entry),
	}
	return
}

func (c *cache) get(pattern string) (rx *regexp.Regexp, err error, ok bool) {
	c.RLock()
	defer c.RUnlock()
	var e *entry
	if e, ok = c.data[pattern]; ok {
		rx = e.Regexp
		err = e.Error
	}
	return
}

func (c *cache) set(pattern string, rx *regexp.Regexp, err error) {
	c.Lock()
	defer c.Unlock()
	c.data[pattern] = &entry{
		Regexp: rx,
		Source: pattern,
		Error:  err,
	}
}

func (c *cache) clear() {
	c.Lock()
	defer c.Unlock()
	c.data = make(map[string]*entry)
}

// Compile will call regexp.Compile and cache the results and any subsequent
// calls to Compile with the same patterns return the cached results
func Compile(pattern string) (rx *regexp.Regexp, err error) {
	var ok bool
	if rx, err, ok = _cache.get(pattern); ok {
		return
	}
	rx, err = regexp.Compile(pattern)
	_cache.set(pattern, rx, err)
	return
}

// ClearCache will purge all cached patterns
func ClearCache() {
	_cache.clear()
}
