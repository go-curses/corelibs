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

func TestLetters(t *testing.T) {

	Convey("ToLetters", t, func() {
		letters := ToLetters(0)
		So(letters, ShouldEqual, "a")
		letters = ToLetters(34)
		So(letters, ShouldEqual, "ai")
	})

	Convey("ToCharacters", t, func() {
		base := "0123456789ABCDEF"
		characters := ToCharacters(1, base)
		So(characters, ShouldEqual, "1")
		characters = ToCharacters(11, base)
		So(characters, ShouldEqual, "B")
		characters = ToCharacters(34, base)
		So(characters, ShouldEqual, "12")
		So(ToCharacters(10, ""), ShouldEqual, "")
	})
}