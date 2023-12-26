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

import "regexp"

const KeywordPattern = `[a-zA-Z0-9]+?[-'a-zA-Z0-9]*[a-zA-Z0-9]+?|[a-zA-Z0-9]+`

var (
	RxNonWord     = regexp.MustCompile(`[^a-zA-Z0-9]`)
	RxEmptySpace  = regexp.MustCompile(`\s+`)
	RxLanguageKey = regexp.MustCompile(`language:(\*|[a-z][-a-zA-Z]+)\s*`)
	RxKeyword     = regexp.MustCompile(`^(` + KeywordPattern + `)$`)
	RxKeywords    = regexp.MustCompile(`\b(` + KeywordPattern + `)\b`)
)