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

package filewriter

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWriter(t *testing.T) {
	tempDir := os.TempDir()
	if tempDir == "" {
		t.Error("unable to get temporary directory")
		return
	}
	tempFile := fmt.Sprintf("%s/filewriter.%d.tmp", tempDir, os.Getpid())
	defer os.Remove(tempFile)

	Convey("New File", t, func() {
		fw, err := New().SetFile(tempFile).Make()
		So(err, ShouldEqual, nil)
		So(fw, ShouldNotEqual, nil)
		So(fw.File(), ShouldEqual, tempFile)
		So(fw.Mode(), ShouldEqual, DefaultFileMode)
		So(fw.Remove(), ShouldEqual, nil)
	})

	Convey("New Temp", t, func() {
		fw, err := New().UseTemp("filewriter.*.tmp").Make()
		So(err, ShouldEqual, nil)
		So(fw, ShouldNotEqual, nil)
		So(fw.Mode(), ShouldEqual, DefaultFileMode)
		file := fw.File()
		So(strings.Contains(file, "/filewriter."), ShouldEqual, true)
		So(strings.HasSuffix(file, ".tmp"), ShouldEqual, true)
		So(fw.Remove(), ShouldEqual, nil)
	})

	Convey("New Default", t, func() {
		fw, err := New().Make()
		So(err, ShouldEqual, nil)
		So(fw, ShouldNotEqual, nil)
		So(fw.Mode(), ShouldEqual, DefaultFileMode)
		file := fw.File()
		So(strings.Contains(file, "/filewriter-"), ShouldEqual, true)
		So(strings.HasSuffix(file, ".tmp"), ShouldEqual, true)
		So(fw.Remove(), ShouldEqual, nil)
	})

	Convey("Custom Mode", t, func() {
		fw, err := New().SetMode(0600).Make()
		So(err, ShouldEqual, nil)
		So(fw, ShouldNotEqual, nil)
		So(fw.Mode(), ShouldEqual, 0600)
		file := fw.File()
		var stat os.FileInfo
		stat, err = os.Stat(file)
		So(err, ShouldEqual, nil)
		So(stat.Mode(), ShouldEqual, 0600)
		So(fw.Remove(), ShouldEqual, nil)
	})

	Convey("Read/Write", t, func() {
		Convey("Functional", func() {
			fw, err := New().SetFile(tempFile).Make()
			So(err, ShouldEqual, nil)
			So(fw, ShouldNotEqual, nil)

			var count int
			So(fw.WalkFile(func(line string) (stop bool) {
				count += 1
				return
			}), ShouldEqual, false)
			So(count, ShouldEqual, 0)

			count, err = fw.Write([]byte("stuff"))
			So(err, ShouldEqual, nil)
			So(count, ShouldEqual, 5)

			count, err = fw.WriteString("\nmoar")
			So(err, ShouldEqual, nil)
			So(count, ShouldEqual, 5)

			var data []byte
			data, err = fw.ReadFile()
			So(err, ShouldEqual, nil)
			So(data, ShouldEqual, []byte("stuff\nmoar"))

			count = 0
			So(fw.WalkFile(func(line string) (stop bool) {
				if count == 0 {
					So(line, ShouldEqual, "stuff")
				} else if count == 1 {
					So(line, ShouldEqual, "moar")
				} else {
					t.Error("too many lines")
				}
				count += 1
				return
			}), ShouldEqual, false)
			So(fw.WalkFile(func(line string) (stop bool) {
				stop = true
				return
			}), ShouldEqual, true)

			So(fw.Close(), ShouldEqual, nil)
			count, err = fw.Write([]byte("nope"))
			So(err, ShouldEqual, os.ErrClosed)
			So(count, ShouldEqual, 0)
			So(fw.Remove(), ShouldEqual, nil)
		})

		Convey("Make Errors", func() {
			// testing os.CreateTemp errors
			fw, err := New().
				UseTemp("/nope/filewriter-*.tmp").
				Make()
			So(err, ShouldNotEqual, nil)
			So(fw, ShouldEqual, nil)

			// TODO: how to test os.Chmod errors?
			// TODO: how to test filepath.Abs errors?
		})

		Convey("Panic Errors", func() {

			tmpFile, _ := os.CreateTemp("", "filewriter-panic-test-*.tmp")
			tmpName := tmpFile.Name()
			_ = tmpFile.Chmod(0000)
			_ = tmpFile.Close()

			fw, err := New().SetFile(tmpName).Make()
			So(err, ShouldEqual, nil)
			_, err = fw.Write([]byte(""))
			So(err, ShouldNotEqual, nil)
			So(func() {
				fw.WalkFile(func(line string) (stop bool) {
					return
				})
			}, ShouldPanic)

			_ = os.Remove(tmpName)
		})

	})
}