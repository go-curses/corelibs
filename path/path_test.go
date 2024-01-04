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
	"path/filepath"
	"syscall"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func cleanup(path string) {
	if IsDir(path) {
		_ = os.Chmod(path, 0770)
		_ = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				_ = os.Chmod(path, 0770)
			} else {
				_ = os.Chmod(path, 0660)
			}
			return nil
		})
		_ = os.RemoveAll(path)
	} else if IsFile(path) {
		_ = os.Remove(path)
	}
}

func mkFile(contents string) (path string, info os.FileInfo, err error) {
	var tmp *os.File
	if tmp, err = os.CreateTemp("", "corelibs-path-*.tmp"); err != nil {
		return
	} else if _, err = tmp.WriteString(contents); err != nil {
		return
	}
	path = tmp.Name()
	info, _ = tmp.Stat()
	_ = tmp.Close()
	return
}

func mkDir(perms os.FileMode) (path string, err error) {
	if path, err = os.MkdirTemp("", "corelibs-path-*.d"); err == nil {
		err = os.Chmod(path, perms)
	}
	return
}

func TestPath(t *testing.T) {
	Convey("Permissions", t, func() {
		perms, err := Permissions(".")
		So(err, ShouldEqual, nil)
		So(perms, ShouldNotEqual, 0)
		perms, err = Permissions("./nope")
		So(err, ShouldNotEqual, nil)
		So(perms, ShouldEqual, 0)
	})

	Convey("Overwrite, OverwriteWithPerms", t, func() {
		var data []byte
		var stat os.FileInfo
		var perms os.FileMode

		oldMask := syscall.Umask(0)
		defer syscall.Umask(oldMask)

		// create a temp file with some data and a specific mode
		path, _, err := mkFile("testing\n")
		So(err, ShouldEqual, nil)
		defer cleanup(path)
		err = os.Chmod(path, 0644)
		So(err, ShouldEqual, nil)
		stat, err = os.Stat(path)
		So(err, ShouldEqual, nil)

		// overwrite the temp file, check for correct mode and contents
		err = Overwrite(path, "success\n")
		perms, err = Permissions(path)
		So(err, ShouldEqual, nil)
		So(perms.Perm(), ShouldEqual, stat.Mode().Perm())
		data, err = os.ReadFile(path)
		So(err, ShouldEqual, nil)
		So(string(data), ShouldEqual, "success\n")

		// remove the file and overwrite again, this time the mode
		// should be the DefaultFileMode
		_ = os.Remove(path)
		err = Overwrite(path, "moar success!\n")
		So(err, ShouldEqual, nil)
		perms, err = Permissions(path)
		So(err, ShouldEqual, nil)
		So(perms.Perm(), ShouldEqual, DefaultFileMode.Perm())
		data, err = os.ReadFile(path)
		So(err, ShouldEqual, nil)
		So(string(data), ShouldEqual, "moar success!\n")

		// verify WriteFile error case
		err = OverwriteWithPerms(".", "nope!", 0640)
		So(err, ShouldNotEqual, nil)

		// verify os.Stat error case
		//_ = os.Remove(path)
		err = OverwriteWithPerms(path, "testing", 0640)
		So(err, ShouldEqual, nil)
	})

	Convey("BackupAndOverwrite", t, func() {
		path, _, err := mkFile("original\n")
		So(err, ShouldEqual, nil)
		defer cleanup(path)
		var b0, b1 string
		b0, err = BackupAndOverwrite(path, "modified\n")
		So(err, ShouldEqual, nil)
		defer cleanup(b0)
		var data []byte
		data, err = os.ReadFile(b0)
		So(err, ShouldEqual, nil)
		So(string(data), ShouldEqual, "original\n")
		data, err = os.ReadFile(path)
		So(err, ShouldEqual, nil)
		So(string(data), ShouldEqual, "modified\n")
		b1, err = BackupAndOverwrite(path, "updated!\n")
		So(err, ShouldEqual, nil)
		defer cleanup(b1)
		data, err = os.ReadFile(path)
		So(string(data), ShouldEqual, "updated!\n")
		data, err = os.ReadFile(b1)
		So(string(data), ShouldEqual, "modified\n")
	})

	Convey("MoveFile", t, func() {
		So(MoveFile("nope", "not-a-thing"), ShouldNotEqual, nil)
		So(MoveFile("/dev/null", "not-a-thing"), ShouldNotEqual, nil)
		path, _, err := mkFile("1\n")
		So(err, ShouldEqual, nil)
		defer cleanup(path)
		So(MoveFile(path, path), ShouldEqual, nil)
		So(MoveFile(path, path+".1"), ShouldEqual, nil)
		defer cleanup(path + ".1")
		So(Overwrite(path, "2\n"), ShouldEqual, nil)
		// this forces os.Rename to fail, but also forces CopyFile to fail too
		var dir string
		dir, err = mkDir(0550)
		So(err, ShouldEqual, nil)
		defer cleanup(dir)
		So(MoveFile(path, dir+"/1"), ShouldNotEqual, nil)
	})

	Convey("PruneEmptyDirs, ChmodAll", t, func() {
		var dirs []string
		// create the parent directory
		dir0, err := mkDir(0770)
		So(err, ShouldEqual, nil)
		defer cleanup(dir0)
		// create the first subdirectory
		dir1 := dir0 + "/one"
		err = os.Mkdir(dir1, 0770)
		So(err, ShouldEqual, nil)
		defer cleanup(dir1)
		// create the second subdirectory
		dir2 := dir0 + "/.two"
		err = os.Mkdir(dir2, 0770)
		So(err, ShouldEqual, nil)
		defer cleanup(dir2)
		// add a file to the first subdirectory
		err = os.WriteFile(dir1+"/file.txt", []byte("test!\n"), 0660)
		So(err, ShouldEqual, nil)
		// check the dirs are there
		dirs, err = ListDirs(dir0, true)
		So(err, ShouldEqual, nil)
		So(dirs, ShouldEqual, []string{dir2, dir1})
		// prune
		err = PruneEmptyDirs(dir0)
		So(err, ShouldEqual, nil)
		// check only dir1 is there
		dirs, err = ListDirs(dir0, true)
		So(err, ShouldEqual, nil)
		So(dirs, ShouldEqual, []string{dir1})
		// create dir2 again
		err = os.Mkdir(dir2, 0770)
		So(err, ShouldEqual, nil)
		dirs, err = ListDirs(dir0, true)
		So(err, ShouldEqual, nil)
		So(dirs, ShouldEqual, []string{dir2, dir1})
		// change dir0 perms to read-only
		_ = os.Chmod(dir0, 0550)
		// prune again, expecting os.Remove call to err
		err = PruneEmptyDirs(dir0)
		So(err, ShouldNotEqual, nil)
		_ = os.Chmod(dir0, 0770)
		_ = os.Chmod(dir1, 0330) // wx
		// prune again, expecting ListFiles call to err
		err = PruneEmptyDirs(dir0)
		So(err, ShouldNotEqual, nil)
		_ = os.Chmod(dir0, 0330)
		// prune again, expecting ListDirs call to err
		err = PruneEmptyDirs(dir0)
		So(err, ShouldNotEqual, nil)
	})

	Convey("CopyFile", t, func() {
		dir, err := mkDir(0770)
		So(err, ShouldEqual, nil)
		defer cleanup(dir)
		file0, file1 := dir+"/file.txt", dir+"/file.txt.1"
		_ = os.WriteFile(file0, []byte("text\n"), 0220)
		_, err = CopyFile(file0, file1)
		So(err, ShouldNotEqual, nil)
		_, err = CopyFile("nope", "not-a-thing")
		So(err, ShouldNotEqual, nil)
		_, err = CopyFile("/dev/null", "/not-a-thing")
		So(err, ShouldNotEqual, nil)
	})
}