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

package convert

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRoman(t *testing.T) {

	Convey("ToRoman", t, func() {
		letters := ToRoman(0)
		So(letters, ShouldEqual, "")
		letters = ToRoman(34)
		So(letters, ShouldEqual, "XXXIV")
		So(ToRoman(4000), ShouldEqual, "4000")
	})

}