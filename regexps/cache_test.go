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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNotify(t *testing.T) {
	Convey("Compile and ClearCache", t, func() {
		var err error
		var rx0, rx1 *regexp.Regexp
		rx0, err = Compile(`.*`)
		So(err, ShouldEqual, nil)
		So(rx0, ShouldNotEqual, nil)
		rx1, err = Compile(`.*`)
		So(rx1, ShouldNotEqual, nil)
		So(rx0, ShouldEqual, rx1)
		_, err = Compile(`[Broken`)
		So(err, ShouldNotEqual, nil)
		So(len(_cache.data), ShouldEqual, 2)
		ClearCache()
		So(len(_cache.data), ShouldEqual, 0)
	})
}
