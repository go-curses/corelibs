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

package spinner

import (
	"time"

	"github.com/go-curses/cdk/lib/sync"
)

var (
	SafeSymbols = []string{
		"/", "|", "\\", "-", "/", "|", "\\", "-",
	}
	BrailleSymbols = []string{
		"⣷", "⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯",
	}
	DefaultSymbols = BrailleSymbols
)

type Callback func(symbol string)

type Spinner interface {
	String() (symbol string)
	Increment()
	Start()
	Stop()
}

type CSpinner struct {
	index    int
	symbols  []string
	ticker   *time.Ticker
	callback Callback

	sync.RWMutex
}

func NewSpinner(symbols []string, fn Callback) (s Spinner) {
	if len(symbols) == 0 {
		symbols = DefaultSymbols[:]
	}
	s = &CSpinner{
		symbols:  symbols,
		callback: fn,
	}
	return
}

func (s *CSpinner) String() (symbol string) {
	s.RLock()
	defer s.RUnlock()
	symbol = s.symbols[s.index]
	return
}

func (s *CSpinner) Increment() {
	s.Lock()
	defer s.Unlock()
	if s.index += 1; s.index >= len(s.symbols) {
		s.index = 0
	}
}

func (s *CSpinner) Start() {
	s.Lock()
	s.ticker = time.NewTicker(time.Millisecond * 250)
	s.Unlock()
	for range s.ticker.C {
		s.callback(s.String())
		s.Increment()
		if s.ticker == nil {
			break
		}
	}
}

func (s *CSpinner) Stop() {
	s.RLock()
	if s.ticker != nil {
		s.RUnlock()
		s.Lock()
		s.ticker.Stop()
		s.ticker = nil
		s.Unlock()
	} else {
		s.RUnlock()
	}
}