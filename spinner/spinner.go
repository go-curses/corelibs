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
	"sync"
	"time"

	"github.com/go-curses/cdk/lib/paint"
)

var (
	DefaultSpinnerRunes, _ = paint.GetSpinners(paint.SevenDotSpinner)
	DefaultSymbols         = DefaultSpinnerRunes.Strings()
)

type Callback func(symbol string)

type Spinner interface {
	String() (symbol string)
	Increment()
	Start()
	StartWith(interval time.Duration)
	Stop()
}

type cSpinner struct {
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
	s = &cSpinner{
		symbols:  symbols,
		callback: fn,
	}
	return
}

func (s *cSpinner) String() (symbol string) {
	s.RLock()
	defer s.RUnlock()
	symbol = s.symbols[s.index]
	return
}

func (s *cSpinner) Increment() {
	s.Lock()
	defer s.Unlock()
	if s.index += 1; s.index >= len(s.symbols) {
		s.index = 0
	}
}

func (s *cSpinner) Start() {
	s.StartWith(time.Millisecond * 250)
}

func (s *cSpinner) StartWith(interval time.Duration) {
	s.Lock()
	s.ticker = time.NewTicker(interval)
	s.Unlock()
	for range s.ticker.C {
		s.callback(s.String())
		s.Increment()
		if s.ticker == nil {
			break
		}
	}
}

func (s *cSpinner) Stop() {
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
