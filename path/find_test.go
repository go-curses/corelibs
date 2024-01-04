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

func TestFind(t *testing.T) {
	Convey("Find", t, func() {
		var err error
		var dir string
		if dir, err = os.MkdirTemp("", "corelibs-find-*.tmpd"); err != nil {
			t.Fatal(err)
		} else if err = os.MkdirAll(dir+"/one/two", 0775); err != nil {
			t.Fatal(err)
		} else if err = os.WriteFile(dir+"/.check", []byte("check\n"), 0644); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		dotCheck := FindFileRelativeToPath(".check", dir+"/one/two")
		So(dotCheck, ShouldEqual, dir+"/.check")
		dotCheck = FindFileRelativeToPath(".check", dir)
		So(dotCheck, ShouldEqual, dir+"/.check")
		dotCheck = FindFileRelativeToPath(".check", "./.../")
		So(dotCheck, ShouldEqual, "")

		cwd, _ := os.Getwd()
		_ = os.Chdir(dir + "/one/two")
		dotCheck = FindFileRelativeToPwd(".check")
		So(dotCheck, ShouldEqual, dir+"/.check")
		_ = os.Chdir(cwd)
	})
}