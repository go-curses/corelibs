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

// Package spinner implements a means of rendering unicode characters that
// give a sense that something is happening and that the user is waiting
// during that process.
//
// The Spinner constructed does not actually write or otherwise output
// anything in particular. Instead, it simply provides the current
// unicode character for that moment or iteration depending on whether
// Spinner.Start was called or if the Spinner.Increment method is used
// to step to the next logical character in the loop.
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

// Callback is the function signature for the `fn` argument to New
type Callback func(symbol string)

// Spinner is the interface for interacting with the New Spinner
// instance, all methods are concurrency safe
type Spinner interface {
	// String returns the current symbol
	String() (symbol string)
	// Increment moves the spinner to the symbol
	Increment()
	// Start simply calls StartWith an interval of 250ms
	Start()
	// StartWith (and Start) are blocking operations, use a goroutine
	// to start a spinner in the background. If using within a Go-Curses
	// environment, use cdk.Go to invoke the goroutine. Once started,
	// a time.Ticker is used to call the New `fn` (if not nil) and
	// increment the spinner state once every interval
	StartWith(interval time.Duration)
	// Stop stops any started incrementing
	Stop()
}

type cSpinner struct {
	index    int
	symbols  []string
	ticker   *time.Ticker
	callback Callback

	sync.RWMutex
}

// New constructs a new Spinner instance
func New(symbols []string, fn Callback) (s Spinner) {
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
