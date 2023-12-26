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

package errors

import (
	"errors"
	"fmt"
)

func New(messages ...string) error {
	text := ""
	if len(messages) > 0 {
		for idx, message := range messages {
			if idx == 0 {
				text = fmt.Sprintf("%v", message)
			} else {
				text = fmt.Sprintf("%v\n%v", text, message)
			}
		}
	}
	return errors.New(text)
}

func NewF(format string, argv ...interface{}) error {
	return errors.New(fmt.Sprintf(format, argv...))
}

func NewPrefixed(prefix string, messages ...string) error {
	text := ""
	if prefix != "" {
		for idx, message := range messages {
			if idx == 0 {
				text = fmt.Sprintf("%v: %v", prefix, message)
			} else {
				text = fmt.Sprintf("%v\n%v: %v", text, prefix, message)
			}
		}
	} else {
		for idx, message := range messages {
			if idx == 0 {
				text = message
			} else {
				text = fmt.Sprintf("%v\n%v", text, message)
			}
		}
	}
	return errors.New(text)
}

func NewPrefixedF(prefix, format string, argv ...interface{}) error {
	fixed := format
	if prefix != "" {
		fixed = fmt.Sprintf("%v: %v", prefix, format)
	}
	return NewF(fixed, argv...)
}