// Copyright (c) 2024  The Go-Curses Authors
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

package path

import (
	"io/fs"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWrappers(t *testing.T) {
	Convey("(all)", t, func() {
		cwd, _ := os.Getwd()
		// Abs
		path, err := Abs("./thing/..")
		So(err, ShouldEqual, nil)
		So(path, ShouldEqual, cwd)
		// Clean
		So(Clean("/thing/.."), ShouldEqual, "/")
		// Dir
		So(Dir("/thing/file.txt"), ShouldEqual, "/thing")
		// Walk
		So(Walk(".", func(path string, info fs.FileInfo, err error) error {
			return nil
		}), ShouldEqual, nil)
		// ReadDir
		_, err = ReadDir(".")
		So(err, ShouldEqual, nil)
		// ReadFile
		_, err = ReadFile("wrappers_test.go")
		So(err, ShouldEqual, nil)
		// Stat
		_, err = Stat(".")
		So(err, ShouldEqual, nil)
		// Which
		path = Which("go")
		So(path, ShouldNotEqual, "")
	})
}