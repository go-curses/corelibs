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

package maths

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClamps(t *testing.T) {

	Convey("Clamp", t, func() {
		So(Clamp(0, 1, 2), ShouldEqual, 1)
		So(Clamp(3, 1, 2), ShouldEqual, 2)
		So(Clamp(1, 1, 2), ShouldEqual, 1)
		So(Clamp(2, 1, 2), ShouldEqual, 2)
	})

	Convey("Floor", t, func() {
		So(Floor(0, 0), ShouldEqual, 0)
		So(Floor(1, 0), ShouldEqual, 1)
		So(Floor(1, 10), ShouldEqual, 10)
	})

	Convey("Ceil", t, func() {
		So(Ceil(0, 0), ShouldEqual, 0)
		So(Ceil(1, 10), ShouldEqual, 1)
		So(Ceil(10, 1), ShouldEqual, 1)
	})

	Convey("Round, Up, Down", t, func() {
		So(Round(1.0), ShouldEqual, 1)
		So(Round(0.5), ShouldEqual, 1)
		So(Round(0.25), ShouldEqual, 0)
		So(RoundUp(0.5), ShouldEqual, 1)
		So(RoundUp(1.0), ShouldEqual, 1)
		So(RoundDown(0.5), ShouldEqual, 0)
		So(RoundDown(1.0), ShouldEqual, 1)
	})

}
