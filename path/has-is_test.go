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
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHasIs(t *testing.T) {
	Convey("IsHidden", t, func() {
		So(IsHidden("/path/.hidden"), ShouldEqual, true)
		So(IsHidden("/path/not-hid"), ShouldEqual, false)
	})

	Convey("IsBackup", t, func() {
		So(IsBackup("/path/name~"), ShouldEqual, true)
		So(IsBackup("/path/nope"), ShouldEqual, false)
	})

	Convey("IsHiddenPath", t, func() {
		So(IsHiddenPath("/path/.hidden/file"), ShouldEqual, true)
		So(IsHiddenPath("/path/not-hid/file"), ShouldEqual, false)
	})

	Convey("IsPlainText", t, func() {
		So(IsPlainText("."), ShouldEqual, false)
		So(IsPlainText("./has-is_test.go"), ShouldEqual, true)
		tmp, _ := os.CreateTemp("", "corelibs-has-is-*.tmp.html")
		tmpName := tmp.Name()
		defer os.Remove(tmpName)
		_, _ = tmp.WriteString(`<!DOCTYPE html><head><title>test</title></head><body><p>test</p></body></html>`)
		_ = tmp.Close()
		So(IsPlainText(tmpName), ShouldEqual, true)
	})

	Convey("Exists", t, func() {
		So(Exists("."), ShouldEqual, true)
		So(Exists("./has-is_test.go"), ShouldEqual, true)
		So(Exists("./not-a-thing"), ShouldEqual, false)
	})

	Convey("IsFile", t, func() {
		So(IsFile("."), ShouldEqual, false)
		So(IsFile("./has-is_test.go"), ShouldEqual, true)
	})

	Convey("IsRegularFile", t, func() {
		So(IsRegularFile("."), ShouldEqual, false)
		So(IsRegularFile("nope"), ShouldEqual, false)
		So(IsRegularFile("/dev/null"), ShouldEqual, false)
		So(IsRegularFile("./has-is_test.go"), ShouldEqual, true)
	})

	Convey("IsDir", t, func() {
		So(IsDir("."), ShouldEqual, true)
		So(IsDir("nope"), ShouldEqual, false)
		So(IsDir("./has-is_test.go"), ShouldEqual, false)
	})

	Convey("HasExt", t, func() {
		So(HasExt("file.txt", ""), ShouldEqual, false)
		So(HasExt("file.txt", ".txt"), ShouldEqual, true)
		So(HasExt("file.txt.tmpl", "txt"), ShouldEqual, true)
		So(HasExt("file.txt.tmpl", "html"), ShouldEqual, false)
	})

	Convey("HasAnyExt", t, func() {
		So(HasAnyExt("file.txt", ""), ShouldEqual, false)
		So(HasAnyExt("file.txt", ".nope", ".txt"), ShouldEqual, true)
	})

}