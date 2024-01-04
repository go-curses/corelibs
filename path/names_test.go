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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNames(t *testing.T) {
	Convey("IncrementLabel", t, func() {
		title, _ := IncrementLabel("", "", "", "", 0)
		So(title, ShouldEqual, "")
		title, _ = IncrementLabel("Title", "", "", "", 0)
		So(title, ShouldEqual, "Title1")
		title, _ = IncrementLabel("Title10", "", "", "", 1)
		So(title, ShouldEqual, "Title11")
		title, _ = IncrementLabel("Title10", "", "", "", -1)
		So(title, ShouldEqual, "Title9")
	})

	Convey("IncrementFileName", t, func() {
		So(IncrementFileName("Name"), ShouldEqual, "Name (1)")
		So(IncrementFileName("Name (1)"), ShouldEqual, "Name (2)")
	})

	Convey("IncrementFilePath", t, func() {
		So(IncrementFilePath("/file.txt"), ShouldEqual, "/file.txt.1")
		So(IncrementFilePath("/file.txt.1"), ShouldEqual, "/file.txt.2")
	})

	Convey("IncrementFileBackup", t, func() {
		So(IncrementFileBackup("/file.txt", ""), ShouldEqual, "/file.txt~")
		So(IncrementFileBackup("/file.txt~", ""), ShouldEqual, "/file.txt.1~")
		So(IncrementFileBackup("/file.txt", ".backup"), ShouldEqual, "/file.txt.backup")
		So(IncrementFileBackup("/file.txt.1.backup", ".backup"), ShouldEqual, "/file.txt.2.backup")
	})

	Convey("BackupName", t, func() {
		Convey("zero argv", func() {
			So(BackupName("file.txt"), ShouldEqual, "file.txt~")
			So(BackupName("file.txt~"), ShouldEqual, "file.txt.1~")
		})

		Convey("one argv", func() {
			So(BackupName("file.txt", 1), ShouldEqual, "file.txt")
			So(BackupName("file.txt", "^"), ShouldEqual, "file.txt^")
			So(BackupName("file.txt^", "^"), ShouldEqual, "file.txt.1^")
		})

		Convey("two argv", func() {
			So(BackupName("file.txt", "^", "~"), ShouldEqual, "file.txt^")
			So(BackupName("file.txt^", "^", "~"), ShouldEqual, "file.txt~1^")
		})

		Convey("three argv", func() {
			So(BackupName("file.txt", "^", "~", "10"), ShouldEqual, "file.txt^")
			So(BackupName("file.txt", "^", "~", -1), ShouldEqual, "file.txt^")
			So(BackupName("file.txt~1^", "^", "~", -10), ShouldEqual, "file.txt^")
			So(BackupName("file.txt", "^", "~", 10), ShouldEqual, "file.txt^")
			So(BackupName("file.txt^", "^", "~", 10), ShouldEqual, "file.txt~10^")
		})

		Convey("four argv", func() {
			So(BackupName("file.txt", "^", "~", "[", "]"), ShouldEqual, "file.txt^")
			So(BackupName("file.txt^", "^", "~", "[", "]"), ShouldEqual, "file.txt~[1]^")
		})

		Convey("five argv", func() {
			So(BackupName("file.txt", "^", "~", "[", "]", 10), ShouldEqual, "file.txt^")
			So(BackupName("file.txt^", "^", "~", "[", "]", 10), ShouldEqual, "file.txt~[10]^")
		})

	})
}