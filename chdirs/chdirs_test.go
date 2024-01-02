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

package chdirs

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("error calling os.Getwd: %v", err)
	}

	Convey("push, stack, pop", t, func() {
		Convey("push", func() {
			So(Push(".."), ShouldEqual, nil)
			So(Push(".."), ShouldEqual, nil)
		})

		Convey("stack", func() {
			items := Stack()
			So(len(items), ShouldEqual, 2)
			So(items[0], ShouldEqual, cwd)
		})

		Convey("pop", func() {
			So(Pop(), ShouldEqual, nil)
			So(len(Stack()), ShouldEqual, 1)
			So(Pop(), ShouldEqual, nil)
			So(len(Stack()), ShouldEqual, 0)
			twd, err := os.Getwd()
			So(err, ShouldEqual, nil)
			So(twd, ShouldEqual, cwd)
		})

		Convey("bad push", func() {
			So(Push("/!not!a!path"), ShouldNotEqual, nil)
		})
	})
}