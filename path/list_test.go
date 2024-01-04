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

func TestList(t *testing.T) {
	Convey("List", t, func() {
		paths, err := List(".", false)
		So(err, ShouldEqual, nil)
		So(paths, ShouldNotBeEmpty)
		paths, err = List("nope", false)
		So(err, ShouldNotEqual, nil)
		So(paths, ShouldBeEmpty)
		paths, err = List("./list_test.go", false)
		So(err, ShouldEqual, nil)
		So(paths, ShouldEqual, []string{"list_test.go"})
	})

	Convey("ListDirs", t, func() {
		dirs, err := ListDirs(".", false)
		So(err, ShouldEqual, nil)
		So(dirs, ShouldBeEmpty)
		dirs, err = ListDirs("nope", false)
		So(err, ShouldNotEqual, nil)
		So(dirs, ShouldBeEmpty)
		dirs, err = ListDirs("..", false)
		So(err, ShouldEqual, nil)
		So(dirs, ShouldNotBeEmpty)
	})

	Convey("ListFiles", t, func() {
		files, err := ListFiles(".", false)
		So(err, ShouldEqual, nil)
		So(files, ShouldNotBeEmpty)
		files, err = ListFiles("nope", false)
		So(err, ShouldNotEqual, nil)
		So(files, ShouldBeEmpty)
	})

	Convey("ListPaths", t, func() {
		l := ListPaths{
			{
				Path: "one",
			},
			{
				Path: ".two",
			},
		}
		l.Sort()
		So(l[0].Path, ShouldEqual, ".two")
		So(l[1].Path, ShouldEqual, "one")
	})

	Convey("ListAllDirs, ListAllFiles", t, func() {
		tmpDir, _ := os.MkdirTemp("", "corelibs-list-all-*.tmp")
		defer func() {
			_ = os.Chmod(tmpDir+"/dir/another", 0775)
			_ = os.RemoveAll(tmpDir)
		}()
		_ = os.WriteFile(tmpDir+"/top", []byte("top"), 0664)
		_ = os.MkdirAll(tmpDir+"/dir/sub-dir", 0775)
		_ = os.MkdirAll(tmpDir+"/dir/another", 0775)
		_ = os.MkdirAll(tmpDir+"/dir/.hiding", 0775)
		_ = os.WriteFile(tmpDir+"/dir/.hidden", []byte("hidden"), 0664)
		_ = os.WriteFile(tmpDir+"/dir/another/file", []byte("file"), 0664)
		things, err := ListAllDirs(tmpDir, false)
		So(err, ShouldEqual, nil)
		So(things, ShouldEqual, []string{
			tmpDir + "/dir",
			tmpDir + "/dir/another",
			tmpDir + "/dir/sub-dir",
		})
		things, err = ListAllDirs(tmpDir, true)
		So(err, ShouldEqual, nil)
		So(things, ShouldEqual, []string{
			tmpDir + "/dir",
			tmpDir + "/dir/.hiding",
			tmpDir + "/dir/another",
			tmpDir + "/dir/sub-dir",
		})
		things, err = ListAllDirs("nope", false)
		So(err, ShouldNotEqual, nil)
		So(things, ShouldBeEmpty)
		things, err = ListAllFiles(tmpDir, false)
		So(err, ShouldEqual, nil)
		So(things, ShouldEqual, []string{
			tmpDir + "/dir/another/file",
			tmpDir + "/top",
		})
		things, err = ListAllFiles(tmpDir, true)
		So(err, ShouldEqual, nil)
		So(things, ShouldEqual, []string{
			tmpDir + "/dir/another/file",
			tmpDir + "/dir/.hidden",
			tmpDir + "/top",
		})
		things, err = ListAllFiles("nope", false)
		So(err, ShouldNotEqual, nil)
		So(things, ShouldBeEmpty)
		_ = os.Chmod(tmpDir+"/dir/another", 0330)
		things, err = ListAllFiles(tmpDir, false)
		So(err, ShouldNotEqual, nil)
		So(things, ShouldBeEmpty)
		things, err = ListAllFiles(tmpDir+"/dir/.hidden", false)
		So(err, ShouldEqual, nil)
		So(things, ShouldBeEmpty)
		things, err = ListAllFiles(tmpDir+"/dir/.hidden", true)
		So(err, ShouldEqual, nil)
		So(things, ShouldEqual, []string{
			tmpDir + "/dir/.hidden",
		})
	})
}